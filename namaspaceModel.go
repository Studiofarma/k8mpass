package main

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type K8mNamespaceModel struct {
	choices                  []string // items on the to-do list
	cursor                   int      // which to-do list item our cursor is pointing at
	selected                 string   // which to-do items are selected
	cluster                  kubernetesCluster
	clusterConnectionSpinner spinner.Model
	isConnected              bool
	command                  NamespaceOperation
	error                    errMsg
}

func initialModel() K8mNamespaceModel {
	return K8mNamespaceModel{
		// Our to-do list is a grocery list
		choices: []string{},

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: "",
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
		return m, tea.Quit
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			m.selected = m.choices[m.cursor]
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

	}
	if !m.isConnected {
		s, cmd := m.clusterConnectionSpinner.Update(msg)
		m.clusterConnectionSpinner = s
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m K8mNamespaceModel) View() string {
	s := ""
	if !m.isConnected {
		s += m.clusterConnectionSpinner.View()
		s += "Connecting to the cluster..."
	} else {
		s = "What should we buy at the market?\n\n"

		// Iterate over our choices
		for i, choice := range m.choices[:10] {

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

		// The footer
		s += "\nPress q to quit.\n"
	}
	return s
}
