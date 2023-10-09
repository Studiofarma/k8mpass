package api

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/studiofarma/k8mpass/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
)

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
