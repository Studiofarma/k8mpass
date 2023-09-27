package main

import (
	"context"
	"encoding/json"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

func clusterConnect() tea.Msg {
	cs, err := getConnection()
	if err != nil {
		return errMsg(err)
	}
	K8sCluster = kubernetesCluster{cs}
	return clusterConnectedMsg{cs}
}

func fetchNamespaces() tea.Msg {
	time.Sleep(500 * time.Millisecond)
	ns, err := getNamespaces(K8sCluster.kubernetes)
	if err != nil {
		return errMsg(err)
	}
	var items []NamespaceItem
	sleepingInfo, err := getReviewAppsSleepingStatus()
	if err != nil {
		return errMsg(err)
	}
	for _, n := range ns.Items {
		var isSleeping = true
		for _, ra := range sleepingInfo {
			if strings.HasPrefix(ra.Metric.ExportedService, n.Name) {
				isSleeping = ra.IsAsleep()
			}
		}
		items = append(items, NamespaceItem{n, isSleeping})
	}
	return namespacesRetrievedMsg{items}
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

func checkIfReviewAppIsAsleep(namespace string) tea.Cmd {
	return func() tea.Msg {
		client := &http.Client{}
		req, err := http.NewRequest("GET", os.Getenv("THANOS_URL")+"/api/v1/query", nil)
		if err != nil {
			return errMsg(err)
		}
		q := req.URL.Query()
		query := os.Getenv("THANOS_QUERY")
		q.Add("query", strings.Replace(query, "%NS%", namespace, 1))
		req.URL.RawQuery = q.Encode()
		resp, err := client.Do(req)
		if err != nil {
			return errMsg(err)
		}
		var thResponse ThanosResponse
		err = json.NewDecoder(resp.Body).Decode(&thResponse)
		if thResponse.IsAsleep() {
			return noOutputResultMsg{false, "Review app is sleeping"}
		} else {
			return noOutputResultMsg{true, "Review app is awake"}
		}
	}
}

func getReviewAppsSleepingStatus() ([]ThanosResult, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", os.Getenv("THANOS_URL")+"/api/v1/query", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	query := os.Getenv("THANOS_QUERY_ALL_NS")
	q.Add("query", query)
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var thResponse ThanosResponse
	err = json.NewDecoder(resp.Body).Decode(&thResponse)
	return thResponse.Data.Result, nil
}

func (r ThanosResponse) IsAsleep() bool {
	if len(r.Data.Result) == 0 {
		return true
	}
	if r.Data.Result[0].Value[1] == "" || r.Data.Result[0].Value[1] == "0" {
		return true
	}
	return false
}
func (r ThanosResult) IsAsleep() bool {
	if r.Value[1] == "" || r.Value[1] == "0" {
		return true
	}
	return false
}

type K8mpassCommand func(model *kubernetes.Clientset, namespace string) tea.Cmd

type NamespaceOperation struct {
	Name    string
	Command K8mpassCommand
}

var WakeUpReviewOperation = NamespaceOperation{
	Name: "Wake up review app",
	Command: func(clientset *kubernetes.Clientset, namespace string) tea.Cmd {
		return func() tea.Msg {
			err := wakeupReview(clientset, namespace)
			if err != nil {
				return noOutputResultMsg{false, err.Error()}
			}
			return noOutputResultMsg{true, "We woke it up!"}
		}
	},
}
var CheckSleepingStatusOperation = NamespaceOperation{
	Name: "Check if review app is asleep",
	Command: func(clientset *kubernetes.Clientset, namespace string) tea.Cmd {
		return checkIfReviewAppIsAsleep(namespace)
	},
}

var PodsOperation = NamespaceOperation{
	Name: "Get list of pods",
	Command: func(clientset *kubernetes.Clientset, namespace string) tea.Cmd {
		return func() tea.Msg {
			p, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				return errMsg(err)
			}
			s := ""
			for _, pod := range p.Items {

				if pod.Status.Phase == v1.PodSucceeded || pod.Status.Phase == v1.PodFailed {
					continue
				}
				s += fmt.Sprintf("  %s\n", styleString(pod.Name, podStyle(pod.Status)))
			}
			return operationResultMsg{body: s}
		}
	},
}

var OpenDbmsOperation = NamespaceOperation{
	Name: "Open DBMS in browser",
	Command: func(clientset *kubernetes.Clientset, namespace string) tea.Cmd {
		return func() tea.Msg {
			ingresses, err := clientset.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})

			if err != nil {
				return noOutputResultMsg{false, err.Error()}
			}

			var dbmsUrl string

			for _, i := range ingresses.Items {
				host := i.Spec.Rules[0].Host
				if strings.HasPrefix(host, "dbms") {
					dbmsUrl = host
				}
			}
			if dbmsUrl == "" {
				return noOutputResultMsg{false, "Ingress not found"}
			}
			Openbrowser("https://" + dbmsUrl)

			return noOutputResultMsg{
				true,
				"DBeaver is better 🦦",
			}
		}
	},
}

var OpenApplicationOperation = NamespaceOperation{
	Name: "Open application in browser",
	Command: func(clientset *kubernetes.Clientset, namespace string) tea.Cmd {
		return func() tea.Msg {
			ingresses, err := clientset.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})

			if err != nil {
				return noOutputResultMsg{false, err.Error()}
			}

			var dbmsUrl string

			for _, i := range ingresses.Items {
				if strings.HasPrefix(i.Name, "g3pharmacy") {
					dbmsUrl = i.Spec.Rules[0].Host
				}
			}
			if dbmsUrl == "" {
				return noOutputResultMsg{false, "Ingress not found"}
			}
			Openbrowser("https://" + dbmsUrl)

			return noOutputResultMsg{true, "App is ready"}
		}
	},
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

func styleString(s string, style lipgloss.Style) lipgloss.Style {
	return style.SetString(s)
}

func podStyle(status v1.PodStatus) lipgloss.Style {

	switch status.Phase {
	case v1.PodRunning:
		var ready = true
		for _, c := range status.ContainerStatuses {
			ready = ready && c.Ready
		}
		if !ready {
			return lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ff6666"))
		}
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#66ffc2"))
	case v1.PodFailed, v1.PodPending:
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff6666"))
	default:
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#a6a6a6"))
	}

}
