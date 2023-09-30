package main

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"runtime"
)

type K8mpassModel struct {
	error          errMsg
	state          modelState
	namespaceModel NamespaceSelectionModel
	operationModel OperationModel
}

type modelState int32

const (
	NamespaceSelection modelState = 1
	OperationSelection modelState = 2
)

func initialModel() K8mpassModel {
	s := spinner.New()
	s.Spinner = spinner.Line
	ops := []NamespaceOperation{CheckSleepingStatusOperation, WakeUpReviewOperation, PodsOperation, OpenDbmsOperation, OpenApplicationOperation}
	return K8mpassModel{
		state: NamespaceSelection,
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
	return tea.Batch(clusterConnect, func() tea.Msg {
		return startupMsg{}
	})
}

func (m K8mpassModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case errMsg:
		m.error = msg
		return m, tea.Quit
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		var correction int
		switch runtime.GOOS {
		case "windows":
			correction = -1 // -1 is for Windows which doesn't handle the size correctly
		default:
			correction = 0
		}
		m.namespaceModel.namespaces.SetHeight(msg.Height + correction)
		m.namespaceModel.namespaces.SetWidth(msg.Width + correction)
		m.operationModel.operations.SetHeight(msg.Height + correction)
		m.operationModel.operations.SetWidth(msg.Width + correction)
	case startupMsg:
		cmds = append(cmds, m.namespaceModel.namespaces.StartSpinner())
	case clusterConnectedMsg:
		cmds = append(cmds, fetchNamespaces)
	case namespacesRetrievedMsg:
		m.state = NamespaceSelection
	case namespaceSelectedMsg:
		m.state = OperationSelection
	case backToNamespaceSelectionMsg:
		m.state = NamespaceSelection
		m.operationModel.Reset()
		//m.namespaceModel.Reset()
	case backToOperationSelectionMsg:
		m.state = OperationSelection
		m.operationModel.Reset()
	}
	switch m.state {
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
	if m.error != nil {
		return m.error.Error()
	}
	switch m.state {
	case NamespaceSelection:
		return m.namespaceModel.View()
	case OperationSelection:
		return m.operationModel.View()
	}
	return ""
}
