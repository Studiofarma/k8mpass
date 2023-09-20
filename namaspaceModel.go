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
	selected                 string // which to-do items are selected
	cluster                  kubernetesCluster
	clusterConnectionSpinner spinner.Model
	isConnected              bool
	command                  NamespaceOperation
	error                    errMsg
	filter                   textinput.Model
}

func initialModel() K8mNamespaceModel {
	filter := textinput.New()
	filter.Focus()
	filter.Placeholder = "Search"

	return K8mNamespaceModel{
		// Our to-do list is a grocery list
		choices:         []string{},
		filteredChoices: []string{},

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: "",
		filter:   filter,
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
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			m.selected = m.filteredChoices[m.cursor]
			cmd := WakeUpReviewOperation.Command(m, m.selected)
			cmds = append(cmds, cmd)
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

			// Is this choice selected?
			checked := " " // not selected
			if m.selected == choice {
				checked = "x" // selected!
			}

			// Render the row
			s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
		}

		if m.error != nil {
			s += m.error.Error()
		}

		// The footer
		s += "\nPress q to quit.\n"
	}
	return s
}
