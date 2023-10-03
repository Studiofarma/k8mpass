package namespace

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	k8mpasskube "github.com/studiofarma/k8mpass/kubernetes"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type MessageHandler struct {
	service    k8mpasskube.NamespaceService
	extensions []Extension
}

func (nh MessageHandler) NextEvent() tea.Msg {
	event := nh.service.GetEvent()
	item, ok := event.Object.(*v1.Namespace)
	if !ok {
		return NextEventMsg{}
	}
	namespace := Item{
		K8sNamespace:       *item,
		ExtendedProperties: make([]Property, 0),
	}
	switch event.Type {
	case watch.Deleted:
		log.Printf("Deleted namespace: %s ", item.Name)
		return RemovedNamespaceMsg{
			Namespace: namespace,
		}
	case watch.Added:
		namespace.LoadCustomProperties(nh.extensions...)
		log.Printf("Added namespace: %s ", item.Name)
		return AddedNamespaceMsg{
			Namespace: namespace,
		}
	default:
		return NextEventMsg{}
	}
}

func (nh *MessageHandler) WatchNamespaces(cs *kubernetes.Clientset, resourceVersion string) tea.Cmd {
	return func() tea.Msg {
		err := nh.service.Subscribe(cs, resourceVersion)
		if err != nil {
			return ErrorMsg{err}
		}
		return WatchingNamespacesMsg{}
	}
}

type namespaceName string

func (nh *MessageHandler) GetNamespaces(cs *kubernetes.Clientset) tea.Cmd {
	res, err := nh.service.GetNamespaces(cs)
	return tea.Sequence(
		func() tea.Msg {
			if err != nil {
				return ErrorMsg{err}
			}
			var namespaces []Item
			namespaceProperties := make(map[namespaceName][]Property)
			for idx, e := range nh.extensions {
				fn := e.ExtendList
				if fn == nil {
					continue
				}
				nsToValue := fn(res.Items)
				for ns, value := range nsToValue {
					if namespaceProperties[namespaceName(ns)] == nil {
						namespaceProperties[namespaceName(ns)] = make([]Property, 0)
					}
					p := Property{
						Key:   e.Name,
						Value: string(value),
						Order: idx,
					}
					namespaceProperties[namespaceName(ns)] = append(namespaceProperties[namespaceName(ns)], p)
				}
			}
			for _, n := range res.Items {
				namespaces = append(namespaces, Item{n, namespaceProperties[namespaceName(n.Name)]})
			}
			return NamespaceListMsg{
				Namespaces:      namespaces,
				ResourceVersion: res.ResourceVersion,
			}
		},
		nh.WatchNamespaces(cs, res.ResourceVersion),
	)
}

func NewHandler(exts ...Extension) *MessageHandler {
	return &MessageHandler{
		service:    k8mpasskube.NamespaceService{},
		extensions: exts,
	}
}
