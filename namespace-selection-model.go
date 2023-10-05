package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/studiofarma/k8mpass/namespace"
)

type NamespaceSelectionModel struct {
	messageHandler *namespace.MessageHandler
	namespaces     list.Model
}

func (m NamespaceSelectionModel) Init() tea.Cmd {
	return nil
}

func (m NamespaceSelectionModel) Update(msg tea.Msg) (NamespaceSelectionModel, tea.Cmd) {
	var cmds []tea.Cmd
	var routedCmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.namespaces.SetSize(msg.Width, msg.Height)
	case startupMsg:
		routedCmds = append(cmds, m.namespaces.StartSpinner())
	case clusterConnectedMsg:
		cmds = append(routedCmds, m.messageHandler.GetNamespaces(K8sCluster.kubernetes))
	case namespace.WatchingMsg:
		cmds = append(cmds, m.messageHandler.NextEvent)
		cmds = append(cmds, namespace.Refresh())
	case namespace.ListMsg:
		items := make([]list.Item, len(msg.Namespaces))
		for i, ns := range msg.Namespaces {
			items[i] = ns
		}
		cmds = append(routedCmds, m.namespaces.SetItems(items))
		m.namespaces.StopSpinner()
		m.namespaces.SetShowPagination(true) //This is needed to overcome an annoying graphical bug https://github.com/charmbracelet/bubbles/issues/405
		m.namespaces.Title = "Select a namespace"
	case namespace.AddedMsg:
		cmds = append(routedCmds, m.namespaces.InsertItem(0, msg.Namespace))
		ns := m.namespaces.Items()
		sort.SliceStable(ns, func(i, j int) bool {
			return ns[i].FilterValue() < ns[j].FilterValue()
		})
		cmds = append(routedCmds, m.namespaces.SetItems(ns))
		cmds = append(cmds, m.messageHandler.NextEvent)
		cmds = append(routedCmds, m.namespaces.NewStatusMessage(fmt.Sprintf("ADDED: %s", msg.Namespace.K8sNamespace.Name)))
	case namespace.ModifiedMsg:
		var idx = namespace.FindNamespace(m.namespaces.Items(), msg.Namespace)
		cmds = append(routedCmds, m.namespaces.SetItem(idx, msg.Namespace))
		cmds = append(cmds, m.messageHandler.NextEvent)
	case namespace.RemovedMsg:
		var idx = namespace.FindNamespace(m.namespaces.Items(), msg.Namespace)
		m.namespaces.RemoveItem(idx)
		cmds = append(cmds, m.messageHandler.NextEvent)
		cmds = append(routedCmds, m.namespaces.NewStatusMessage(fmt.Sprintf("REMOVED: %s", msg.Namespace.K8sNamespace.Name)))
	case namespace.NextEventMsg:
		cmds = append(cmds, m.messageHandler.NextEvent)
	case namespace.ErrorMsg:
		m.namespaces.NewStatusMessage(msg.Err.Error())
	case namespace.ReloadTick:
		var namespaces []namespace.Item
		for _, item := range m.namespaces.Items() {
			namespaces = append(namespaces, item.(namespace.Item))
		}
		cmds = append(cmds, m.messageHandler.ReloadExtensions(namespaces))
	case namespace.ReloadExtensionsMsg:
		items := m.namespaces.Items()
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
		cmds = append(routedCmds, m.namespaces.SetItems(items))
		cmds = append(cmds, namespace.Refresh())
		cmds = append(routedCmds, m.namespaces.NewStatusMessage("Reloaded"))
		log.Printf("Namespace - Width: %d; Height: %d", m.namespaces.Width(), m.namespaces.Height())
	case namespace.RoutedMsg:
		model, cmd := m.namespaces.Update(msg.Embedded)
		m.namespaces = model
		cmds = append(routedCmds, cmd)
	case tea.KeyMsg:
		if m.namespaces.FilterState() == list.Filtering {
			break
		}
		switch keypress := msg.String(); keypress {
		case "enter":
			i, ok := m.namespaces.SelectedItem().(namespace.Item)
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

	lm, lmCmd := m.namespaces.Update(msg)
	m.namespaces = lm
	cmds = append(cmds, lmCmd)
	cmds = append(cmds, namespace.Route(routedCmds...)...)
	return m, tea.Batch(cmds...)
}

func (m NamespaceSelectionModel) View() string {
	return m.namespaces.View()
}

func (m *NamespaceSelectionModel) Reset() {
	m.namespaces.ResetSelected()
	m.namespaces.ResetFilter()
}
