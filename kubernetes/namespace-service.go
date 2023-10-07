package kubernetes

import (
	"context"
	"log"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type INamespaceService interface {
	GetNamespaces() (*v1.NamespaceList, error)
	GetEvent() NamespaceEvent
	Watch(resourceVersion string) error
}

func (c *Cluster) GetNamespaces() (*v1.NamespaceList, error) {
	res, err := c.cs.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Cluster) GetEvent() NamespaceEvent {
	e, closed := <-c.namespaceWatch.ResultChan()
	if closed {
		return NamespaceEvent{
			Type: Closed,
		}
	}
	switch e.Type {
	case watch.Added:
		return NamespaceEvent{
			Type:      Added,
			Namespace: e.Object.(*v1.Namespace),
		}
	case watch.Modified:
		return NamespaceEvent{
			Type:      Modified,
			Namespace: e.Object.(*v1.Namespace),
		}
	case watch.Deleted:
		return NamespaceEvent{
			Type:      Deleted,
			Namespace: e.Object.(*v1.Namespace),
		}
	case watch.Bookmark:
		return NamespaceEvent{
			Type: Unhandled,
		}
	case watch.Error:
		return NamespaceEvent{
			Type: Error,
		}
	default:
		return NamespaceEvent{}
	}
}

func (c *Cluster) Watch(resourceVersion string) error {
	opt := metav1.ListOptions{
		ResourceVersion: resourceVersion,
	}
	watcher, err := c.cs.CoreV1().Namespaces().Watch(context.TODO(), opt)
	if err != nil {
		return err
	}
	c.namespaceWatch = watcher
	return nil
}

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
