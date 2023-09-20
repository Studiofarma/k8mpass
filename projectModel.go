package main

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type sessionState int

const (
	mainView       sessionState = 0
	namespacesView sessionState = 1
	podsView       sessionState = 2
	cronjobsView                = 3
)

type K8mpassModel struct {
	state                    sessionState
	entry                    tea.Model
	error                    errMsg
	cluster                  kubernetesCluster
	clusterConnectionSpinner spinner.Model
	isConnected              bool
	command                  NamespaceOperation
}

func initialProjectModel(m K8mpassModel) K8mpassModel {
	if m == nil {
		s := spinner.New()
		s.Spinner = spinner.Line

		return K8mpassModel{
			state:                    mainView,
			clusterConnectionSpinner: s,
		}
	} else {
		return K8mpassModel{
			state:                    m.state,
			entry:                    m.entry,
			error:                    m.error,
			cluster:                  m.cluster,
			clusterConnectionSpinner: m.clusterConnectionSpinner,
			isConnected:              m.isConnected,
			command:                  m.command,
		}
	}
}

func (m K8mpassModel) Init() tea.Cmd {
	return tea.EnterAltScreen
	//return tea.Batch(m.clusterConnectionSpinner.Tick, clusterConnect)
}

func (m K8mpassModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case errMsg:
		m.error = msg
		return m, tea.Quit

	case clusterConnectedMsg:
		m.isConnected = true
		m.cluster.kubernetes = msg.clientset
		command := m.command.Command(m, "review-hack-cgmgpharm-47203-be")
		cmds = append(cmds, command)
	}

	if !m.isConnected {
		s, cmd := m.clusterConnectionSpinner.Update(msg)
		m.clusterConnectionSpinner = s
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m K8mpassModel) View() string {
	s := "Sono bello bravo e funzionante"
	if m.entry == nil {
		n := initialNamespaceModel()
		m.entry = n
		return appStyle.Render(n.View())
	} else {
		switch m.entry.(type) {
		case K8mpassModel:
			//appStyle.Render(s)
			n := initialNamespaceModel()
			m.entry = n
			initialProjectModel(m)
			return appStyle.Render(n.View())
		}

	}

	return s
}
