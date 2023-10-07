package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) clusterConnect() tea.Msg {
	err := m.cluster.Connect()
	if err != nil {
		return errMsg(err)
	}
	return clusterConnectedMsg{
		context: m.cluster.GetContext(),
	}
}
