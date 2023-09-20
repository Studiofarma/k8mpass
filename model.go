package main

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

type K8mpassModel struct {
	error                    errMsg
	cluster                  kubernetesCluster
	clusterConnectionSpinner spinner.Model
	isConnected              bool
	command                  NamespaceOperation
	namespaces               list.Model
}

type item struct {
	title, desc string
}

func (i item) FilterValue() string {
	//TODO implement me
	panic("implement me")
}

func initialModel() K8mpassModel {
	s := spinner.New()
	s.Spinner = spinner.Line

	config, err := clientcmd.BuildConfigFromFlags("", defaultKubeConfigFilePath())
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	pods, err := clientset.CoreV1().Pods("review-hack-cgmgpharm-47203-be").List(context.TODO(), metav1.ListOptions{})
	numPods := len(pods.Items)
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	items := make([]list.Item, numPods)
	for i := 0; i < numPods; i++ {
		items[i] = item{title: pods.Items[i].Name, desc: pods.Items[i].Namespace}
		fmt.Printf("POD NAME" + pods.Items[i].Name + "\n")
	}

	groceryList := list.New(items, list.NewDefaultDelegate(), 0, 0)
	groceryList.Title = "Stocazzo"
	groceryList.Styles.Title = titleStyle

	return K8mpassModel{
		clusterConnectionSpinner: s,
		command:                  WakeUpReviewOperation,
		namespaces:               groceryList,
	}
}

func (m K8mpassModel) Init() tea.Cmd {
	return tea.Batch(m.clusterConnectionSpinner.Tick, clusterConnect)
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
		}
	case clusterConnectedMsg:
		m.isConnected = true
		m.cluster.kubernetes = msg.clientset
		command := m.command.Command(m, "review-hack-cgmgpharm-47203-be")
		cmds = append(cmds, command)
	}
	if !m.isConnected {
		s, cmd := m.clusterConnectionSpinner.Update(msg)
		m.clusterConnectionSpinner = s
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m K8mpassModel) View() string {
	s := ""
	if !m.isConnected {
		s += m.clusterConnectionSpinner.View()
		s += "Connecting to the cluster..."
	} else {
		s += "Connection successful! Press esc to quit"
	}
	s += "\n"
	return s
}
