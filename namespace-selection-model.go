package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type NamespaceSelectionModel struct {
	namespaces []string
	list       list.Model
}

func (n NamespaceSelectionModel) Init() tea.Cmd {
	return nil
}

func (n NamespaceSelectionModel) Update(msg tea.Msg) (NamespaceSelectionModel, tea.Cmd) {
	var cmds []tea.Cmd
	lm, lmCmd := n.list.Update(msg)
	n.list = lm
	cmds = append(cmds, lmCmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			i, ok := n.list.SelectedItem().(NamespaceItem)
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

	return n, tea.Batch(cmds...)
}

func (n NamespaceSelectionModel) View() string {
	return n.list.View()
}
