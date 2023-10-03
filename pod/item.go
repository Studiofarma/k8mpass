package pod

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
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

func (n ItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (n ItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Item)
	if !ok {
		return
	}
	fmt.Fprintf(w, "  %s", styleString(i.K8sPod.Name, podStyle(i.K8sPod.Status)))
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
