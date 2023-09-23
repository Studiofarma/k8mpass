package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type navBarState int32

const (
	connectingToCluster navBarState = 0
	fetchingNamespaces  navBarState = 1

	refreshingNamespaces navBarState = 2
	namespaceSelection   navBarState = 3
	operationSelection   navBarState = 4
	operationOutput      navBarState = 5
)

type NavigationBarModel struct {
	state     navBarState
	context   string
	namespace string
	operation string
	spinner   spinner.Model
	loading   bool
}

func initNavBar() NavigationBarModel {
	s := spinner.Model{}
	s.Spinner = spinner.Dot
	return NavigationBarModel{
		state:   connectingToCluster,
		spinner: s,
		loading: true,
	}
}

func (m NavigationBarModel) Init() tea.Cmd {
	return nil
}

func (m NavigationBarModel) Update(msg tea.Msg) (NavigationBarModel, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case clusterConnectedMsg:
		m.state = fetchingNamespaces
		m.context = "k8s-context"
	case namespacesFetchedMsg:
		m.state = namespaceSelection
		m.loading = false
	case namespaceSelectedMsg:
		m.state = operationSelection
		m.namespace = msg.namespace
	case refreshNamespacesMsg:
		m.state = refreshingNamespaces
		m.loading = true
		cmds = append(cmds, m.spinner.Tick)
	case backToNamespaceSelectionMsg:
		m.state = namespaceSelection
		m.namespace = ""
	case backToOperationSelectionMsg:
		m.state = operationSelection
		m.operation = ""
	}
	if m.loading {
		s, sCmd := m.spinner.Update(msg)
		m.spinner = s
		cmds = append(cmds, sCmd)
	}
	return m, tea.Batch(cmds...)
}

func (m NavigationBarModel) View() string {
	var spinner string
	if m.loading {
		spinner = fmt.Sprintf(" %s ", m.spinner.View())
	} else {
		spinner = "    "
	}
	switch m.state {
	case connectingToCluster:
		return fmt.Sprintf("%sConnecting to the cluster...\n", spinner)
	case fetchingNamespaces:
		return fmt.Sprintf("%sFetching namespaces\n", spinner)
	case refreshingNamespaces:
		return fmt.Sprintf("%sRefreshing namespaces...\n", spinner)
	case namespaceSelection:
		return fmt.Sprintf("%sSelect namespace\n", spinner)
	case operationSelection:
		return fmt.Sprintf("%sSelection operation on %s\n", spinner, m.namespace)
	case operationOutput:
		return fmt.Sprintf("%sOutput of operation %s\n", spinner, m.operation)
	default:
		return ""
	}
}
