package kubernetes

import (
	"context"
	"log"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type NamespaceService struct {
	Events <-chan watch.Event
}

func (s *NamespaceService) Subscribe(cs *kubernetes.Clientset, resourceVersion string) error {
	opt := metav1.ListOptions{
		ResourceVersion: resourceVersion,
	}
	watcher, err := cs.CoreV1().Namespaces().Watch(context.TODO(), opt)
	if err != nil {
		return err
	}
	s.Events = watcher.ResultChan()
	return nil
}

func (s NamespaceService) GetEvent() watch.Event {
	event := <-s.Events
	log.Println("Received namespace event of type ", event.Type)
	return event
}

func (s NamespaceService) GetNamespaces(cs *kubernetes.Clientset) (*v1.NamespaceList, error) {
	res, err := cs.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return res, nil
}
