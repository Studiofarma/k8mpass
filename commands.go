package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/studiofarma/k8mpass/api"
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

func CheckConditionsThatApply(cs *kubernetes.Clientset, namespace string, operations []api.INamespaceOperation) tea.Cmd {
	return func() tea.Msg {
		var availableOps []api.INamespaceOperation
		for _, operation := range operations {
			if operation.GetCondition() == nil {
				continue
			}
			if operation.GetCondition()(cs, namespace) {
				availableOps = append(availableOps, operation)
			}
		}
		return api.AvailableOperationsMsg{Operations: availableOps}
	}
}
