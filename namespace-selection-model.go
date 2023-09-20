package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type NamespaceSelectionModel struct {
}

func (n NamespaceSelectionModel) Init() tea.Cmd {
	return nil
}

func (n NamespaceSelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return n, nil
}

func (n NamespaceSelectionModel) View() string {
	return "We are selecting a namespace now\n"
}
