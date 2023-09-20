package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"os"
	"path/filepath"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type kubernetesCluster struct {
	kubernetes *kubernetes.Clientset
}

func getConnection() (*kubernetes.Clientset, error) {
	args := os.Args
	var kubeConfigPath string = defaultKubeConfigFilePath()
	if len(args) > 1 {
		kubeConfigPath = args[1]
	}
	// To add a minimim spinner time
	sleep := make(chan string)
	go func(c chan string) {
		time.Sleep(1000 * time.Millisecond)
		close(c)
	}(sleep)

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		fmt.Printf("Error getting kubernetes config: %v\n", err)
		return nil, err
	}
	<-sleep
	k8s, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}
	res := k8s.RESTClient().Get().AbsPath("/healthz").Do(context.TODO())
	if res.Error() != nil {
		return nil, errors.New(res.Error().Error())
	}
	return k8s, nil
}

func defaultKubeConfigFilePath() string {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic("error getting user home dir: %v\n")
	}
	return filepath.Join(userHomeDir, ".kube", "config_reviewapp_next.dist")
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

func getConfig() *rest.Config {
	config, err := clientcmd.BuildConfigFromFlags("", defaultKubeConfigFilePath())
	if err != nil {
		panic(err.Error())
	}
	return config
}

func createClientSet(config *rest.Config) *kubernetes.Clientset {
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil
	}
	return clientSet
}

func getPods(clientset *kubernetes.Clientset, nameSpace string) *v1.PodList {
	pods, err := clientset.CoreV1().Pods(nameSpace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	return pods
}
