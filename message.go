package main

import (
	"k8s.io/client-go/kubernetes"
)

type errMsg error

type clusterConnectedMsg struct {
	clientset *kubernetes.Clientset
}

type operationResultMsg struct {
	body string
}

type startupMsg struct{}

type namespacesFetchedMsg struct {
	namespaces []string
}

type namespaceSelectedMsg struct {
	namespace string
}

type operationSelectedMsg struct {
	operation string
}

type backToNamespaceSelectionMsg struct{}
type backToOperationSelectionMsg struct{}

type refreshNamespacesMsg struct{}
