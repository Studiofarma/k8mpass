package main

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type K8mpassModel struct {
	error                    errMsg
	cluster                  kubernetesCluster
	clusterConnectionSpinner spinner.Model
	textInput                textinput.Model
	isConnected              bool
	inputRequired            bool
	command                  NamespaceOperation
}

func initialModel() K8mpassModel {
	s := spinner.New()
	s.Spinner = spinner.Line

	txt := textinput.New()
	txt.Placeholder = "write namespace name here"

	return K8mpassModel{
		clusterConnectionSpinner: s,
		textInput:                txt,
		command:                  WakeUpReviewOperation,
	}
}

func (m K8mpassModel) Init() tea.Cmd {
	return tea.Batch(m.clusterConnectionSpinner.Tick, clusterConnect)
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
		case "enter":
			m.inputRequired = false
			command := m.command.Command(m, m.textInput.Value())
			cmds = append(cmds, command)
		}

		if m.inputRequired {
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			cmds = append(cmds, cmd)
		}

	case clusterConnectedMsg:
		m.isConnected = true
		m.cluster.kubernetes = msg.clientset
		m.inputRequired = true
		m.textInput.Focus()
		//command := m.command.Command(m, "review-devops-new-filldata")
		//cmds = append(cmds, command)
	}
	if !m.isConnected {
		s, cmd := m.clusterConnectionSpinner.Update(msg)
		m.clusterConnectionSpinner = s
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m K8mpassModel) View() string {
	s := ""
	if !m.isConnected {
		s += m.clusterConnectionSpinner.View()
		s += "Connecting to the cluster..."
	} else {
		if m.error != nil {
			s += m.error.Error()
		} else {
			if m.inputRequired {
				return "namespace name:" + m.textInput.View()
			}
			s += "Connection successful! Press esc to quit"
		}
	}
	return s
}
