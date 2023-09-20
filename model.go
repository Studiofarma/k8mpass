package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"io"
	"strings"
)

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

type K8mpassModel struct {
	error                    errMsg
	cluster                  kubernetesCluster
	clusterConnectionSpinner spinner.Model
	isConnected              bool
	namespacePodsInfo        namespacePodsInfo
	command                  NamespaceOperation
	list                     list.Model
	listItems                []list.Item
	nameSpace                string
}

type item struct {
	name string
}

func (i item) Title() string {
	return i.name
}
func (i item) FilterValue() string {
	return i.name
}

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	mItem, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, mItem.name)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func initialModel() K8mpassModel {
	s := spinner.New()
	s.Spinner = spinner.Line
	model := K8mpassModel{
		clusterConnectionSpinner: s,
		namespacePodsInfo: namespacePodsInfo{
			podsInfo: []podInfo{},
		},
		command: GetAllNamespacesOperation,
		list:    list.New([]list.Item{}, itemDelegate{}, 80, 15),
	}
	return model
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
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			command := m.command.Command(m, m.list.SelectedItem().FilterValue())
			cmds = append(cmds, command)
		}
	case nameSpaceSelectedMsg:
		m.nameSpace = msg.body
		operations := []list.Item{item{"op1"}, item{"op2"}}
		m.list = list.New(operations, itemDelegate{}, 80, 15)
		m.list.Title = "Select the operation"
		m.command = OperationSelected
	case operationSelectedMsg:
		fmt.Printf("%s", msg.body)

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case clusterConnectedMsg:
		m.isConnected = true
		m.cluster.kubernetes = msg.clientset
		command := m.command.Command(m, "")
		cmds = append(cmds, command)
	case namespacesNamesMsg:
		namespacesList := []list.Item{}
		for _, val := range msg.body {
			namespacesList = append(namespacesList, item{val})
		}
		m.list = list.New(namespacesList, itemDelegate{}, 80, 15)
		m.list.Title = "Select the namespace"
		m.command = NameSpaceSelected
	case podsInfoMsg:
		m.namespacePodsInfo = msg.body
		m.namespacePodsInfo.calculateNamespaceStatus()
		fmt.Println(m)
	}
	if !m.isConnected {
		s, cmd := m.clusterConnectionSpinner.Update(msg)
		m.clusterConnectionSpinner = s
		cmds = append(cmds, cmd)
	} else {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
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
		return docStyle.Render(m.list.View())
		// s += "Connection successful! Press esc to quit"
	}
	s += "\n"
	return s
}
