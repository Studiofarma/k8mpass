package main

import (
	"k8s.io/client-go/kubernetes"
)

type errMsg error

type clusterConnectedMsg struct {
	clientset *kubernetes.Clientset
}

type fetchedNamespacesMsg struct {
	namespaces []Namespace
}

type wakeUpReviewMsg struct {
	body string
}
