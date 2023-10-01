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

func (nh *NamespaceMessageHandler) WatchNamespaces(cs *kubernetes.Clientset) tea.Cmd {
	return tea.Sequence(
		func() tea.Msg {
			err := nh.service.Subscribe(cs)
			if err != nil {
				return ErrorMsg{err}
			}
			return WatchingNamespacesMsg{}
		},
		nh.NextEvent,
	)
}

func NewHandler(exts ...NamespaceExtension) *NamespaceMessageHandler {
	return &NamespaceMessageHandler{
		service:    k8mpasskube.NamespaceService{},
		extensions: exts,
	}
}
