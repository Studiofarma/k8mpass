package namespace

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/studiofarma/k8mpass/api"
	k8mpasskube "github.com/studiofarma/k8mpass/kubernetes"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"log"
	"time"
)

type MessageHandler struct {
	service    k8mpasskube.NamespaceService
	Extensions []api.INamespaceExtension
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
		return RemovedMsg{
			Namespace: namespace,
		}
	case watch.Added:
		namespace.LoadCustomProperties(nh.Extensions...)
		log.Printf("Added namespace: %s ", item.Name)
		return AddedMsg{
			Namespace: namespace,
		}
	case watch.Modified:
		namespace.LoadCustomProperties(nh.Extensions...)
		log.Printf("Modified namespace: %s ", item.Name)
		return ModifiedMsg{
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
		return WatchingMsg{}
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
			namespaces := LoadExtensions(nh.Extensions, res.Items)
			return ListMsg{
				Namespaces:      namespaces,
				ResourceVersion: res.ResourceVersion,
			}
		},
		nh.WatchNamespaces(cs, res.ResourceVersion),
	)
}

func (nh *MessageHandler) ReloadExtensions(namespaces []Item) tea.Cmd {
	return func() tea.Msg {
		extensions := GetReloadedExtensions(nh.Extensions, namespaces)
		return ReloadExtensionsMsg{
			Properties: extensions,
		}
	}
}

func Refresh() tea.Cmd {
	return tea.Tick(time.Minute, func(t time.Time) tea.Msg {
		return ReloadTick{}
	})
}

func LoadExtensions(extensions []api.INamespaceExtension, res []v1.Namespace) []Item {
	var namespaces []Item
	namespaceProperties := make(map[namespaceName][]Property)
	for idx, e := range extensions {
		fn := e.GetExtendList()
		if fn == nil {
			continue
		}
		nsToValue := fn(res)
		for ns, value := range nsToValue {
			if namespaceProperties[namespaceName(ns)] == nil {
				namespaceProperties[namespaceName(ns)] = make([]Property, 0)
			}
			p := Property{
				Key:   e.GetName(),
				Value: string(value),
				Order: idx,
			}
			namespaceProperties[namespaceName(ns)] = append(namespaceProperties[namespaceName(ns)], p)
		}
	}
	for _, n := range res {
		namespaces = append(namespaces, Item{n, namespaceProperties[namespaceName(n.Name)]})
	}
	return namespaces
}

func GetReloadedExtensions(extensions []api.INamespaceExtension, res []Item) map[string][]Property {
	namespaceProperties := make(map[string][]Property)
	for idx, e := range extensions {
		fn := e.GetExtendList()
		if fn == nil {
			continue
		}
		var v1Namespaces []v1.Namespace
		for _, v := range res {
			v1Namespaces = append(v1Namespaces, v.K8sNamespace)
		}
		nsToValue := fn(v1Namespaces)
		for ns, value := range nsToValue {
			if namespaceProperties[string(ns)] == nil {
				namespaceProperties[string(ns)] = make([]Property, 0)
			}
			p := Property{
				Key:   e.GetName(),
				Value: string(value),
				Order: idx,
			}
			namespaceProperties[string(ns)] = append(namespaceProperties[string(ns)], p)
		}
	}
	return namespaceProperties
}

func NewHandler(exts ...api.INamespaceExtension) *MessageHandler {
	return &MessageHandler{
		service:    k8mpasskube.NamespaceService{},
		Extensions: exts,
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
