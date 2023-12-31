package main

import (
	"github.com/studiofarma/k8mpass/api"
	"github.com/studiofarma/k8mpass/kubernetes"
	"github.com/studiofarma/k8mpass/log"
	"time"

	"github.com/studiofarma/k8mpass/pod"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/studiofarma/k8mpass/namespace"
)

type Model struct {
	error          errMsg
	cluster        kubernetes.ICluster
	namespaceModel NamespaceSelectionModel
	podModel       PodSelectionModel
	state          state
}

type state int32

const (
	NamespaceSelection state = 0
	PodSelection       state = 1
)

func initialModel(plugins api.IPlugins) Model {
	cluster := kubernetes.Cluster{}
	return Model{
		cluster: &cluster,
		state:   NamespaceSelection,
		namespaceModel: NamespaceSelectionModel{
			namespaces: namespace.New(),
			messageHandler: namespace.NewHandler(
				&cluster,
				plugins.GetNamespaceExtensions()...,
			),
		},
		podModel: PodSelectionModel{
			pods: pod.New(),
			messageHandler: pod.NewHandler(
				&cluster,
				plugins.GetPodExtensions(),
				plugins.GetNamespaceOperations(),
			),
			operations: initializeOperationList(),
			logs: NewLogModel(
				&cluster,
			),
		},
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Sequence(
		func() tea.Msg {
			return startupMsg{}
		},
		m.clusterConnect,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case errMsg:
		m.error = msg
		return m, tea.Tick(time.Second*5, func(t time.Time) tea.Msg { return tea.Quit })
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.namespaceModel.userService.Persist()
			return m, tea.Quit
		case "f5":
			cmds = append(cmds, tea.ClearScreen)
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
	case log.Message:
		lm, lmCmd := m.podModel.logs.Update(msg)
		m.podModel.logs = lm
		cmds = append(cmds, lmCmd)
	default:
		nm, nmCmd := m.namespaceModel.Update(msg)
		pm, pmCmd := m.podModel.Update(msg)
		m.namespaceModel = nm
		m.podModel = pm
		cmds = append(cmds, pmCmd, nmCmd)
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

func (m Model) View() string {
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
