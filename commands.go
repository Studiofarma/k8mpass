package main

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
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
	var nsNames []string
	for _, n := range ns.Items {
		nsNames = append(nsNames, n.Name)
	}
	return namespacesRetrievedMsg{nsNames}
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
				return errMsg(err)
			}
			return operationResultMsg{"  We woke it up!"}
		}
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
			time.Sleep(200 * time.Millisecond)
			ingresses, err := clientset.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})

			if err != nil {
				return errMsg(err)
			}

			var dbmsUrl string

			for _, i := range ingresses.Items {
				host := i.Spec.Rules[0].Host
				if strings.HasPrefix(host, "dbms") {
					dbmsUrl = host
				}
			}

			Openbrowser("https://" + dbmsUrl)

			return operationResultMsg{body: "DBeaver is better 🦦"}
		}
	},
}

var OpenApplicationOperation = NamespaceOperation{
	Name: "Open application in browser",
	Command: func(clientset *kubernetes.Clientset, namespace string) tea.Cmd {
		return func() tea.Msg {
			time.Sleep(200 * time.Millisecond)
			ingresses, err := clientset.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})

			if err != nil {
				return errMsg(err)
			}

			var dbmsUrl string

			for _, i := range ingresses.Items {
				if strings.HasPrefix(i.Name, "g3pharmacy") {
					dbmsUrl = i.Spec.Rules[0].Host
				}
			}

			Openbrowser("https://" + dbmsUrl)

			return operationResultMsg{body: "App is ready. Who is Bruno Marotta?"}
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
