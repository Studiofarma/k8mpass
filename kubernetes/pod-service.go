package kubernetes

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/studiofarma/k8mpass/api"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type IPodService interface {
	GetPods(namespace string) (*v1.PodList, error)
	GetPodEvent() PodEvent
	WatchPods(ctx context.Context, namespace string, resourceVersion string) error
	StopWatchingPods()
	RunK8mpassCondition(fn api.K8mpassCondition, namespace string) bool
	RunK8mpassCommand(fn api.K8mpassCommand, namespace string) tea.Cmd
}

func (c *Cluster) GetPods(namespace string) (*v1.PodList, error) {
	res, err := c.cs.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Cluster) GetPodEvent() PodEvent {
	e, closed := <-c.podWatch.ResultChan()
	if !closed {
		return PodEvent{
			Type: Closed,
		}
	}
	switch e.Type {
	case watch.Added:
		return PodEvent{
			Type: Added,
			Pod:  e.Object.(*v1.Pod),
		}
	case watch.Modified:
		return PodEvent{
			Type: Modified,
			Pod:  e.Object.(*v1.Pod),
		}
	case watch.Deleted:
		return PodEvent{
			Type: Deleted,
			Pod:  e.Object.(*v1.Pod),
		}
	case watch.Bookmark:
		return PodEvent{
			Type: Unhandled,
		}
	case watch.Error, "":
		return PodEvent{
			Type: Error,
		}
	default:
		return PodEvent{}
	}

}

func (c *Cluster) WatchPods(ctx context.Context, namespace string, resourceVersion string) error {
	opt := metav1.ListOptions{
		ResourceVersion: resourceVersion,
	}
	watcher, err := c.cs.CoreV1().Pods(namespace).Watch(ctx, opt)
	if err != nil {
		return err
	}

	c.podWatch = watcher
	return nil
}

func (c *Cluster) StopWatchingPods() {
	c.podWatch.Stop()
}

func (c *Cluster) RunK8mpassCondition(fn api.K8mpassCondition, namespace string) bool {
	return fn(c.cs, namespace)
}

func (c *Cluster) RunK8mpassCommand(fn api.K8mpassCommand, namespace string) tea.Cmd {
	return fn(c.cs, namespace)
}
