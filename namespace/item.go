package namespace

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	v1 "k8s.io/api/core/v1"
)

type NamespaceItem struct {
	K8sNamespace       v1.Namespace
	ExtendedProperties []Property
}

type Property struct {
	Key   string
	Value string
	Order int
}

func (n NamespaceItem) FilterValue() string {
	return n.K8sNamespace.Name
}

func (n *NamespaceItem) LoadCustomProperties(properties ...NamespaceExtension) {
	for idx, p := range properties {
		fn := p.ExtendSingle
		if fn == nil {
			log.Println(fmt.Sprintf("Missing extention function for %s", p.Name), "namespace:", n.K8sNamespace.Name)
			continue
		}
		value, err := fn(n.K8sNamespace)
		if err != nil {
			log.Println(fmt.Sprintf("Error while computing extension %s", p.Name), "namespace:", n.K8sNamespace.Name)
			continue
		}
		n.ExtendedProperties = append(n.ExtendedProperties, Property{Key: p.Name, Value: string(value), Order: idx})
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
	for _, property := range i.ExtendedProperties {
		customProperties += customPropertiesStyle.Render(fmt.Sprintf(property.Value))
	}
	fn := unselectedItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(namespace)+customProperties)
}
