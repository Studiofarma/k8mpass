package namespace

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	k8mpasskube "github.com/studiofarma/k8mpass/kubernetes"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type NamespaceMessageHandler struct {
	service    k8mpasskube.NamespaceService
	extensions []NamespaceExtension
}

func (nh NamespaceMessageHandler) NextEvent() tea.Msg {
	event := nh.service.GetEvent()
	item := event.Object.(*v1.Namespace)
	namespace := NamespaceItem{
		*item,
		make(map[string]string),
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

func (nh *NamespaceMessageHandler) WatchNamespaces(cs *kubernetes.Clientset, resourceVersion string) tea.Cmd {
	return tea.Sequence(
		func() tea.Msg {
			err := nh.service.Subscribe(cs, resourceVersion)
			if err != nil {
				return ErrorMsg{err}
			}
			return WatchingNamespacesMsg{}
		},
		nh.NextEvent,
	)
}

type namespaceName string

func (nh *NamespaceMessageHandler) GetNamespaces(cs *kubernetes.Clientset) tea.Cmd {
	res, err := nh.service.GetNamespaces(cs)
	return tea.Sequence(
		func() tea.Msg {
			if err != nil {
				return ErrorMsg{err}
			}
			var namespaces []NamespaceItem
			namespaceProperties := make(map[namespaceName]map[string]string)
			for _, e := range nh.extensions {
				fn := e.ExtendList
				if fn == nil {
					continue
				}
				nsToValue := fn(res.Items)
				for ns, value := range nsToValue {
					if namespaceProperties[namespaceName(ns)] == nil {
						namespaceProperties[namespaceName(ns)] = make(map[string]string)
					}
					namespaceProperties[namespaceName(ns)][e.Name] = string(value)
				}
			}
			for _, n := range res.Items {
				namespaces = append(namespaces, NamespaceItem{n, namespaceProperties[namespaceName(n.Name)]})
			}
			return NamespaceListMsg{
				Namespaces:      namespaces,
				ResourceVersion: res.ResourceVersion,
			}
		},
		nh.WatchNamespaces(cs, res.ResourceVersion),
	)
}

func NewHandler(exts ...NamespaceExtension) *NamespaceMessageHandler {
	return &NamespaceMessageHandler{
		service:    k8mpasskube.NamespaceService{},
		extensions: exts,
	}
}
