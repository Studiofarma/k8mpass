package main

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
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
		}
	case clusterConnectedMsg:
		m.isConnected = true
		m.cluster.kubernetes = msg.clientset

		namespaces, err := m.cluster.kubernetes.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
		if err != nil {
			fmt.Println(err)
		}
		for _, item := range namespaces.Items {
			m.choices = append(m.choices, item.ObjectMeta.Name)
		}

		fmt.Println(m.choices)

		//command := m.command.Command(m, "review-devops-new-filldata")
		//cmds = append(cmds, command)
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
		s += "Connection successful! Press esc to quit"
	}
	return s
}

func ListNameSpaces(coreClient kubernetes.Interface) []string {
	nsList, err := coreClient.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{}) //checkErr(err)
	fmt.Println(err)

	for _, n := range nsList.Items {
		fmt.Println(n.Name)
	}

	return []string{"a", "b"}
}
