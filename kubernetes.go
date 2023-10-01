package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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

func wakeupReview(clientset *kubernetes.Clientset, namespace string) error {
	cronjobs := clientset.BatchV1().CronJobs(namespace)
	cronjob, err := cronjobs.Get(context.TODO(), "scale-to-zero-wakeup", metav1.GetOptions{})
	if err != nil {
		return err
	}

	newUUID, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	jobSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("k8mpass-wakeup-%s", newUUID.String()),
			Namespace: namespace,
		},
		Spec: cronjob.Spec.JobTemplate.Spec,
	}
	jobs := clientset.BatchV1().Jobs(namespace)

	_, err = jobs.Create(context.TODO(), jobSpec, metav1.CreateOptions{})

	return err
}
