package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	ops := []NamespaceOperation{WakeUpReviewOperation, PodsOperation}
	return K8mpassModel{
		clusterConnectionSpinner: s,
		command:                  WakeUpReviewOperation,
		state:                    Connection,
		namespaceModel: NamespaceSelectionModel{
			namespaces: initializeList(),
		},
		operationModel: OperationModel{
			operations: initializeOperationList(ops),
			helpFooter: initializeHelpFooter(),
		},
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
		m.cluster.kubernetes = msg.clientset
		m.operationModel.clientset = msg.clientset
		c := func() tea.Msg {
			ns, err := getNamespaces(m.cluster.kubernetes)
			if err != nil {
				m.error = errMsg(err)
			}
			var nsNames []string
			for _, n := range ns.Items {
				nsNames = append(nsNames, n.Name)
			}
			return namespacesRetrievedMsg{nsNames}
		}
		cmds = append(cmds, c)
	case namespacesRetrievedMsg:
		m.state = NamespaceSelection
		var items []list.Item
		for _, n := range msg.namespaces {
			items = append(items, NamespaceItem{n})
		}
		m.namespaceModel.namespaces.SetItems(items)
	case namespaceSelectedMsg:
		m.state = OperationSelection
		m.operationModel.namespace = msg.namespace
		styledNamespace := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("170")).Render(msg.namespace)
		m.operationModel.operations.NewStatusMessage(styledNamespace)
	case backToNamespaceSelectionMsg:
		m.state = NamespaceSelection
		m.operationModel.Reset()
		m.namespaceModel.Reset()
	case backToOperationSelectionMsg:
		m.state = OperationSelection
		m.operationModel.Reset()
	}
	switch m.state {
	case Connection:
		sm, smCmd := m.clusterConnectionSpinner.Update(msg)
		m.clusterConnectionSpinner = sm
		cmds = append(cmds, smCmd)
	case NamespaceSelection:
		nm, nmCmd := m.namespaceModel.Update(msg)
		m.namespaceModel = nm
		cmds = append(cmds, nmCmd)
	case OperationSelection:
		om, omCmd := m.operationModel.Update(msg)
		m.operationModel = om
		cmds = append(cmds, omCmd)
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
