package main

import "k8s.io/client-go/kubernetes"

type errMsg error

type clusterConnectedMsg struct {
	clientset *kubernetes.Clientset
}

type operationResultMsg struct {
	body string
}

type noOutputResultMsg struct {
	message string
}

type startupMsg struct{}

type namespacesRetrievedMsg struct {
	namespaces []string
}

type namespaceSelectedMsg struct {
	namespace string
}

type backToNamespaceSelectionMsg struct{}
type backToOperationSelectionMsg struct{}
