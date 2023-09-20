package main

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type NamespacesModel struct {
	state                    sessionState
	entry                    tea.Model
	error                    errMsg
	cluster                  kubernetesCluster
	clusterConnectionSpinner spinner.Model
	isConnected              bool
	command                  NamespaceOperation
	list                     list.Model
	keys                     listKeyMap
	delegateKeys             delegateKeyMap
}

func initialNamespaceModel() NamespacesModel {

	s := spinner.New()
	s.Spinner = spinner.Line
	var (
		//delegateKeys = newDelegateKeyMap()
		listKeys = newListKeyMap()
	)
	namespaces := getNamespaces(createClientSet(getConfig()))
	delegate := newItemDelegate(newDelegateKeyMap())
	namespaceList := list.New(namespaces, delegate, 80, 15)
	namespaceList.Title = "Namespaces"
	namespaceList.Styles.Title = titleStyle
	namespaceList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.toggleSpinner,
			listKeys.insertItem,
			listKeys.toggleTitleBar,
			listKeys.toggleStatusBar,
			listKeys.togglePagination,
			listKeys.toggleHelpMenu,
		}
	}

	return NamespacesModel{
		state:                    namespacesView,
		clusterConnectionSpinner: s,
		list:                     namespaceList,
	}
}

func (m NamespacesModel) Init() tea.Cmd {

	return tea.EnterAltScreen
	//return tea.Batch(m.clusterConnectionSpinner.Tick, clusterConnect)
}

func (m NamespacesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case errMsg:
		m.error = msg
		return m, tea.Quit
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			//torna a main
			return m, tea.Quit
		}
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.toggleSpinner):
			cmd := m.list.ToggleSpinner()
			return m, cmd

		case key.Matches(msg, m.keys.toggleTitleBar):
			v := !m.list.ShowTitle()
			m.list.SetShowTitle(v)
			m.list.SetShowFilter(v)
			m.list.SetFilteringEnabled(v)
			return m, nil

		case key.Matches(msg, m.keys.toggleStatusBar):
			m.list.SetShowStatusBar(!m.list.ShowStatusBar())
			return m, nil

		case key.Matches(msg, m.keys.togglePagination):
			m.list.SetShowPagination(!m.list.ShowPagination())
			return m, nil

		case key.Matches(msg, m.keys.toggleHelpMenu):
			m.list.SetShowHelp(!m.list.ShowHelp())
			return m, nil

		}
	}

	if !m.isConnected {
		s, cmd := m.clusterConnectionSpinner.Update(msg)
		m.clusterConnectionSpinner = s
		cmds = append(cmds, cmd)
	}

	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m NamespacesModel) View() string {
	if m.state == namespacesView {
		return appStyle.Render(m.list.View())
		//return s
	} else {
		return ""
	}

}
