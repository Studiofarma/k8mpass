package main

import (
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type errMsg error

type clusterConnectedMsg struct {
	clientset *kubernetes.Clientset
}

type operationResultMsg struct {
	body string
}

type noOutputResultMsg struct {
	success bool
	message string
}

type startupMsg struct{}

type namespaceSelectedMsg struct {
	namespace string
}

type backToNamespaceSelectionMsg struct{}
type backToOperationSelectionMsg struct{}

type watchNamespaceMsg struct {
	channel <-chan watch.Event
}

type nextEventMsg struct{}
