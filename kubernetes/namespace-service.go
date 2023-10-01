package kubernetes

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type NamespaceService struct {
	Events <-chan watch.Event
}

func (s *NamespaceService) Subscribe(cs *kubernetes.Clientset) error {
	opt := metav1.ListOptions{}
	watcher, err := cs.CoreV1().Namespaces().Watch(context.TODO(), opt)
	if err != nil {
		return err
	}
	s.Events = watcher.ResultChan()
	return nil
}

func (s NamespaceService) GetEvent() watch.Event {
	return <-s.Events
}

func (s NamespaceService) GetEvents() []watch.Event {
	events := make([]watch.Event, 0)
poll:
	for {
		select {
		case e := <-s.Events:
			events = append(events, e)
		default:
			break poll
		}
	}
	return events
}
