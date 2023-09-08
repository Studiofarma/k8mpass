package main

import tea "github.com/charmbracelet/bubbletea"

func clusterConnect() tea.Msg {
	cs, err := getConnection()
	if err != nil {
		return errMsg(err)
	}
	return clusterConnectedMsg{cs}
}
