package main

import (
	"k8s.io/client-go/kubernetes"
)

type errMsg error

type clusterConnectedMsg struct {
	clientset *kubernetes.Clientset
}

type startupMsg struct{}

type namespaceSelectedMsg struct {
	namespace string
}

type backToNamespaceSelectionMsg struct{}
