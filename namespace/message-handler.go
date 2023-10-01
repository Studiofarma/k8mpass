package namespace

import (
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
	extendedProperties := make(map[string]string)

	for _, ext := range nh.extensions {
		extendedProperties[ext.Name] = ext.ExtendSingle(*item)
	}

	switch event.Type {
	case watch.Deleted:
		return RemovedNamespaceMsg{
			Namespace: NamespaceItem{K8sNamespace: *item},
		}
	case watch.Added:
		return AddedNamespaceMsg{
			Namespace: NamespaceItem{K8sNamespace: *item, ExtendedProperties: extendedProperties},
		}
	default:
		return NextEventMsg{}
	}
}

func (nh *NamespaceMessageHandler) WatchNamespaces(cs *kubernetes.Clientset, resourceversion string) tea.Cmd {
	return tea.Sequence(
		func() tea.Msg {
			err := nh.service.Subscribe(cs, resourceversion)
			if err != nil {
				return ErrorMsg{err}
			}
			return WatchingNamespacesMsg{}
		},
		nh.NextEvent,
	)
}

func (nh *NamespaceMessageHandler) GetNamespaces(cs *kubernetes.Clientset) tea.Cmd {
	res, err := nh.service.GetNamespaces(cs)
	return tea.Sequence(
		func() tea.Msg {
			if err != nil {
				return ErrorMsg{err}
			}
			var namespaces []NamespaceItem
			for _, n := range res.Items {
				extendedProperties := make(map[string]string)

				for _, ext := range nh.extensions {
					extendedProperties[ext.Name] = ext.ExtendSingle(n)
				}
				namespaces = append(namespaces, NamespaceItem{n, extendedProperties})
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
