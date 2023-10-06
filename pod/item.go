package pod

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
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

type ItemDelegate struct {
	IsFocused bool
}

func (n ItemDelegate) Height() int {
	return 1
}

func (n ItemDelegate) Spacing() int {
	return 0
}

func (n ItemDelegate) Update(_ tea.Msg, m *list.Model) tea.Cmd {
	selectedItem := m.SelectedItem()
	if selectedItem != nil {
		m.Title = selectedItem.FilterValue()
	}
	if len(m.Items()) == 0 {
		m.SetShowStatusBar(false)
	} else {
		m.SetShowStatusBar(true)
	}
	return nil
}

func (n ItemDelegate) Render(w io.Writer, l list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Item)
	if !ok {
		return
	}
	propertyWidth := uint(15)
	style := podStyle(i.K8sPod.Status, int(maxPodNameLength))
	propertiesStyle := customPropertiesStyle.Copy()
	if n.IsFocused && index == l.Index() {
		style = style.Background(lipgloss.Color("#444852"))
		propertiesStyle = propertiesStyle.Background(lipgloss.Color("#444852"))
		propertyWidth = uint(20)
	}
	customProperties := ""
	for _, property := range i.ExtendedProperties {
		truncatedProperty := truncate.StringWithTail(property.Value, propertyWidth, "...")
		prop := lipgloss.PlaceHorizontal(15, lipgloss.Left, truncatedProperty)
		customProperties += propertiesStyle.Render(prop)
	}
	truncatedName := truncate.StringWithTail(i.K8sPod.Name, maxPodNameLength, "...")
	_, _ = fmt.Fprintf(w, "  %s%s", style.Render(truncatedName), customProperties)
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
