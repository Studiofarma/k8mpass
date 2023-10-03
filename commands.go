package main

import (
	"context"
	"fmt"
	"github.com/studiofarma/k8mpass/api"
	"log"
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

func WakeUpReviewCondition(cs *kubernetes.Clientset, namespace string) bool {
	_, err := cs.BatchV1().CronJobs(namespace).Get(context.TODO(), "scale-to-zero-wakeup", metav1.GetOptions{})
	if err != nil {
		return false
	}
	return true
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
func CheckConditionsThatApply(cs *kubernetes.Clientset, namespace string, operations []api.NamespaceOperation) tea.Cmd {
	return func() tea.Msg {
		var availableOps []api.NamespaceOperation
		for _, operation := range operations {
			if operation.Condition == nil {
				continue
			}
			if operation.Condition(cs, namespace) {
				availableOps = append(availableOps, operation)
			}
		}
		return api.AvailableOperationsMsg{availableOps}
	}
}
