package kubernetes

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"log"
)

type PodService struct {
	Events <-chan watch.Event
}

func (s PodService) GetPods(ctx context.Context, cs *kubernetes.Clientset, namespace string) (*v1.PodList, error) {
	res, err := cs.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *PodService) Subscribe(ctx context.Context, cs *kubernetes.Clientset, resourceVersion string, namespace string) error {
	opt := metav1.ListOptions{
		ResourceVersion: resourceVersion,
	}
	watcher, err := cs.CoreV1().Pods(namespace).Watch(ctx, opt)
	if err != nil {
		return err
	}
	s.Events = watcher.ResultChan()
	return nil
}

func (s PodService) GetEvent() watch.Event {
	event := <-s.Events
	log.Println("Received pod event of type ", event.Type)
	return event
}
