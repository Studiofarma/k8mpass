package main

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func clusterConnect() tea.Msg {
	cs, err := getConnection()
	if err != nil {
		return errMsg(err)
	}

	err = wakeupReview(cs, "review-davinci-cgmgpharm-46374-be")
	if err != nil {
		return errMsg(err)
	}

	return clusterConnectedMsg{cs}
}

type K8mpassExecutable func(*kubernetes.Clientset, string) tea.Msg
type K8mpassCommand func(model K8mpassModel, namespace string) tea.Cmd

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
