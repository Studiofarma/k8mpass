package pod

import (
	"context"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	k8mpasskube "github.com/studiofarma/k8mpass/kubernetes"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type MessageHandler struct {
	service k8mpasskube.PodService
}

func (handler MessageHandler) NextEvent() tea.Msg {
	event := handler.service.GetEvent()

	switch event.Type {
	case watch.Deleted:
		item := event.Object.(*v1.Pod)
		log.Printf("Deleted pod: %s ", item.Name)
		return RemovedPodMsg{
			Pod: Item{
				K8sPod:             *item,
				ExtendedProperties: make([]Property, 0),
			},
		}
	case watch.Added:
		item := event.Object.(*v1.Pod)
		log.Printf("Added pod: %s ", item.Name)
		return AddedPodMsg{
			Pod: Item{
				K8sPod: *item,
			},
		}
	case watch.Modified:
		item := event.Object.(*v1.Pod)
		return ModifiedPodMsg{
			Pod: Item{
				K8sPod: *item,
			},
		}
	case watch.Error, "":
		log.Printf("Error event for pods")
		return nil
	default:
		log.Printf("Event not handled")
		return NextEventMsg{}
	}
}

func (handler *MessageHandler) WatchPods(ctx context.Context, cs *kubernetes.Clientset, resourceVersion string, namespace string) tea.Cmd {
	return func() tea.Msg {
		err := handler.service.Subscribe(ctx, cs, resourceVersion, namespace)
		if err != nil {
			return ErrorMsg{err}
		}
		return WatchingPodsMsg{}
	}
}

func (handler *MessageHandler) GetPods(ctx context.Context, cs *kubernetes.Clientset, namespace string) tea.Cmd {
	res, err := handler.service.GetPods(ctx, cs, namespace)
	return tea.Sequence(
		func() tea.Msg {
			if err != nil {
				return ErrorMsg{err}
			}
			var pods []Item
			for _, n := range res.Items {
				pods = append(pods, Item{K8sPod: n})
			}
			return ListMsg{
				Pods:            pods,
				ResourceVersion: res.ResourceVersion,
			}
		},
		handler.WatchPods(ctx, cs, res.ResourceVersion, namespace),
	)
}

func (handler MessageHandler) StopWatching() {
	handler.service.Watcher.Stop()
}

func NewHandler() *MessageHandler {
	return &MessageHandler{
		service: k8mpasskube.PodService{},
	}
}
