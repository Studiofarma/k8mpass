package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"io"
	v1 "k8s.io/api/core/v1"
	"math/rand"
	"strings"
	"time"
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

type listKeyMap struct {
	toggleSpinner    key.Binding
	toggleTitleBar   key.Binding
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
	insertItem       key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		toggleSpinner: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "toggle spinner"),
		),
		toggleTitleBar: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "toggle title"),
		),
		toggleStatusBar: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "toggle status"),
		),
		togglePagination: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle pagination"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
	}
}

type delegateKeyMap struct {
	choose key.Binding
	remove key.Binding
}

type K8mpassModel struct {
	error                    errMsg
	cluster                  kubernetesCluster
	clusterConnectionSpinner spinner.Model
	isConnected              bool
	command                  NamespaceOperation
	namespaces               list.Model
	keys                     *listKeyMap
	delegateKeys             *delegateKeyMap
}

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }
func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "choose"),
		),
	}
}

func getPodLogs(nameSpace string, podName string) []list.Item {
	podLogOpts := v1.PodLogOptions{}
	config := getConfig()
	// creates the clientset
	clientset := createClientSet(config)

	req := clientset.CoreV1().Pods(nameSpace).GetLogs(podName, &podLogOpts)
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		panic("error in opening stream")
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		panic("error in copy information from podLogs to buf")
	}
	str := buf.String()

	string_split := strings.Split(str, "\n")
	to_return := make([]list.Item, len(string_split))
	for i := 0; i < len(string_split); i++ {
		to_return[i] = item{title: string_split[i]}
	}

	return to_return
}

func namespaceItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var title string

		if i, ok := m.SelectedItem().(item); ok {
			title = i.Title()
		} else {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.choose):
				podList := getPods(createClientSet(getConfig()), title)

				podNames := make([]list.Item, len(podList.Items))
				for i := 0; i < len(podList.Items); i++ {
					podNames[i] = item{title: podList.Items[i].Name}
				}

				m.SetDelegate(podItemDelegate(newDelegateKeyMap(), title))
				m.Title = "Selected namespace: " + title + " pod to see logs"
				return m.SetItems(podNames)

			case key.Matches(msg, keys.remove):
				index := m.Index()
				m.RemoveItem(index)
				if len(m.Items()) == 0 {
					keys.remove.SetEnabled(false)
				}
				return m.NewStatusMessage(statusMessageStyle("Deleted " + title))
			}
		}

		return nil
	}

	help := []key.Binding{keys.choose, keys.remove}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}

func podItemDelegate(keys *delegateKeyMap, namespace string) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var title string

		if i, ok := m.SelectedItem().(item); ok {
			title = i.Title()
		} else {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.choose):
				return m.SetItems(getPodLogs(namespace, title))

			case key.Matches(msg, keys.remove):
				index := m.Index()
				m.RemoveItem(index)
				if len(m.Items()) == 0 {
					keys.remove.SetEnabled(false)
				}
				return m.NewStatusMessage(statusMessageStyle("Deleted " + title))
			}
		}

		return nil
	}

	help := []key.Binding{keys.choose, keys.remove}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}

func initialModel() K8mpassModel {
	var (
		delegateKeys = newDelegateKeyMap()
		listKeys     = newListKeyMap()
	)
	s := spinner.New()
	s.Spinner = spinner.Line

	delegate := namespaceItemDelegate(newDelegateKeyMap())
	groceryList := list.New(getNamespaces(createClientSet(getConfig())), delegate, 80, 15)
	groceryList.Title = "----"
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
		clusterConnectionSpinner: s,
		command:                  WakeUpReviewOperation,
		namespaces:               groceryList,
		keys:                     listKeys,
		delegateKeys:             delegateKeys,
	}
}

func setPodList() []list.Item {
	clientSet := createClientSet(getConfig())
	pods := getPods(clientSet, "review-hack-cgmgpharm-47203-be")
	numPods := len(pods.Items)
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	items := make([]list.Item, numPods)
	for i := 0; i < numPods; i++ {
		items[i] = item{title: pods.Items[i].Name}
	}
	return items
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
	rand.Seed(time.Now().UTC().UnixNano())
	/*s := ""
	if !m.isConnected {
		s += m.clusterConnectionSpinner.View()
		s += "Connecting to the cluster..."
	} else {
		s += "Connection successful! Press esc to quit"
	}
	s += "\n"
	s += appStyle.Render(m.namespaces.View())*/
	return appStyle.Render(m.namespaces.View())
	//return s
}
