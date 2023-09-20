package main

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"strconv"
)

type K8mpassModel struct {
	error       errMsg
	cluster     kubernetesCluster
	spinner     spinner.Model
	isConnected bool
	command     NamespaceOperation
	namespaces  []Namespace
}

func initialModel() K8mpassModel {
	s := spinner.New()
	s.Spinner = spinner.Line
	return K8mpassModel{
		spinner: s,
		command: WakeUpReviewOperation,
	}
}

func (m K8mpassModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, connectToClusterCmd)
}

func (m K8mpassModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case errMsg:
		m.error = msg
		return m, tea.Quit
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	case clusterConnectedMsg:
		m.isConnected = true
		m.cluster.kubernetes = msg.clientset
		cmds = append(cmds, fetchNamespacesCmd(m.cluster.kubernetes))
	case fetchedNamespacesMsg:
		m.namespaces = msg.namespaces
	}
	if !m.isConnected {
		s, cmd := m.spinner.Update(msg)
		m.spinner = s
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m K8mpassModel) View() string {
	s := ""
	if !m.isConnected {
		s += m.spinner.View()
		s += "Connecting to the cluster..."
	}
	if m.isConnected && m.namespaces != nil {
		s += "Namespaces: " + strconv.Itoa(len(m.namespaces))
	}

	return s + "\n"
}
