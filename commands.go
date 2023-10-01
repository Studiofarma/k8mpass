package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	v1 "k8s.io/api/core/v1"
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

// func fetchNamespaces() tea.Msg {
// 	time.Sleep(500 * time.Millisecond)
// 	ns, err := getNamespaces(K8sCluster.kubernetes)
// 	if err != nil {
// 		return errMsg(err)
// 	}
// 	var items []namespace.NamespaceItem
// 	sleepingInfo, err := getReviewAppsSleepingStatus()
// 	for _, n := range ns.Items {
// 		var isAwake = false
// 		if err != nil {
// 			for _, ra := range sleepingInfo {
// 				if strings.HasPrefix(ra.Metric.ExportedService, n.Name) {
// 					isAwake = ra.IsAwake() || isAwake
// 				}
// 			}
// 		}
// 		items = append(items, namespace.NamespaceItem{n, isAwake})
// 	}
// 	return namespacesRetrievedMsg{items}
// }

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
