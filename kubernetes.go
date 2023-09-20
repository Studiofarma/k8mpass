package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"path/filepath"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type kubernetesCluster struct {
	kubernetes *kubernetes.Clientset
}

type namespacePodsInfo struct {
	status            string
	numberRunningPods int
	numberErrorPods   int
	podsInfo          []podInfo
}
type podInfo struct {
	name             string
	status           string
	startTime        string
	numberOfRestarts int
}

func (n *namespacePodsInfo) calculateNamespaceStatus() {
	n.numberErrorPods = 0
	n.numberRunningPods = 0
	n.status = "Ok"
	for _, value := range n.podsInfo {
		if value.status == "Running" {
			n.numberRunningPods++
		} else {
			n.numberErrorPods++
		}
	}
	if n.numberRunningPods < 5 {
		n.status = "To wake up"
	}
	if n.numberErrorPods > 0 {
		n.status = "Pods with error"
	}
}

func getConnection() (*kubernetes.Clientset, error) {
	args := os.Args
	var kubeConfigPath = ""
	if len(args) > 1 {
		kubeConfigPath = args[1]
	} else {
		kubeConfigPath = defaultKubeConfigFilePath()
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
	return filepath.Join(userHomeDir, ".kube", "config")
}

func getNameSpaces(clientset *kubernetes.Clientset) ([]string, error) {
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})

	var namespacesNames []string
	for _, value := range namespaces.Items {
		namespacesNames = append(namespacesNames, value.ObjectMeta.Name)
	}
	if err != nil {
		return []string{}, err
	}

	return namespacesNames, err
}

func getPods(clientset *kubernetes.Clientset, namespace string) ([]podInfo, error) {
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})

	var podsInfo []podInfo
	for _, value := range pods.Items {
		restartCount := -1
		if len(value.Status.ContainerStatuses) > 0 {
			restartCount = int(value.Status.ContainerStatuses[0].RestartCount)
		}
		podsInfo = append(podsInfo, podInfo{
			name:             value.ObjectMeta.Name,
			status:           string(value.Status.Phase),
			startTime:        value.Status.StartTime.String(),
			numberOfRestarts: restartCount,
		})
	}

	if err != nil {
		return []podInfo{}, err
	}

	return podsInfo, err
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
