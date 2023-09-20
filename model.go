package main

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type K8mpassModel struct {
	error                    errMsg
	cluster                  kubernetesCluster
	clusterConnectionSpinner spinner.Model
	command                  NamespaceOperation
	namespaceModel           NamespaceSelectionModel
	operationModel           OperationModel
	state                    modelState
}

type modelState int32

const (
	Connection         modelState = 0
	NamespaceSelection modelState = 1
	OperationSelection modelState = 2
)

func initialModel() K8mpassModel {
	s := spinner.New()
	s.Spinner = spinner.Line
	return K8mpassModel{
		clusterConnectionSpinner: s,
		command:                  WakeUpReviewOperation,
		state:                    Connection,
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
		}
	case clusterConnectedMsg:
		m.state = NamespaceSelection
		//m.cluster.kubernetes = msg.clientset
		//command := m.command.Command(m, "review-devops-new-filldata")
		//cmds = append(cmds, command)
	}
	switch m.state {
	case Connection:
		sm, smCmd := m.clusterConnectionSpinner.Update(msg)
		m.clusterConnectionSpinner = sm
		cmds = append(cmds, smCmd)
	}
	return m, tea.Batch(cmds...)
}

func (m K8mpassModel) View() string {
	switch m.state {
	case Connection:
		return m.clusterConnectionSpinner.View() + " Connecting to Kubernetes Cluster...\n"
	case NamespaceSelection:
		return m.namespaceModel.View()
	case OperationSelection:
		return m.operationModel.View()
	}
	return ""
}
