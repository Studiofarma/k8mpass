package namespace

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	v1 "k8s.io/api/core/v1"
)

type NamespaceItem struct {
	K8sNamespace     v1.Namespace
	IsAwake          bool
	CustomProperties map[string]string
}

func (n NamespaceItem) FilterValue() string {
	return n.K8sNamespace.Name
}

func (n *NamespaceItem) LoadCustomProperties(properties ...NamespaceCustomProperty) {
	for _, p := range properties {
		n.CustomProperties[p.Name] = p.Func(&n.K8sNamespace)
	}
}

type NamespaceItemDelegate struct{}

func (n NamespaceItemDelegate) Height() int {
	return 1
}

func (n NamespaceItemDelegate) Spacing() int {
	return 0
}

func (n NamespaceItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (n NamespaceItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(NamespaceItem)
	if !ok {
		return
	}

	namespace := i.K8sNamespace.Name
	customProperties := ""
	for _, property := range i.CustomProperties {
		customProperties += customPropertiesStyle.Render(property)
	}
	fn := unselectedItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(namespace)+customProperties)
}
