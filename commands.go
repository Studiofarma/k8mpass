package main

import (
	"context"
	"fmt"
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

	return clusterConnectedMsg{cs}
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
			return wakeUpReviewMsg{"  We woke it up!"}
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
			return wakeUpReviewMsg{body: s}
		}
	},
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
