package main

import (
	"context"
	"github.com/studiofarma/k8mpass/api"
	"runtime"
	"time"

	"github.com/studiofarma/k8mpass/pod"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/studiofarma/k8mpass/namespace"
)

type K8mpassModel struct {
	error          errMsg
	state          modelState
	namespaceModel NamespaceSelectionModel
	podModel       PodSelectionModel
}

type modelState int32

const (
	NamespaceSelection modelState = 1
	PodSelection       modelState = 3
)

func initialModel(extensions []api.IExtension, operations []api.INamespaceOperation) K8mpassModel {
	s := spinner.New()
	s.Spinner = spinner.Line
	return K8mpassModel{
		state: NamespaceSelection,
		namespaceModel: NamespaceSelectionModel{
			namespaces: namespace.New(),
			messageHandler: namespace.NewHandler(
				extensions...,
			),
		},
		podModel: PodSelectionModel{
			pods:                pod.New(),
			messageHandler:      pod.NewHandler(),
			availableOperations: operations,
			operations:          initializeOperationList(),
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
		case "s":
			m.podModel.messageHandler.StopWatching()

		}
	case tea.WindowSizeMsg:
		var correction int
		switch runtime.GOOS {
		case "windows":
			correction = -1 // -1 is for Windows which doesn't handle the size correctly
		default:
			correction = 0
		}
		m.namespaceModel.namespaces.SetSize(msg.Width, msg.Height+correction)
		m.podModel.dimensions = struct {
			width  int
			height int
		}{width: msg.Width, height: msg.Height}
		m.podModel.UpdateSize()
		m.podModel.operations.SetWidth(msg.Width)
		m.podModel.pods.SetSize(msg.Width, msg.Height+correction-m.podModel.operations.Height())
	case startupMsg:
		cmds = append(cmds, m.namespaceModel.namespaces.StartSpinner())
	case clusterConnectedMsg:
		cmds = append(cmds, m.namespaceModel.messageHandler.GetNamespaces(K8sCluster.kubernetes))
	case namespaceSelectedMsg:
		m.podModel.namespace = msg.namespace
		cmds = append(cmds, m.podModel.Init())
		cmds = append(cmds, m.podModel.messageHandler.GetPods(context.TODO(), K8sCluster.kubernetes, msg.namespace))
		m.podModel.operations.Title = msg.namespace
		cmds = append(cmds, m.podModel.operations.StartSpinner())
		m.state = PodSelection
	case backToNamespaceSelectionMsg:
		m.state = NamespaceSelection
		cmds = append(cmds, m.podModel.Reset())
	}
	// Model specific messages
	switch msg.(type) {
	case namespace.Message:
		nm, nmCmd := m.namespaceModel.Update(msg)
		m.namespaceModel = nm
		cmds = append(cmds, nmCmd)
	case pod.Message:
		pm, pmCmd := m.podModel.Update(msg)
		m.podModel = pm
		cmds = append(cmds, pmCmd)
	}

	switch m.state {
	case NamespaceSelection:
		if _, ok := msg.(namespace.Message); ok {
			break
		}
		nm, nmCmd := m.namespaceModel.Update(msg)
		m.namespaceModel = nm
		cmds = append(cmds, nmCmd)
	case PodSelection:
		if _, ok := msg.(pod.Message); ok {
			break
		}
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
	case PodSelection:
		return m.podModel.View()
	}
	return ""
}
