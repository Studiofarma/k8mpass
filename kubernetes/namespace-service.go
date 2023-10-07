package kubernetes

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type INamespaceService interface {
	GetNamespaces() (*v1.NamespaceList, error)
	GetNamespaceEvent() NamespaceEvent
	WatchNamespaces(resourceVersion string) error
}

func (c *Cluster) GetNamespaces() (*v1.NamespaceList, error) {
	res, err := c.cs.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Cluster) GetNamespaceEvent() NamespaceEvent {
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

func (c *Cluster) WatchNamespaces(resourceVersion string) error {
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
