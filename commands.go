package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"k8s.io/client-go/kubernetes"
)

func connectToClusterCmd() tea.Msg {
	cs, err := getConnection()
	if err != nil {
		return errMsg(err)
	}

	return clusterConnectedMsg{cs}
}

func fetchNamespacesCmd(clientset *kubernetes.Clientset) tea.Cmd {
	return func() tea.Msg {
		k8Namespaces, err := getNamespaces(clientset)
		if err != nil {
			return errMsg(err)
		}

		var namespaces []Namespace
		for i := 0; i < len(k8Namespaces.Items); i++ {
			namespaces = append(namespaces, Namespace{
				id:   k8Namespaces.Items[i].UID,
				name: k8Namespaces.Items[i].ObjectMeta.Name,
			})
		}

		return fetchedNamespacesMsg{namespaces}
	}
}

type K8mpassCommand func(model K8mpassModel, namespace string) tea.Cmd

type NamespaceOperation struct {
	Name    string
	Command K8mpassCommand
}

var WakeUpReviewOperation = NamespaceOperation{
	Name: "Wake up review app",
	Command: func(model K8mpassModel, namespace string) tea.Cmd {
		return func() tea.Msg {
			err := wakeupReview(model.cluster.kubernetes, namespace)
			if err != nil {
				return errMsg(err)
			}
			return wakeUpReviewMsg{}
		}
	},
}
