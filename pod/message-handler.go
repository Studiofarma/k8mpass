package pod

import (
	"context"
	"fmt"
	"github.com/studiofarma/k8mpass/api"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	k8mpasskube "github.com/studiofarma/k8mpass/kubernetes"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type MessageHandler struct {
	service    k8mpasskube.PodService
	Extensions []api.IPodExtension
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
		pod := Item{K8sPod: *item, ExtendedProperties: make([]Property, 0)}
		pod.LoadCustomProperties(handler.Extensions...)
		log.Printf("Added pod: %s ", item.Name)
		return AddedPodMsg{
			Pod: pod,
		}
	case watch.Modified:
		item := event.Object.(*v1.Pod)
		pod := Item{K8sPod: *item, ExtendedProperties: make([]Property, 0)}
		pod.LoadCustomProperties(handler.Extensions...)
		return ModifiedPodMsg{
			Pod: pod,
		}
	case watch.Error, "":
		log.Printf("Error event for pods")
		return nil
	default:
		log.Printf("NamespaceEvent not handled")
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
			pods := LoadExtensions(handler.Extensions, res.Items)
			return ListMsg{
				Pods:            pods,
				ResourceVersion: res.ResourceVersion,
			}
		},
		handler.WatchPods(ctx, cs, res.ResourceVersion, namespace),
	)
}

func LoadExtensions(extensions []api.IPodExtension, res []v1.Pod) []Item {
	var pods []Item
	podProperties := make(map[string][]Property)
	for idx, e := range extensions {
		fn := e.GetExtendList()
		if fn == nil {
			continue
		}
		pToValue := fn(res)
		for ns, value := range pToValue {
			if podProperties[ns] == nil {
				podProperties[ns] = make([]Property, 0)
			}
			p := Property{
				Key:   e.GetName(),
				Value: value,
				Order: idx,
			}
			podProperties[ns] = append(podProperties[ns], p)
		}
	}
	for _, n := range res {
		pods = append(pods, Item{n, podProperties[n.Name]})
	}
	return pods
}

func (n *Item) LoadCustomProperties(properties ...api.IPodExtension) {
	n.ExtendedProperties = make([]Property, 0)
	for idx, p := range properties {
		fn := p.GetExtendSingle()
		if fn == nil {
			log.Println(fmt.Sprintf("Missing extention function for %s", p.GetName()), "namespace:", n.K8sPod.Name)
			continue
		}
		value, err := fn(n.K8sPod)
		if err != nil {
			log.Println(fmt.Sprintf("Error while computing extension %s", p.GetName()), "namespace:", n.K8sPod.Name)
			continue
		}
		n.ExtendedProperties = append(n.ExtendedProperties, Property{Key: p.GetName(), Value: value, Order: idx})
	}
}

func (handler MessageHandler) StopWatching() {
	handler.service.Watcher.Stop()
}

func NewHandler(extensions ...api.IPodExtension) *MessageHandler {
	return &MessageHandler{
		service:    k8mpasskube.PodService{},
		Extensions: extensions,
	}
}

func Route(cmds ...tea.Cmd) []tea.Cmd {
	var ret []tea.Cmd
	for _, cmd := range cmds {
		if cmd == nil {
			continue
		}
		ret = append(ret, func() tea.Msg {
			return RoutedMsg{Embedded: cmd()}
		})
	}
	return ret
}
