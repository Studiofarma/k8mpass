package main

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/studiofarma/k8mpass/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
)

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
