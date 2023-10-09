package pod

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/studiofarma/k8mpass/api"
	k8mpasskube "github.com/studiofarma/k8mpass/kubernetes"
	v1 "k8s.io/api/core/v1"
	"log"
)

type MessageHandler struct {
	service                      k8mpasskube.IPodService
	logs                         k8mpasskube.ILogService
	Extensions                   []api.IPodExtension
	AvailableNamespaceOperations []api.INamespaceOperation
}

func (handler *MessageHandler) NextEvent() tea.Msg {
	event := handler.service.GetPodEvent()

	switch event.Type {
	case k8mpasskube.Deleted:
		log.Printf("Deleted pod: %s ", event.Pod)
		return RemovedPodMsg{
			Pod: Item{
				K8sPod:             *event.Pod,
				ExtendedProperties: make([]Property, 0),
			},
		}
	case k8mpasskube.Added:
		pod := Item{K8sPod: *event.Pod, ExtendedProperties: make([]Property, 0)}
		pod.LoadCustomProperties(handler.Extensions...)
		log.Printf("Added pod: %s ", event.Pod.Name)
		return AddedPodMsg{
			Pod: pod,
		}
	case k8mpasskube.Modified:
		pod := Item{K8sPod: *event.Pod, ExtendedProperties: make([]Property, 0)}
		pod.LoadCustomProperties(handler.Extensions...)
		return ModifiedPodMsg{
			Pod: pod,
		}
	case k8mpasskube.Closed, k8mpasskube.Error:
		return nil
	default:
		return NextEventMsg{}
	}
}

func (handler *MessageHandler) WatchPods(ctx context.Context, resourceVersion string, namespace string) tea.Cmd {
	return func() tea.Msg {
		err := handler.service.WatchPods(ctx, namespace, resourceVersion)
		if err != nil {
			return ErrorMsg{err}
		}
		return WatchingPodsMsg{}
	}
}

func (handler *MessageHandler) GetPods(namespace string) tea.Cmd {
	res, err := handler.service.GetPods(namespace)
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
		handler.WatchPods(context.TODO(), res.ResourceVersion, namespace),
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

func (handler *MessageHandler) StopWatching() {
	handler.service.StopWatchingPods()
}

func (handler *MessageHandler) CheckConditionsThatApply(namespace string) tea.Cmd {
	return func() tea.Msg {
		var availableOps []api.INamespaceOperation
		for _, operation := range handler.AvailableNamespaceOperations {
			if operation.GetCondition() == nil {
				continue
			}
			if handler.service.RunK8mpassCondition(operation.GetCondition(), namespace) {
				availableOps = append(availableOps, operation)
			}
		}
		return api.AvailableOperationsMsg{Operations: availableOps}
	}
}

func (handler *MessageHandler) RunComand(command api.INamespaceOperation, namespace string) tea.Cmd {
	return handler.service.RunK8mpassCommand(command.GetCommand(), namespace)
}

func NewHandler(service k8mpasskube.IPodService, extensions []api.IPodExtension, ops []api.INamespaceOperation, logs k8mpasskube.ILogService) *MessageHandler {
	return &MessageHandler{
		service:                      service,
		logs:                         logs,
		Extensions:                   extensions,
		AvailableNamespaceOperations: ops,
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

func (handler *MessageHandler) FollowLogs(namespace string, pod string, maxWidth int) {
	err := handler.logs.GetLogReader(namespace, pod, maxWidth)
	if err != nil {
		return
	}
}

func (handler *MessageHandler) GetNextLogLine() tea.Cmd {
	return func() tea.Msg {
		lines, closed := handler.logs.GetNextLogs()
		if closed {
			return nil
		}
		return NextLogLineMsg{
			NextLines: lines,
		}
	}
}
