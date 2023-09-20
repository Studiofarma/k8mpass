package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

package main

import (
"context"
"fmt"
"github.com/charmbracelet/bubbles/key"
"github.com/charmbracelet/bubbles/list"
"github.com/charmbracelet/bubbles/spinner"
tea "github.com/charmbracelet/bubbletea"
metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
"k8s.io/client-go/kubernetes"
"k8s.io/client-go/tools/clientcmd"
)
type sessionState int

const (
	mainView       sessionState = 0
	namespacesView sessionState = 1
	podsView       sessionState = 2
	cronjobsView                = 3
)

type K8mpassModel struct {
	state                    sessionState
	entry tea.Model
	error                    errMsg
	cluster                  kubernetesCluster
	clusterConnectionSpinner spinner.Model
	isConnected              bool
	command                  NamespaceOperation
	namespaces               list.Model
	keys                     *listKeyMap
	delegateKeys             *delegateKeyMap
}



func initialModel() K8mpassModel {
	var (
		delegateKeys = newDelegateKeyMap()
		listKeys     = newListKeyMap()
	)
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
		items[i] = item{title: pods.Items[i].Name}
		//fmt.Printf("POD NAME" + pods.Items[i].Name + "\n")
	}

	delegate := newItemDelegate(newDelegateKeyMap())
	groceryList := list.New(items, delegate, 80, 15)
	groceryList.Title = "Stocazzo"
	groceryList.Styles.Title = titleStyle
	groceryList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.toggleSpinner,
			listKeys.insertItem,
			listKeys.toggleTitleBar,
			listKeys.toggleStatusBar,
			listKeys.togglePagination,
			listKeys.toggleHelpMenu,
		}
	}

	return K8mpassModel{
		state: mainView,
		entry: s,
		clusterConnectionSpinner: s,
		command:                  WakeUpReviewOperation,
		namespaces:               groceryList,
		keys:                     listKeys,
		delegateKeys:             delegateKeys,
	}
}

func (m K8mpassModel) Init() tea.Cmd {
	return tea.EnterAltScreen
	//return tea.Batch(m.clusterConnectionSpinner.Tick, clusterConnect)
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
		if m.namespaces.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.toggleSpinner):
			cmd := m.namespaces.ToggleSpinner()
			return m, cmd

		case key.Matches(msg, m.keys.toggleTitleBar):
			v := !m.namespaces.ShowTitle()
			m.namespaces.SetShowTitle(v)
			m.namespaces.SetShowFilter(v)
			m.namespaces.SetFilteringEnabled(v)
			return m, nil

		case key.Matches(msg, m.keys.toggleStatusBar):
			m.namespaces.SetShowStatusBar(!m.namespaces.ShowStatusBar())
			return m, nil

		case key.Matches(msg, m.keys.togglePagination):
			m.namespaces.SetShowPagination(!m.namespaces.ShowPagination())
			return m, nil

		case key.Matches(msg, m.keys.toggleHelpMenu):
			m.namespaces.SetShowHelp(!m.namespaces.ShowHelp())
			return m, nil

			//case key.Matches(msg, m.keys.insertItem):
			//	m.delegateKeys.remove.SetEnabled(true)
			//	newItem := m..next()
			//	insCmd := m.namespaces.InsertItem(0, newItem)
			//	statusCmd := m.namespaces.NewStatusMessage(statusMessageStyle("Added " + newItem.Title()))
			//	return m, tea.Batch(insCmd, statusCmd)
			//
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

	newListModel, cmd := m.namespaces.Update(msg)
	m.namespaces = newListModel
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m K8mpassModel) View() string {
	switch m.state {
	case mainView:
		return appStyle.Render()
	}
	return appStyle.Render(m.namespaces.View())
	//return s
}
