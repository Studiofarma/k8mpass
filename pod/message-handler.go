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

func (nh MessageHandler) NextEvent() tea.Msg {
	event := nh.service.GetEvent()
	item := event.Object.(*v1.Pod)
	pod := Item{
		K8sPod:             *item,
		ExtendedProperties: make([]Property, 0),
	}
	switch event.Type {
	case watch.Deleted:
		log.Printf("Deleted pod: %s ", item.Name)
		return RemovedPodMsg{
			Pod: pod,
		}
	case watch.Added:
		log.Printf("Added pod: %s ", item.Name)
		return AddedPodMsg{
			Pod: pod,
		}
	case watch.Modified:
		log.Printf("Modified pod: %s ", item.Name)
		return AddedPodMsg{
			Pod: pod,
		}
	default:
		return NextEventMsg{}
	}
}

func (nh *MessageHandler) WatchPods(ctx context.Context, cs *kubernetes.Clientset, resourceVersion string, namespace string) tea.Cmd {
	return tea.Sequence(
		func() tea.Msg {
			err := nh.service.Subscribe(ctx, cs, resourceVersion, namespace)
			if err != nil {
				return ErrorMsg{err}
			}
			return WatchingPodsMsg{}
		},
		nh.NextEvent,
	)
}

func (nh *MessageHandler) GetPods(ctx context.Context, cs *kubernetes.Clientset, namespace string) tea.Cmd {
	res, err := nh.service.GetPods(ctx, cs, namespace)
	return tea.Sequence(
		func() tea.Msg {
			if err != nil {
				return ErrorMsg{err}
			}
			var pods []Item
			for _, n := range res.Items {
				pods = append(pods, Item{K8sPod: n})
			}
			return PodListMsg{
				Pods:            pods,
				ResourceVersion: res.ResourceVersion,
			}
		},
		nh.WatchPods(ctx, cs, res.ResourceVersion, namespace),
	)
}

func NewHandler() *MessageHandler {
	return &MessageHandler{
		service: k8mpasskube.PodService{},
	}
}
