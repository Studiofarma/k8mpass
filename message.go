package main

import "k8s.io/client-go/kubernetes"

type errMsg error

type clusterConnectedMsg struct {
	clientset *kubernetes.Clientset
}

type namespacesNamesMsg struct {
	body []string
}

type wakeUpReviewMsg struct {
	body string
}
