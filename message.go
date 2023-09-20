package main

import (
	"container/list"
	"k8s.io/client-go/kubernetes"
)

type errMsg error

type clusterConnectedMsg struct {
	clientset *kubernetes.Clientset
}

type fetchNamespacesMsg struct {
	namespacesList *list.List
}

type wakeUpReviewMsg struct {
	body string
}
