package main

import (
	"k8s.io/client-go/kubernetes"
)

type kubernetesCluster struct {
	kubernetes *kubernetes.Clientset
}

var K8sCluster kubernetesCluster
