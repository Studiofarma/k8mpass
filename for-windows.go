package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"github.com/studiofarma/k8mpass/api"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func Age(ns v1.Namespace) (api.ExtensionValue, error) {
	res := fmt.Sprintf("Age: %0.f minutes", time.Since(ns.CreationTimestamp.Time).Minutes())
	return api.ExtensionValue(res), nil
}

func AgeList(ns []v1.Namespace) map[api.Name]api.ExtensionValue {
	values := make(map[api.Name]api.ExtensionValue)
	for _, n := range ns {
		age, _ := Age(n)
		values[api.Name(n.Name)] = age
	}
	return values
}

var AgeProperty = api.Extension{
	Name:         "age",
	ExtendSingle: Age,
	ExtendList:   AgeList,
}

var OpenApplicationOperation = api.NamespaceOperation{
	Name:      "Open application in browser",
	Condition: OpenApplicationCondition,
	Command: func(clientset *kubernetes.Clientset, namespace string) tea.Cmd {
		return func() tea.Msg {
			ingresses, err := clientset.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})

			if err != nil {
				return api.NoOutputResultMsg{false, err.Error()}
			}

			var dbmsUrl string

			for _, i := range ingresses.Items {
				if strings.HasPrefix(i.Name, "g3pharmacy") {
					dbmsUrl = i.Spec.Rules[0].Host
				}
			}
			if dbmsUrl == "" {
				return api.NoOutputResultMsg{false, "Ingress not found"}
			}
			Openbrowser("https://" + dbmsUrl)

			return api.NoOutputResultMsg{true, "App is ready"}
		}
	},
}

func OpenApplicationCondition(cs *kubernetes.Clientset, namespace string) bool {
	ingresses, err := cs.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return false
	}
	res := false
	for _, i := range ingresses.Items {
		host := i.Spec.Rules[0].Host
		res = res || strings.HasPrefix(host, "g3pharmacy")
	}
	return res
}

var OpenDbmsOperation = api.NamespaceOperation{
	Name:      "Open DBMS in browser",
	Condition: OpenDbmsCondition,
	Command: func(clientset *kubernetes.Clientset, namespace string) tea.Cmd {
		return func() tea.Msg {
			ingresses, err := clientset.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})

			if err != nil {
				return api.NoOutputResultMsg{false, err.Error()}
			}

			var dbmsUrl string

			for _, i := range ingresses.Items {
				host := i.Spec.Rules[0].Host
				if strings.HasPrefix(host, "dbms") {
					dbmsUrl = host
				}
			}
			if dbmsUrl == "" {
				return api.NoOutputResultMsg{false, "Ingress not found"}
			}
			Openbrowser("https://" + dbmsUrl)

			return api.NoOutputResultMsg{
				true,
				"DBeaver is better 🦦",
			}
		}
	},
}

func OpenDbmsCondition(cs *kubernetes.Clientset, namespace string) bool {
	ingresses, err := cs.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return false
	}
	res := false
	for _, i := range ingresses.Items {
		host := i.Spec.Rules[0].Host
		res = res || strings.HasPrefix(host, "dbms")
	}
	return res
}

const (
	sleeping = "Sleeping..."
	awake    = "Awake!"
)

func GetNamespaceExtensions() []api.IExtension {
	return namespaceExtensions
}

func GetNamespaceOperations() []api.INamespaceOperation {
	return namespaceOperations
}

var namespaceOperations = []api.INamespaceOperation{
	WakeUpReviewOperation,
	OpenDbmsOperation,
	OpenApplicationOperation,
}

var namespaceExtensions = []api.IExtension{
	ReviewAppSleepStatus,
	//AgeProperty,
}

var ReviewAppSleepStatus = api.Extension{
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

func IsReviewAppSleeping(ns v1.Namespace) (api.ExtensionValue, error) {
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
	return api.ExtensionValue(thResponse.Status()), nil
}
func AreReviewAppsSleeping(ns []v1.Namespace) map[api.Name]api.ExtensionValue {
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
	values := make(map[api.Name]api.ExtensionValue, len(ns))
	for _, n := range ns {
		values[api.Name(n.Name)] = api.ExtensionValue(thResponse.StatusByNamespace(n.Name))
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

var WakeUpReviewOperation = api.NamespaceOperation{
	Name:      "Wake up review app",
	Condition: WakeUpReviewCondition,
	Command: func(clientset *kubernetes.Clientset, namespace string) tea.Cmd {
		return func() tea.Msg {
			err := wakeupReview(clientset, namespace)
			if err != nil {
				return api.NoOutputResultMsg{false, err.Error()}
			}
			return api.NoOutputResultMsg{true, "We woke it up!"}
		}
	},
}

func wakeupReview(clientset *kubernetes.Clientset, namespace string) error {
	cronjobs := clientset.BatchV1().CronJobs(namespace)
	cronjob, err := cronjobs.Get(context.TODO(), "scale-to-zero-wakeup", metav1.GetOptions{})
	if err != nil {
		return err
	}

	newUUID, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	jobSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("k8mpass-wakeup-%s", newUUID.String()),
			Namespace: namespace,
		},
		Spec: cronjob.Spec.JobTemplate.Spec,
	}
	jobs := clientset.BatchV1().Jobs(namespace)

	_, err = jobs.Create(context.TODO(), jobSpec, metav1.CreateOptions{})

	return err
}

func WakeUpReviewCondition(cs *kubernetes.Clientset, namespace string) bool {
	_, err := cs.BatchV1().CronJobs(namespace).Get(context.TODO(), "scale-to-zero-wakeup", metav1.GetOptions{})
	if err != nil {
		return false
	}
	return true
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
