package main

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type K8mpassModel struct {
	error                    errMsg
	cluster                  kubernetesCluster
	clusterConnectionSpinner spinner.Model
	isConnected              bool
	command                  NamespaceOperation
	namespace                string
}

/*
func initialModel() K8mpassModel {
	s := spinner.New()
	s.Spinner = spinner.Line
	return K8mpassModel{
		clusterConnectionSpinner: s,
		command:                  WakeUpReviewOperation,
	}
}
*/

func (m K8mpassModel) Init() tea.Cmd {
	return tea.Batch(m.clusterConnectionSpinner.Tick, clusterConnect)
}

func (m K8mpassModel) View() string {
	s := ""
	if !m.isConnected {
		s += m.clusterConnectionSpinner.View()
		s += "Connecting to the cluster..."
	} else {
		s += "Connection successful! Press esc to quit"
	}
	return s
}
