package main

import (
	"sort"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/studiofarma/k8mpass/namespace"
	"k8s.io/apimachinery/pkg/watch"
)

type NamespaceSelectionModel struct {
	namespaces list.Model
	nsChannel  <-chan watch.Event
}

func (n NamespaceSelectionModel) Init() tea.Cmd {
	return nil
}

func (n NamespaceSelectionModel) Update(msg tea.Msg) (NamespaceSelectionModel, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case watchNamespaceMsg:
		n.nsChannel = msg.channel
		cmds = append(cmds, nextEvent(n.nsChannel))
		n.namespaces.StopSpinner()
		n.namespaces.Title = "Select a namespace"
	case namespace.AddedNamespaceMsg:
		cmds = append(cmds, n.namespaces.InsertItem(0, msg.Namespace))
		ns := n.namespaces.Items()
		sort.SliceStable(ns, func(i, j int) bool {
			return ns[i].FilterValue() < ns[j].FilterValue()
		})
		cmds = append(cmds, n.namespaces.SetItems(ns))
	case namespace.RemovedNamespaceMsg:
		var idx = -1
		for i, v := range n.namespaces.Items() {
			if v.FilterValue() == msg.Namespace.FilterValue() {
				idx = i
			}
		}
		n.namespaces.RemoveItem(idx)
	case tea.KeyMsg:
		if n.namespaces.FilterState() == list.Filtering {
			break
		}
		switch keypress := msg.String(); keypress {
		case "enter":
			i, ok := n.namespaces.SelectedItem().(namespace.NamespaceItem)
			if ok {
				nsCommand := func() tea.Msg {
					return namespaceSelectedMsg{i.K8sNamespace.Name}
				}
				cmds = append(cmds, nsCommand)

			} else {
				panic("Casting went wrong")
			}
		}
	}

	lm, lmCmd := n.namespaces.Update(msg)
	n.namespaces = lm
	cmds = append(cmds, lmCmd)

	return n, tea.Batch(cmds...)
}

func (n NamespaceSelectionModel) View() string {
	return n.namespaces.View()
}

func (n *NamespaceSelectionModel) Reset() {
	n.namespaces.ResetSelected()
	n.namespaces.ResetFilter()
}
