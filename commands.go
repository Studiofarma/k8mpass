package main

import (
	"fmt"
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

var GetAllNamespacesOperation = NamespaceOperation{
	Name: "Get all namespaces",
	Command: func(model K8mpassModel, namespace string) tea.Cmd {
		return func() tea.Msg {
			namespacesNames, err := getNameSpaces(model.cluster.kubernetes)
			if err != nil {
				return errMsg(err)
			}
			return namespacesNamesMsg{
				body: namespacesNames,
			}
		}
	},
}

var NameSpaceSelected = NamespaceOperation{
	Name: "Namespaces selection",
	Command: func(model K8mpassModel, namespace string) tea.Cmd {
		return func() tea.Msg {
			return nameSpaceSelectedMsg{body: namespace}
		}
	},
}

var OperationSelected = NamespaceOperation{
	Name: "Operations selection",
	Command: func(model K8mpassModel, operation string) tea.Cmd {
		return func() tea.Msg {
			return operationSelectedMsg{body: operation}
		}
	},
}

var GetAllPodsOperation = NamespaceOperation{
	Name: "Get pods",
	Command: func(model K8mpassModel, namespace string) tea.Cmd {
		return func() tea.Msg {
			podsInfo, err := getPods(model.cluster.kubernetes, namespace)
			if err != nil {
				return errMsg(err)
			}
			fmt.Println(podsInfo) //just for checking it out
			return podsInfoMsg{}
		}
	},
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
