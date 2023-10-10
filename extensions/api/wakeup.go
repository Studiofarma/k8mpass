package api

import (
	"encoding/json"
	"errors"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/studiofarma/k8mpass/api"
	"k8s.io/client-go/kubernetes"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	v1 "k8s.io/api/core/v1"
)

const (
	sleeping = "Sleeping..."
	awake    = "Awake!"
)

var NamespaceOperations = []api.INamespaceOperation{
	api.NewCronjobOperation("Wake up!", "scale-to-zero-wakeup", "Wake up triggered"),
	api.NewIngressOperation("Open DBMS in browser", "dbms", "Now go delete the database"),
	api.NewIngressOperation("Open application in browser", "g3pharmacy", "It's not a bug, it's a feature"),
}

var NamespaceExtensions = []api.INamespaceExtension{
	ReviewAppSleepStatus,
	NamespaceAgeProperty,
	//PodVersion,
}

var PodExtensions = []api.IPodExtension{
	PodAgeProperty,
	PodVersion,
}

var ReviewAppSleepStatus = api.NamespaceExtension{
	Name:         "sleeping",
	ExtendSingle: IsReviewAppSleeping,
	ExtendList:   AreReviewAppsSleeping,
}

type ThanosResponse struct {
	Data ThanosData `json:"data"`
}

type ThanosData struct {
	Result []ThanosResult `json:"result"`
}

type ThanosMetric struct {
	ExportedService string `json:"exported_service"`
}

type ThanosResult struct {
	Metric ThanosMetric  `json:"metric"`
	Value  []interface{} `json:"value"`
}

func (r ThanosResponse) IsAwake() bool {
	if len(r.Data.Result) == 0 {
		return false
	}
	if r.Data.Result[0].Value[1] == "" || r.Data.Result[0].Value[1] == "0" {
		return false
	}
	return true
}

func (r ThanosResponse) IsAwakeByNamespace(namespace string) bool {
	var isAwake = false
	for _, ra := range r.Data.Result {
		if strings.HasPrefix(ra.Metric.ExportedService, namespace) {
			isAwake = ra.IsAwake() || isAwake
		}
	}
	return isAwake
}

func IsReviewApp(namespace string) bool {
	return strings.HasPrefix(namespace, "review")
}

func (r ThanosResponse) StatusByNamespace(namespace string) string {
	if !IsReviewApp(namespace) {
		return ""
	}
	if r.IsAwakeByNamespace(namespace) {
		return awake
	} else {
		return sleeping
	}
}

func (r ThanosResult) IsAwake() bool {
	if r.Value[1] == "" || r.Value[1] == "0" {
		return false
	}
	return true
}

func (r ThanosResult) Status() string {
	if r.IsAwake() {
		return awake
	} else {
		return sleeping
	}
}

func (r ThanosResponse) Status() string {
	if r.IsAwake() {
		return awake
	} else {
		return sleeping
	}
}

func IsReviewAppSleeping(ns v1.Namespace) (string, error) {
	if !IsReviewApp(ns.Name) {
		return "", nil
	}
	thanosUrl, isPresent := os.LookupEnv("THANOS_URL")
	if !isPresent {
		return "", errors.New("env var THANOS_URL not present")
	}
	query, isPresent := os.LookupEnv("THANOS_QUERY")
	if !isPresent {
		return "", errors.New("env var THANOS_QUERY not present")
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", thanosUrl+"/api/v1/query", nil)
	if err != nil {
		return "", err
	}
	q := req.URL.Query()
	q.Add("query", strings.Replace(query, "%NS%", ns.Name, 1))
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	var thResponse ThanosResponse
	err = json.NewDecoder(resp.Body).Decode(&thResponse)
	if err != nil {
		return "", err
	}
	return thResponse.Status(), nil
}
func AreReviewAppsSleeping(ns []v1.Namespace) map[string]string {
	thanosUrl, isPresent := os.LookupEnv("THANOS_URL")
	if !isPresent {
		return nil
	}
	query, isPresent := os.LookupEnv("THANOS_QUERY_ALL_NS")
	if !isPresent {
		return nil
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", thanosUrl+"/api/v1/query", nil)
	if err != nil {
		return nil
	}
	q := req.URL.Query()
	q.Add("query", query)
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	var thResponse ThanosResponse
	err = json.NewDecoder(resp.Body).Decode(&thResponse)
	if err != nil {
		return nil
	}
	values := make(map[string]string, len(ns))
	for _, n := range ns {
		values[n.Name] = thResponse.StatusByNamespace(n.Name)
	}

	return values
}

var CheckSleepingStatusOperation = api.NamespaceOperation{
	Name:      "Check if review app is asleep",
	Condition: CheckSleepingStatusCondition,
	Command: func(clientset *kubernetes.Clientset, namespace string) tea.Cmd {
		return checkIfReviewAppIsAsleep(namespace)
	},
}

func checkIfReviewAppIsAsleep(namespace string) tea.Cmd {
	return func() tea.Msg {
		client := &http.Client{}
		req, err := http.NewRequest("GET", os.Getenv("THANOS_URL")+"/api/v1/query", nil)
		if err != nil {
			return api.NoOutputResultMsg{Message: err.Error()}
		}
		q := req.URL.Query()
		query := os.Getenv("THANOS_QUERY")
		q.Add("query", strings.Replace(query, "%NS%", namespace, 1))
		req.URL.RawQuery = q.Encode()
		resp, err := client.Do(req)
		if err != nil {
			return api.NoOutputResultMsg{Message: err.Error()}
		}
		var thResponse ThanosResponse
		err = json.NewDecoder(resp.Body).Decode(&thResponse)
		if err != nil {
			return api.NoOutputResultMsg{Message: err.Error()}
		}
		return api.NoOutputResultMsg{Success: true, Message: thResponse.Status()}
	}
}

func CheckSleepingStatusCondition(*kubernetes.Clientset, string) bool {
	_, ok := os.LookupEnv("THANOS_URL")
	return ok
}

func Openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}
