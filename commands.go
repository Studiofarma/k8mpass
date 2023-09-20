package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

func clusterConnect() tea.Msg {
	cs, err := getConnection()
	if err != nil {
		return errMsg(err)
	}

	return clusterConnectedMsg{cs}
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
			err := wakeupReview(model.cluster.GetClientset(), namespace)
			if err != nil {
				return errMsg(err)
			}
			return wakeUpReviewMsg{}
		}
	},
}

var NamespacesListOperation = NamespaceOperation{
	Name: "Wake up review app",
	Command: func(model K8mpassModel, namespace string) tea.Cmd {
		return func() tea.Msg {
			err := wakeupReview(model.cluster.GetClientset(), namespace)
			if err != nil {
				return errMsg(err)
			}
			return wakeUpReviewMsg{}
		}
	},
}
