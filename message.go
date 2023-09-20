package main

import "k8s.io/client-go/kubernetes"

type errMsg error

type clusterConnectedMsg struct {
	clientset *kubernetes.Clientset
}

type wakeUpReviewMsg struct {
	body string
}

type namespacesRetrievedMsg struct {
	namespaces []string
}

type namespaceSelectedMsg struct {
	namespace string
}
