package namespace

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
	"github.com/studiofarma/k8mpass/api"
	"io"
	v1 "k8s.io/api/core/v1"
	"log"
)

type Item struct {
	K8sNamespace       v1.Namespace
	ExtendedProperties []Property
}

type Property struct {
	Key   string
	Value string
	Order int
}

func (n Item) FilterValue() string {
	return n.K8sNamespace.Name
}

func (n *Item) LoadCustomProperties(properties ...api.INamespaceExtension) {
	n.ExtendedProperties = make([]Property, 0)
	for idx, p := range properties {
		fn := p.GetExtendSingle()
		if fn == nil {
			log.Println(fmt.Sprintf("Missing extention function for %s", p.GetName()), "namespace:", n.K8sNamespace.Name)
			continue
		}
		value, err := fn(n.K8sNamespace)
		if err != nil {
			log.Println(fmt.Sprintf("Error while computing extension %s", p.GetName()), "namespace:", n.K8sNamespace.Name)
			continue
		}
		n.ExtendedProperties = append(n.ExtendedProperties, Property{Key: p.GetName(), Value: value, Order: idx})
	}
}

type ItemDelegate struct{}

func (n ItemDelegate) Height() int {
	return 1
}

func (n ItemDelegate) Spacing() int {
	return 0
}

func (n ItemDelegate) Update(tea.Msg, *list.Model) tea.Cmd {
	return nil
}

func (n ItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Item)
	if !ok {
		return
	}

	maxLength := 0
	propertyWidth := 12
	for _, item := range m.VisibleItems() {
		maxLength = max(maxLength, len(item.FilterValue()))
	}

	namespace := i.K8sNamespace.Name
	propertiesStyle := customPropertiesStyle.Copy()
	style := unselectedItemStyle
	if index == m.Index() {
		style = selectedItemStyle
		style = style.Background(lipgloss.Color("#444852"))
		propertiesStyle = propertiesStyle.Background(lipgloss.Color("#444852"))
	}

	customProperties := ""
	for _, property := range i.ExtendedProperties {
		prop := lipgloss.PlaceHorizontal(propertyWidth, lipgloss.Left, ellipsis(property.Value, propertyWidth))
		customProperties += propertiesStyle.Render(prop)
	}

	if i.K8sNamespace.Status.Phase == v1.NamespaceTerminating {
		style = terminatingNamespace
	}

	_, _ = fmt.Fprint(w, style.Width(maxLength+3).Render(namespace)+customProperties)
}

func FindNamespace(items []list.Item, search Item) int {
	var idx = -1
	for i, item := range items {
		if ns, ok := item.(Item); ok {
			if ns.K8sNamespace.Name == search.K8sNamespace.Name {
				idx = i
			}
		}
	}
	return idx
}

func ellipsis(s string, length int) string {
	if len(s) > length {
		return truncate.StringWithTail(s, uint(length), "..")
	}
	return s
}
