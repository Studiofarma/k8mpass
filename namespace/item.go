package namespace

import (
	"fmt"
	"github.com/studiofarma/k8mpass/api"
	"io"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	v1 "k8s.io/api/core/v1"
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

func (n *Item) LoadCustomProperties(properties ...api.IExtension) {
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
		n.ExtendedProperties = append(n.ExtendedProperties, Property{Key: p.GetName(), Value: string(value), Order: idx})
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

	namespace := i.K8sNamespace.Name
	customProperties := ""
	for _, property := range i.ExtendedProperties {
		customProperties += customPropertiesStyle.Render(fmt.Sprintf(property.Value))
	}
	fn := unselectedItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render(strings.Join(s, " "))
		}
	}

	_, _ = fmt.Fprint(w, fn(namespace)+customProperties)
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
