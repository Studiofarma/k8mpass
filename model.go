package main

import (
	"context"
	"github.com/studiofarma/k8mpass/pod"
	"runtime"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/studiofarma/k8mpass/namespace"
)

type K8mpassModel struct {
	error          errMsg
	state          modelState
	namespaceModel NamespaceSelectionModel
	operationModel OperationModel
	podModel       PodSelectionModel
	ctx            *context.Context
}

type modelState int32

const (
	NamespaceSelection modelState = 1
	OperationSelection modelState = 2
	PodSelection       modelState = 3
)

func initialModel() K8mpassModel {
	s := spinner.New()
	s.Spinner = spinner.Line
	ops := []NamespaceOperation{CheckSleepingStatusOperation, WakeUpReviewOperation, PodsOperation, OpenDbmsOperation, OpenApplicationOperation}
	return K8mpassModel{
		state: NamespaceSelection,
		namespaceModel: NamespaceSelectionModel{
			namespaces: namespace.New(),
			messageHandler: namespace.NewHandler(
				//namespace.NamespaceAgeProperty,
				ReviewAppSleepStatus,
			),
		},
		operationModel: OperationModel{
			operations: initializeOperationList(ops),
			helpFooter: initializeHelpFooter(),
		},
		podModel: PodSelectionModel{
			pods:           pod.New(),
			messageHandler: pod.NewHandler(),
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
		return m, tea.Tick(time.Second*5, func(t time.Time) tea.Msg { return tea.Quit })
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
		m.podModel.pods.SetHeight(msg.Height + correction)
		m.podModel.pods.SetWidth(msg.Width + correction)
	case startupMsg:
		cmds = append(cmds, m.namespaceModel.namespaces.StartSpinner())
	case clusterConnectedMsg:
		cmds = append(cmds, m.namespaceModel.messageHandler.GetNamespaces(K8sCluster.kubernetes))
	case namespaceSelectedMsg:
		cmds = append(cmds, m.podModel.messageHandler.GetPods(context.TODO(), K8sCluster.kubernetes, msg.namespace))
		m.state = PodSelection
	case backToNamespaceSelectionMsg:
		m.state = NamespaceSelection
		m.operationModel.Reset()
	case backToOperationSelectionMsg:
		m.state = OperationSelection
		m.operationModel.Reset()
	}
	// Model specific messages
	switch msg.(type) {
	case namespace.Message:
		nm, nmCmd := m.namespaceModel.Update(msg)
		m.namespaceModel = nm
		cmds = append(cmds, nmCmd)
	}

	switch m.state {
	case NamespaceSelection:
		if _, ok := msg.(namespace.Message); ok {
			break
		}
		nm, nmCmd := m.namespaceModel.Update(msg)
		m.namespaceModel = nm
		cmds = append(cmds, nmCmd)
	case OperationSelection:
		om, omCmd := m.operationModel.Update(msg)
		m.operationModel = om
		cmds = append(cmds, omCmd)
	case PodSelection:
		pm, pmCmd := m.podModel.Update(msg)
		m.podModel = pm
		cmds = append(cmds, pmCmd)
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
	case PodSelection:
		return m.podModel.View()
	}
	return ""
}
