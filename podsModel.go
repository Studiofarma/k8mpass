package main

//
//import (
//	"context"
//	"fmt"
//	"github.com/charmbracelet/bubbles/key"
//	"github.com/charmbracelet/bubbles/list"
//	"github.com/charmbracelet/bubbles/spinner"
//	tea "github.com/charmbracelet/bubbletea"
//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//	"k8s.io/client-go/kubernetes"
//)
//
//type podsModel struct {
//	error        errMsg
//	namespace    string
//	isConnected  bool
//	command      NamespaceOperation
//	keys         *listKeyMap
//	delegateKeys *delegateKeyMap
//	clientset    kubernetes.Clientset
//}
//
//func initialModel(m podsModel) podsModel {
//	var (
//		delegateKeys = newDelegateKeyMap()
//		listKeys     = newListKeyMap()
//	)
//	s := spinner.New()
//	s.Spinner = spinner.Line
//
//	//config, err := clientcmd.BuildConfigFromFlags("", defaultKubeConfigFilePath())
//	//if err != nil {
//	//	panic(err.Error())
//	//}
//	//clientset, err := kubernetes.NewForConfig(config)
//	pods, err := m.clientset.CoreV1().Pods(m.namespace).List(context.TODO(), metav1.ListOptions{})
//	numPods := len(pods.Items)
//	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
//
//	items := make([]list.Item, numPods)
//	for i := 0; i < numPods; i++ {
//		items[i] = item{title: pods.Items[i].Name}
//		//fmt.Printf("POD NAME" + pods.Items[i].Name + "\n")
//	}
//
//	delegate := newItemDelegate(newDelegateKeyMap())
//	groceryList := list.New(items, delegate, 80, 15)
//	groceryList.Title = "Pods"
//	groceryList.Styles.Title = titleStyle
//	groceryList.AdditionalFullHelpKeys = func() []key.Binding {
//		return []key.Binding{
//			listKeys.toggleSpinner,
//			listKeys.insertItem,
//			listKeys.toggleTitleBar,
//			listKeys.toggleStatusBar,
//			listKeys.togglePagination,
//			listKeys.toggleHelpMenu,
//		}
//	}
//
//	return K8mpassModel{
//		clusterConnectionSpinner: s,
//		command:                  WakeUpReviewOperation,
//		namespaces:               groceryList,
//		keys:                     listKeys,
//		delegateKeys:             delegateKeys,
//	}
//}
//
//func (m K8mpassModel) Init() tea.Cmd {
//	return tea.EnterAltScreen
//	//return tea.Batch(m.clusterConnectionSpinner.Tick, clusterConnect)
//}
//
//func (m K8mpassModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
//	var cmds []tea.Cmd
//	switch msg := msg.(type) {
//	case errMsg:
//		m.error = msg
//		return m, tea.Quit
//	case tea.KeyMsg:
//		switch msg.String() {
//		case "ctrl+c", "esc":
//			return m, tea.Quit
//		}
//		if m.namespaces.FilterState() == list.Filtering {
//			break
//		}
//
//		switch {
//		case key.Matches(msg, m.keys.toggleSpinner):
//			cmd := m.namespaces.ToggleSpinner()
//			return m, cmd
//
//		case key.Matches(msg, m.keys.toggleTitleBar):
//			v := !m.namespaces.ShowTitle()
//			m.namespaces.SetShowTitle(v)
//			m.namespaces.SetShowFilter(v)
//			m.namespaces.SetFilteringEnabled(v)
//			return m, nil
//
//		case key.Matches(msg, m.keys.toggleStatusBar):
//			m.namespaces.SetShowStatusBar(!m.namespaces.ShowStatusBar())
//			return m, nil
//
//		case key.Matches(msg, m.keys.togglePagination):
//			m.namespaces.SetShowPagination(!m.namespaces.ShowPagination())
//			return m, nil
//
//		case key.Matches(msg, m.keys.toggleHelpMenu):
//			m.namespaces.SetShowHelp(!m.namespaces.ShowHelp())
//			return m, nil
//
//			//case key.Matches(msg, m.keys.insertItem):
//			//	m.delegateKeys.remove.SetEnabled(true)
//			//	newItem := m..next()
//			//	insCmd := m.namespaces.InsertItem(0, newItem)
//			//	statusCmd := m.namespaces.NewStatusMessage(statusMessageStyle("Added " + newItem.Title()))
//			//	return m, tea.Batch(insCmd, statusCmd)
//			//
//		}
//	case clusterConnectedMsg:
//		m.isConnected = true
//		m.cluster.kubernetes = msg.clientset
//		command := m.command.Command(m, "review-hack-cgmgpharm-47203-be")
//		cmds = append(cmds, command)
//	}
//
//	if !m.isConnected {
//		s, cmd := m.clusterConnectionSpinner.Update(msg)
//		m.clusterConnectionSpinner = s
//		cmds = append(cmds, cmd)
//	}
//
//	newListModel, cmd := m.namespaces.Update(msg)
//	m.namespaces = newListModel
//	cmds = append(cmds, cmd)
//	return m, tea.Batch(cmds...)
//}
//
//func (m K8mpassModel) View() string {
//	/*s := ""
//	if !m.isConnected {
//		s += m.clusterConnectionSpinner.View()
//		s += "Connecting to the cluster..."
//	} else {
//		s += "Connection successful! Press esc to quit"
//	}
//	s += "\n"
//	s += appStyle.Render(m.namespaces.View())*/
//	return appStyle.Render(m.namespaces.View())
//	//return s
//}
//func setPodList() []list.Item {
//	clientSet := createClientSet(getConfig())
//	pods := getPods(clientSet, "review-hack-cgmgpharm-47203-be")
//	numPods := len(pods.Items)
//	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
//
//	items := make([]list.Item, numPods)
//	for i := 0; i < numPods; i++ {
//		items[i] = item{title: pods.Items[i].Name}
//	}
//	return items
//}
