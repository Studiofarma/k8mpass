package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type NamespaceSelectionModel struct {
	namespaces list.Model
}

func (n NamespaceSelectionModel) Init() tea.Cmd {
	return nil
}

func (n NamespaceSelectionModel) Update(msg tea.Msg) (NamespaceSelectionModel, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if n.namespaces.FilterState() == list.Filtering {
			break
		}
		switch keypress := msg.String(); keypress {
		case "enter":
			i, ok := n.namespaces.SelectedItem().(NamespaceItem)
			if ok {
				nsCommand := func() tea.Msg {
					return namespaceSelectedMsg{i.name}
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
