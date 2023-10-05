package pod

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"io"
	v1 "k8s.io/api/core/v1"
)

type Item struct {
	K8sPod             v1.Pod
	ExtendedProperties []Property
}

type Property struct {
	Key   string
	Value string
	Order int
}

func (n Item) FilterValue() string {
	return n.K8sPod.Name
}

type ItemDelegate struct{}

func (n ItemDelegate) Height() int {
	return 1
}

func (n ItemDelegate) Spacing() int {
	return 0
}

func (n ItemDelegate) Update(_ tea.Msg, m *list.Model) tea.Cmd {
	if len(m.Items()) == 0 {
		m.SetShowStatusBar(false)
	} else {
		m.SetShowStatusBar(true)
	}
	return nil
}

func (n ItemDelegate) Render(w io.Writer, l list.Model, _ int, listItem list.Item) {
	i, ok := listItem.(Item)
	if !ok {
		return
	}
	customProperties := ""
	for _, property := range i.ExtendedProperties {
		customProperties += customPropertiesStyle.Render(fmt.Sprintf(property.Value))
	}
	maxLength := 0
	for _, item := range l.Items() {
		maxLength = max(maxLength, len(item.FilterValue()))
	}
	_, _ = fmt.Fprintf(w, "  %s%s", podStyle(i.K8sPod.Status, maxLength).Render(i.K8sPod.Name), customProperties)
}

func FindPod(items []list.Item, search Item) int {
	var idx = -1
	for i, item := range items {
		if pod, ok := item.(Item); ok {
			if pod.K8sPod.Name == search.K8sPod.Name {
				idx = i
			}
		}
	}
	return idx
}
