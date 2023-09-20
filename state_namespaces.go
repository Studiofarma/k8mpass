package main

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type NamespacesModel struct {
	error      errMsg
	cluster    kubernetesCluster
	spinner    spinner.Model
	namespaces []Namespace
}

func createNamespacesModel(cluster kubernetesCluster) NamespacesModel {
	s := spinner.New()
	s.Spinner = spinner.Line
	return NamespacesModel{
		cluster: cluster,
		spinner: s,
	}
}

func (m NamespacesModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, fetchNamespacesCmd(m.cluster.kubernetes))
}

func (m NamespacesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	case fetchedNamespacesMsg:
		m.namespaces = msg.namespaces
	}

	if m.namespaces == nil {
		s, cmd := m.spinner.Update(msg)
		m.spinner = s
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m NamespacesModel) View() string {
	s := ""
	if m.namespaces == nil {
		s += m.spinner.View()
		s += "Loading namespaces..."
	} else {
		s += "I got the LIST!"
	}

	return s + "\n"
}
