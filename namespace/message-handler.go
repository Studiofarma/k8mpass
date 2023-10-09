package namespace

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/studiofarma/k8mpass/api"
	k8mpasskube "github.com/studiofarma/k8mpass/kubernetes"
	v1 "k8s.io/api/core/v1"
	"log"
	"time"
)

type MessageHandler struct {
	service    k8mpasskube.INamespaceService
	Extensions []api.INamespaceExtension
}

func (nh MessageHandler) NextEvent() tea.Msg {
	event := nh.service.GetNamespaceEvent()
	switch event.Type {
	case k8mpasskube.Deleted:
		log.Printf("Deleted namespace: %s ", event.Namespace.Name)
		namespace := Item{
			K8sNamespace: *event.Namespace,
		}
		return RemovedMsg{
			Namespace: namespace,
		}
	case k8mpasskube.Added:
		namespace := Item{
			K8sNamespace: *event.Namespace,
		}
		namespace.LoadCustomProperties(nh.Extensions...)
		log.Printf("Added namespace: %s ", event.Namespace.Name)
		return AddedMsg{
			Namespace: namespace,
		}
	case k8mpasskube.Modified:
		namespace := Item{
			K8sNamespace: *event.Namespace,
		}
		namespace.LoadCustomProperties(nh.Extensions...)
		log.Printf("Modified namespace: %s ", event.Namespace.Name)
		return ModifiedMsg{
			Namespace: namespace,
		}
	case k8mpasskube.Closed, k8mpasskube.Error:
		return nil
	default:
		return NextEventMsg{}
	}
}

func (nh *MessageHandler) WatchNamespaces(resourceVersion string) tea.Cmd {
	return func() tea.Msg {
		err := nh.service.WatchNamespaces(resourceVersion)
		if err != nil {
			return ErrorMsg{err}
		}
		return WatchingMsg{}
	}
}

func (nh *MessageHandler) GetNamespaces() tea.Cmd {
	res, err := nh.service.GetNamespaces()
	return tea.Sequence(
		func() tea.Msg {
			if err != nil {
				return ErrorMsg{err}
			}
			namespaces := LoadExtensions(nh.Extensions, res.Items)
			return ListMsg{
				Namespaces: namespaces,
			}
		},
		nh.WatchNamespaces(res.ResourceVersion),
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
	namespaceProperties := make(map[string][]Property)
	for idx, e := range extensions {
		fn := e.GetExtendList()
		if fn == nil {
			continue
		}
		nsToValue := fn(res)
		for ns, value := range nsToValue {
			if namespaceProperties[ns] == nil {
				namespaceProperties[ns] = make([]Property, 0)
			}
			p := Property{
				Key:   e.GetName(),
				Value: value,
				Order: idx,
			}
			namespaceProperties[ns] = append(namespaceProperties[ns], p)
		}
	}
	for _, n := range res {
		namespaces = append(namespaces, Item{n, namespaceProperties[n.Name]})
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
			if namespaceProperties[ns] == nil {
				namespaceProperties[ns] = make([]Property, 0)
			}
			p := Property{
				Key:   e.GetName(),
				Value: value,
				Order: idx,
			}
			namespaceProperties[ns] = append(namespaceProperties[ns], p)
		}
	}
	return namespaceProperties
}

func NewHandler(service k8mpasskube.INamespaceService, exts ...api.INamespaceExtension) *MessageHandler {
	return &MessageHandler{
		service:    service,
		Extensions: exts,
	}
}

func Route(cmds []tea.Cmd) []tea.Cmd {
	var ret []tea.Cmd
	for _, cmd := range cmds {
		fn := cmd //This is needed to avoid passing the reference to the iteration element, in which case all the functions would point to the last cmd in the iteration
		if fn == nil {
			continue
		}
		ret = append(ret, func() tea.Msg {
			result := fn()
			if result == nil {
				return nil
			} else {
				return RoutedMsg{Embedded: result}
			}
		})
	}
	return ret
}
