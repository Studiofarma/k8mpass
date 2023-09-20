package main

import "k8s.io/client-go/kubernetes"

type errMsg error

type clusterConnectedMsg struct {
	clientset *kubernetes.Clientset
}

type nameSpaceSelectedMsg struct{ body string }
type operationSelectedMsg struct{ body string }

type namespacesNamesMsg struct {
	body []string
}

type wakeUpReviewMsg struct {
	body string
}
