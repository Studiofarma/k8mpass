package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type OperationModel struct {
	namespace string
}

func (o OperationModel) Init() tea.Cmd {
	return nil
}

func (o OperationModel) Update(msg tea.Msg) (OperationModel, tea.Cmd) {
	switch msg := msg.(type) {
	case namespaceSelectedMsg:
		o.namespace = msg.namespace
	}
	return o, nil
}

func (o OperationModel) View() string {
	return "You selected the namespace" + o.namespace
}
