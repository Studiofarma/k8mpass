package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func clusterConnect() tea.Msg {
	cs, err := getConnection()
	if err != nil {
		return errMsg(err)
	}
	K8sCluster = kubernetesCluster{cs}
	return clusterConnectedMsg{cs}
}

type K8mpassCommand func(model *kubernetes.Clientset, namespace string) tea.Cmd
type K8mpassCondition func(cs *kubernetes.Clientset, namespace string) bool

type NamespaceOperation struct {
	Name      string
	Command   K8mpassCommand
	Condition K8mpassCondition
}

var WakeUpReviewOperation = NamespaceOperation{
	Name:      "Wake up review app",
	Condition: WakeUpReviewCondition,
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

func WakeUpReviewCondition(cs *kubernetes.Clientset, namespace string) bool {
	_, err := cs.BatchV1().CronJobs(namespace).Get(context.TODO(), "scale-to-zero-wakeup", metav1.GetOptions{})
	if err != nil {
		return false
	}
	return true
}

var CheckSleepingStatusOperation = NamespaceOperation{
	Name:      "Check if review app is asleep",
	Condition: CheckSleepingStatusCondition,
	Command: func(clientset *kubernetes.Clientset, namespace string) tea.Cmd {
		return checkIfReviewAppIsAsleep(namespace)
	},
}

func CheckSleepingStatusCondition(cs *kubernetes.Clientset, namespace string) bool {
	_, ok := os.LookupEnv("THANOS_URL")
	return ok
}

var OpenDbmsOperation = NamespaceOperation{
	Name:      "Open DBMS in browser",
	Condition: OpenDbmsCondition,
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

var OpenApplicationOperation = NamespaceOperation{
	Name:      "Open application in browser",
	Condition: OpenApplicationCondition,
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
func CheckConditionsThatApply(cs *kubernetes.Clientset, namespace string, operations []NamespaceOperation) tea.Cmd {
	return func() tea.Msg {
		var availableOps []NamespaceOperation
		for _, operation := range operations {
			if operation.Condition == nil {
				continue
			}
			if operation.Condition(cs, namespace) {
				availableOps = append(availableOps, operation)
			}
		}
		return AvailableOperationsMsg{availableOps}
	}
}

type AvailableOperationsMsg struct {
	operations []NamespaceOperation
}
