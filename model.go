package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strconv"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type K8mpassModel struct {
	error             errMsg
	cluster           kubernetesCluster
	spinner           spinner.Model
	isConnected       bool
	command           NamespaceOperation
	namespaces        []Namespace
	namespacesList    list.Model
	selectedNamespace string
}

func initialModel() K8mpassModel {
	s := spinner.New()
	s.Spinner = spinner.Line
	return K8mpassModel{
		spinner: s,
		command: WakeUpReviewOperation,
	}
}

func (m K8mpassModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, connectToClusterCmd)
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
			if len(m.namespacesList.Items()) > 0 {
				m.selectedNamespace = m.namespacesList.SelectedItem().FilterValue()
			} else {
				m.selectedNamespace = ""
			}
		}
	case clusterConnectedMsg:
		m.isConnected = true
		m.cluster.kubernetes = msg.clientset
		cmds = append(cmds, fetchNamespacesCmd(m.cluster.kubernetes))
	case fetchedNamespacesMsg:
		m.namespaces = msg.namespaces
		var items []list.Item
		for i := 0; i < len(msg.namespaces); i++ {
			items = append(items, msg.namespaces[i])
		}
		m.namespacesList = list.New(items, list.NewDefaultDelegate(), 50, 15)
		m.namespacesList.SetShowStatusBar(false)
		m.namespacesList.SetShowHelp(false)
	case tea.WindowSizeMsg:
		if len(m.namespacesList.Items()) > 0 {
			h, v := docStyle.GetFrameSize()
			m.namespacesList.SetSize(msg.Width-h, msg.Height-v)
		}
	}
	if !m.isConnected {
		s, cmd := m.spinner.Update(msg)
		m.spinner = s
		cmds = append(cmds, cmd)
	}
	if len(m.namespacesList.Items()) > 0 {
		var cmd tea.Cmd
		m.namespacesList, cmd = m.namespacesList.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m K8mpassModel) View() string {
	s := ""
	if !m.isConnected {
		s += m.spinner.View()
		s += "Connecting to the cluster..."
	}
	if m.isConnected && len(m.namespacesList.Items()) > 0 {
		s += "Namespaces: " + strconv.Itoa(len(m.namespacesList.Items())) + "\n"
		s += docStyle.Render(m.namespacesList.View())
		s += "\n Selected: " + m.selectedNamespace
	}

	return s + "\n"
}
