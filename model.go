package main

import (
	"github.com/studiofarma/k8mpass/api"
	"time"

	"github.com/studiofarma/k8mpass/pod"

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
	NamespaceSelection modelState = 0
	PodSelection       modelState = 1
)

func initialModel(plugins api.IPlugins) K8mpassModel {
	return K8mpassModel{
		state: NamespaceSelection,
		namespaceModel: NamespaceSelectionModel{
			namespaces: namespace.New(),
			messageHandler: namespace.NewHandler(
				plugins.GetNamespaceExtensions()...,
			),
		},
		podModel: PodSelectionModel{
			pods: pod.New(),
			messageHandler: pod.NewHandler(
				plugins.GetPodExtensions()...),
			availableOperations: plugins.GetNamespaceOperations(),
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
		default:
			switch m.state {
			case NamespaceSelection:
				model, cmd := m.namespaceModel.Update(msg)
				m.namespaceModel = model
				cmds = append(cmds, cmd)
			case PodSelection:
				model, cmd := m.podModel.Update(msg)
				m.podModel = model
				cmds = append(cmds, cmd)
			}
		}
	case namespace.Message:
		nm, nmCmd := m.namespaceModel.Update(msg)
		m.namespaceModel = nm
		cmds = append(cmds, nmCmd)
	case pod.Message:
		pm, pmCmd := m.podModel.Update(msg)
		m.podModel = pm
		cmds = append(cmds, pmCmd)
	default:
		nm, nmCmd := m.namespaceModel.Update(msg)
		m.namespaceModel = nm
		cmds = append(cmds, nmCmd)
		pm, pmCmd := m.podModel.Update(msg)
		m.podModel = pm
		cmds = append(cmds, pmCmd)
	}

	switch msg.(type) {
	case namespaceSelectedMsg:
		m.state = PodSelection
	case backToNamespaceSelectionMsg:
		m.state = NamespaceSelection
		m.podModel.Reset()
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
