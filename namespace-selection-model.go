package main

import (
	"fmt"
	"github.com/studiofarma/k8mpass/config"
	"slices"
	"sort"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/studiofarma/k8mpass/namespace"
)

type NamespaceSelectionModel struct {
	messageHandler *namespace.MessageHandler
	userService    config.IUserService
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
		routedCmds = append(routedCmds, m.namespaces.StartSpinner())
	case clusterConnectedMsg:
		m.namespaces.Title = msg.context
		m.userService = config.New(msg.context)
		m.namespaces.SetDelegate(namespace.ItemDelegate{Pinned: m.userService.GetPinnedNamespaces()})
		cmds = append(cmds, m.messageHandler.GetNamespaces())
	case namespace.WatchingMsg:
		cmds = append(cmds, m.messageHandler.NextEvent)
		cmds = append(cmds, namespace.Refresh())
	case namespace.ListMsg:
		items := make([]list.Item, len(msg.Namespaces))
		for i, ns := range msg.Namespaces {
			items[i] = ns
		}
		routedCmds = append(routedCmds, m.namespaces.SetItems(SortWithFavourites(items, m.userService.GetPinnedNamespaces())))
		m.WorkaroundForGraphicalBug()
		m.namespaces.StopSpinner()
	case namespace.AddedMsg:
		_ = m.namespaces.InsertItem(0, msg.Namespace)
		ns := m.namespaces.Items()
		sort.SliceStable(ns, func(i, j int) bool {
			return ns[i].FilterValue() < ns[j].FilterValue()
		})
		routedCmds = append(routedCmds, m.namespaces.SetItems(SortWithFavourites(ns, m.userService.GetPinnedNamespaces())))
		m.WorkaroundForGraphicalBug()
		cmds = append(cmds, m.messageHandler.NextEvent)
		routedCmds = append(routedCmds, m.namespaces.NewStatusMessage(fmt.Sprintf("ADDED: %s", msg.Namespace.K8sNamespace.Name)))
	case namespace.ModifiedMsg:
		var idx = FindItem(m.namespaces.Items(), msg.Namespace)
		routedCmds = append(routedCmds, m.namespaces.SetItem(idx, msg.Namespace))
		m.WorkaroundForGraphicalBug()
		cmds = append(cmds, m.messageHandler.NextEvent)
	case namespace.RemovedMsg:
		var idx = FindItem(m.namespaces.Items(), msg.Namespace)
		m.namespaces.RemoveItem(idx)
		cmds = append(cmds, m.messageHandler.NextEvent)
		routedCmds = append(routedCmds, m.namespaces.NewStatusMessage(fmt.Sprintf("REMOVED: %s", msg.Namespace.K8sNamespace.Name)))
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
		routedCmds = append(routedCmds, m.namespaces.SetItems(items))
		routedCmds = append(routedCmds, m.namespaces.NewStatusMessage("Reloaded"))
		cmds = append(cmds, namespace.Refresh())
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
		case "p":
			m.userService.Pin(m.namespaces.SelectedItem().FilterValue())
			m.namespaces.SetDelegate(namespace.ItemDelegate{Pinned: m.userService.GetPinnedNamespaces()})
			routedCmds = append(routedCmds, m.namespaces.SetItems(SortWithFavourites(m.namespaces.Items(), m.userService.GetPinnedNamespaces())))
			m.WorkaroundForGraphicalBug()
		case "u":
			m.userService.Unpin(m.namespaces.SelectedItem().FilterValue())
			ns := m.namespaces.Items()
			sort.SliceStable(ns, func(i, j int) bool {
				return ns[i].FilterValue() < ns[j].FilterValue()
			})
			m.namespaces.SetDelegate(namespace.ItemDelegate{Pinned: m.userService.GetPinnedNamespaces()})
			routedCmds = append(routedCmds, m.namespaces.SetItems(SortWithFavourites(ns, m.userService.GetPinnedNamespaces())))
			m.WorkaroundForGraphicalBug()
		}
	}
	switch msg := msg.(type) {
	case namespace.RoutedMsg:
		if batchMsg, ok := msg.Embedded.(tea.BatchMsg); ok {
			cmds = append(routedCmds, func() tea.Msg {
				return batchMsg
			})
			break
		}
		lm, lmCmd := m.namespaces.Update(msg.Embedded)
		m.namespaces = lm
		routedCmds = append(routedCmds, lmCmd)
	default:
		lm, lmCmd := m.namespaces.Update(msg)
		m.namespaces = lm
		routedCmds = append(routedCmds, lmCmd)
	}

	cmds = append(cmds, namespace.Route(routedCmds)...)
	return m, tea.Batch(cmds...)
}

func (m NamespaceSelectionModel) View() string {
	return m.namespaces.View()
}

func (m *NamespaceSelectionModel) Reset() {
	m.namespaces.ResetSelected()
	m.namespaces.ResetFilter()
}

func SortWithFavourites(items []list.Item, pinned []string) []list.Item {
	sortedNames := pinned
	for _, item := range items {
		if slices.Contains(pinned, item.FilterValue()) {
			continue
		}
		sortedNames = append(sortedNames, item.FilterValue())
	}
	mappedItems := make(map[string]list.Item)
	for _, item := range items {
		mappedItems[item.FilterValue()] = item
	}
	var res []list.Item
	for _, n := range sortedNames {
		if mappedItems[n] == nil {
			continue
		}
		res = append(res, mappedItems[n])
	}
	return res
}

// WorkaroundForGraphicalBug This is needed to overcome an annoying graphical bug https://github.com/charmbracelet/bubbles/issues/405
func (m *NamespaceSelectionModel) WorkaroundForGraphicalBug() {
	m.namespaces.SetShowPagination(true)
}
