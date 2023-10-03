package main

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type kubernetesCluster struct {
	kubernetes *kubernetes.Clientset
}

var K8sCluster kubernetesCluster

func getConnection() (*kubernetes.Clientset, error) {
	args := os.Args
	var kubeConfigPath = defaultKubeConfigFilePath()
	if len(args) > 1 {
		kubeConfigPath = args[1]
	}
	// To add a minimum spinner time
	sleep := time.NewTimer(time.Millisecond * 500).C

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, err
	}
	k8s, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}
	res := k8s.RESTClient().Get().AbsPath("/healthz").Do(context.TODO())
	if res.Error() != nil {
		return nil, errors.New(res.Error().Error())
	}
	<-sleep
	return k8s, nil
}

func defaultKubeConfigFilePath() string {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic("error getting user home dir: %v\n")
	}
	return filepath.Join(userHomeDir, ".kube", "config")
}
