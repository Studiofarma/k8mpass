package main

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type K8mNamespaceModel struct {
	choices                  []string // items on the to-do list
	filteredChoices          []string
	cursor                   int    // which to-do list item our cursor is pointing at
	selectedNamespace        string // which to-do items are selectedNamespace
	cluster                  kubernetesCluster
	clusterConnectionSpinner spinner.Model
	isConnected              bool
	command                  NamespaceOperation
	error                    errMsg
	filter                   textinput.Model
	views                    []string
	activeView               string
	commands                 []NamespaceOperation
	podIteams                []string
}

func initialModel() K8mNamespaceModel {
	filter := textinput.New()
	filter.Focus()
	filter.Placeholder = "Search"

	return K8mNamespaceModel{
		// Our to-do list is a grocery list
		choices:         []string{},
		filteredChoices: []string{},

		// A map which indicates which choices are selectedNamespace. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selectedNamespace: "",
		filter:            filter,
		views:             []string{"SHOW_NAMESPACES", "SHOW_COMMANDS"},
		commands:          []NamespaceOperation{WakeUpReviewOperation, ListPodsOperation},
		activeView:        "SHOW_NAMESPACES",
		podIteams:         []string{},
	}
}

func (m K8mNamespaceModel) Init() tea.Cmd {
	return tea.Batch(m.clusterConnectionSpinner.Tick, clusterConnect)
}

func (m K8mNamespaceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case errMsg:
		m.error = msg
		//return m, tea.Quit
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "ctrl+shift+f":
			m.filter.Focus()
		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.filteredChoices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selectedNamespace state for the item that the cursor is pointing at.
		case "enter", " ":
			switch m.activeView {
			case "SHOW_NAMESPACES":
				m.selectedNamespace = m.filteredChoices[m.cursor]
				m.activeView = "SHOW_COMMANDS"
			case "SHOW_COMMANDS":
				m.command = m.commands[m.cursor]
				if m.command.Name == WakeUpReviewOperation.Name {
					cmd := WakeUpReviewOperation.Command(m, m.selectedNamespace)
					cmds = append(cmds, cmd)
				} else {
					cmd := ListPodsOperation.Command(m, m.selectedNamespace)
					cmds = append(cmds, cmd)
				}
			}
			//cmd := WakeUpReviewOperation.Command(m, m.selectedNamespace)
			//cmds = append(cmds, cmd)
			//return m, tea.Batch(cmd)
		}
	case cursor.BlinkMsg:
	default:
		if m.filter.Value() != "" {
			m.filteredChoices = []string{}
			for _, namespace := range m.choices {
				if strings.Contains(namespace, m.filter.Value()) {
					m.filteredChoices = append(m.filteredChoices, namespace)
				}
			}
			m.cursor = 0
		}

	case podListMsg:
		m.activeView = "SHOW_PODS"
		m.podIteams = msg.items
		//return m, nil
	case clusterConnectedMsg:
		m.isConnected = true
		m.cluster.kubernetes = msg.clientset

		namespaces, err := m.cluster.kubernetes.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
		if err != nil {
			m.error = err
		}
		for _, item := range namespaces.Items {
			m.choices = append(m.choices, item.ObjectMeta.Name)
		}
		m.filteredChoices = m.choices
	case namespaceSelectedMsg:
		tea.Batch()
	}
	if !m.isConnected {
		s, cmd := m.clusterConnectionSpinner.Update(msg)
		m.clusterConnectionSpinner = s
		cmds = append(cmds, cmd)
	}

	filterValue, cmd := m.filter.Update(msg)
	m.filter = filterValue
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m K8mNamespaceModel) View() string {
	s := ""
	if !m.isConnected {
		s += m.clusterConnectionSpinner.View()
		s += "Connecting to the cluster..."
	} else {
		s = m.filter.View()
		switch m.activeView {
		case "SHOW_NAMESPACES":
			s = listNamespacesView(s, m)
		case "SHOW_COMMANDS":
			s = listCommands(s, m)
		case "SHOW_PODS":
			s = listPodsView(s, m)
		}

	}
	// The footer
	s += "\nPress q to quit.\n"
	return s
}

func listPodsView(s string, m K8mNamespaceModel) string {
	s += "\nSelect a pod\n\n"

	choicesLen := len(m.filteredChoices)
	limit := choicesLen

	if choicesLen > 10 {
		limit = 10
	}
	for i, choice := range m.podIteams[:limit] {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selectedNamespace?
		checked := " " // not selectedNamespace
		if m.selectedNamespace == choice {
			checked = "x" // selectedNamespace!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	if m.error != nil {
		s += m.error.Error()
	}

	return s
}

func listCommands(s string, m K8mNamespaceModel) string {
	s += "\nSelect a command\n\n"

	// Iterate over our choices

	for i, choice := range m.commands {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selectedNamespace?
		checked := " " // not selectedNamespace
		if m.selectedNamespace == choice.Name {
			checked = "x" // selectedNamespace!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	if m.error != nil {
		s += m.error.Error()
	}

	return s
}

func listNamespacesView(s string, m K8mNamespaceModel) string {
	s += "\nSelect a namespace\n\n"

	// Iterate over our choices

	choicesLen := len(m.filteredChoices)
	limit := choicesLen

	if choicesLen > 10 {
		limit = 10
	}
	for i, choice := range m.filteredChoices[:limit] {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selectedNamespace?
		checked := " " // not selectedNamespace
		if m.selectedNamespace == choice {
			checked = "x" // selectedNamespace!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	if m.error != nil {
		s += m.error.Error()
	}

	return s
}
