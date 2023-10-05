package main

import (
	"fmt"
	"sort"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/studiofarma/k8mpass/namespace"
)

type NamespaceSelectionModel struct {
	messageHandler *namespace.MessageHandler
	namespaces     list.Model
}

func (n NamespaceSelectionModel) Init() tea.Cmd {
	return nil
}

func (n NamespaceSelectionModel) Update(msg tea.Msg) (NamespaceSelectionModel, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case namespace.WatchingMsg:
		cmds = append(cmds, n.messageHandler.NextEvent)
		cmds = append(cmds, namespace.Refresh())
	case namespace.ListMsg:
		items := make([]list.Item, len(msg.Namespaces))
		for i, ns := range msg.Namespaces {
			items[i] = ns
		}
		cmds = append(cmds, n.namespaces.SetItems(items))
		//cmds = append(cmds, n.namespaces.SetItems(items))
		n.namespaces.StopSpinner()
		n.namespaces.Title = "Select a namespace"
	case namespace.AddedMsg:
		cmds = append(cmds, n.namespaces.InsertItem(0, msg.Namespace))
		ns := n.namespaces.Items()
		sort.SliceStable(ns, func(i, j int) bool {
			return ns[i].FilterValue() < ns[j].FilterValue()
		})
		cmds = append(cmds, n.namespaces.SetItems(ns))
		cmds = append(cmds, n.messageHandler.NextEvent)
		cmds = append(cmds, n.namespaces.NewStatusMessage(fmt.Sprintf("ADDED: %s", msg.Namespace.K8sNamespace.Name)))
	case namespace.ModifiedMsg:
		var idx = namespace.FindNamespace(n.namespaces.Items(), msg.Namespace)
		cmds = append(cmds, n.namespaces.SetItem(idx, msg.Namespace))
		cmds = append(cmds, n.messageHandler.NextEvent)
	case namespace.RemovedMsg:
		var idx = namespace.FindNamespace(n.namespaces.Items(), msg.Namespace)
		n.namespaces.RemoveItem(idx)
		cmds = append(cmds, n.messageHandler.NextEvent)
		cmds = append(cmds, n.namespaces.NewStatusMessage(fmt.Sprintf("REMOVED: %s", msg.Namespace.K8sNamespace.Name)))
	case namespace.NextEventMsg:
		cmds = append(cmds, n.messageHandler.NextEvent)
	case namespace.ErrorMsg:
		n.namespaces.NewStatusMessage(msg.Err.Error())
	case namespace.ReloadTick:
		var namespaces []namespace.Item
		for _, item := range n.namespaces.Items() {
			namespaces = append(namespaces, item.(namespace.Item))
		}
		cmds = append(cmds, n.messageHandler.ReloadExtensions(namespaces))
	case namespace.ReloadExtensionsMsg:
		items := n.namespaces.Items()
		for idx, item := range items {
			ns := item.(namespace.Item)
			property := msg.Properties[ns.K8sNamespace.Name]
			if property == nil {
				continue
			} else {
				ns.ExtendedProperties = msg.Properties[ns.K8sNamespace.Name]
				items[idx] = ns
			}
		}
		cmds = append(cmds, n.namespaces.SetItems(items))
		cmds = append(cmds, namespace.Refresh())
		cmds = append(cmds, n.namespaces.NewStatusMessage("Reloaded"))
	case tea.KeyMsg:
		if n.namespaces.FilterState() == list.Filtering {
			break
		}
		switch keypress := msg.String(); keypress {
		case "enter":
			i, ok := n.namespaces.SelectedItem().(namespace.Item)
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
