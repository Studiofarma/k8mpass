package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"sort"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/studiofarma/k8mpass/pod"
)

type PodSelectionModel struct {
	messageHandler      *pod.MessageHandler
	availableOperations []NamespaceOperation
	pods                list.Model
	operations          list.Model
	namespace           string
	dimensions          struct {
		width  int
		height int
	}
}

func (m PodSelectionModel) Init() tea.Cmd {
	return CheckConditionsThatApply(K8sCluster.kubernetes, m.namespace, m.availableOperations)
}

func (m PodSelectionModel) Update(msg tea.Msg) (PodSelectionModel, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case pod.WatchingPodsMsg:
		cmds = append(cmds, m.messageHandler.NextEvent)
	case pod.PodListMsg:
		items := make([]list.Item, len(msg.Pods))
		for i, ns := range msg.Pods {
			items[i] = ns
		}
		cmds = append(cmds, m.pods.SetItems(items))
	case pod.AddedPodMsg:
		cmds = append(cmds, m.pods.InsertItem(0, msg.Pod))
		ns := m.pods.Items()
		sort.SliceStable(ns, func(i, j int) bool {
			return ns[i].FilterValue() < ns[j].FilterValue()
		})
		cmds = append(cmds, m.pods.SetItems(ns))
		cmds = append(cmds, m.messageHandler.NextEvent)
		cmds = append(cmds, m.pods.NewStatusMessage(fmt.Sprintf("ADDED: %s", msg.Pod.K8sPod.Name)))
	case pod.RemovedPodMsg:
		var idx = pod.FindPod(m.pods.Items(), msg.Pod)
		m.pods.RemoveItem(idx)
		cmds = append(cmds, m.messageHandler.NextEvent)
		cmds = append(cmds, m.pods.NewStatusMessage(fmt.Sprintf("REMOVED: %s", msg.Pod.K8sPod.Name)))
	case pod.ModifiedPodMsg:
		var idx = pod.FindPod(m.pods.Items(), msg.Pod)
		m.pods.SetItem(idx, msg.Pod)
		cmds = append(cmds, m.messageHandler.NextEvent)
	case pod.NextEventMsg:
		cmds = append(cmds, m.messageHandler.NextEvent)
	case pod.ErrorMsg:
		m.pods.NewStatusMessage(msg.Err.Error())
	case AvailableOperationsMsg:
		var ops []list.Item
		for _, operation := range msg.operations {
			ops = append(ops, operation)
		}
		cmd := m.operations.SetItems(ops)
		m.UpdateSize()
		cmds = append(cmds, cmd)
	case noOutputResultMsg:
		var style lipgloss.Style
		if msg.success {
			style = statusMessageGreen
		} else {
			style = statusMessageRed
		}
		cmd := m.operations.NewStatusMessage(style.Render(msg.message))
		cmds = append(cmds, cmd)
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			i, ok := m.operations.SelectedItem().(NamespaceOperation)
			if ok {
				opCommand := i.Command(K8sCluster.kubernetes, m.namespace)
				cmds = append(cmds, opCommand)
			} else {
				panic("Casting went wrong")
			}
		case "backspace", "esc":
			m.pods.NewStatusMessage("")
			cmds = append(cmds, func() tea.Msg {
				return backToNamespaceSelectionMsg{}
			})
		}
	}

	lm, lmCmd := m.pods.Update(msg)
	m.pods = lm
	cmds = append(cmds, lmCmd)
	om, omCmd := m.operations.Update(msg)
	m.operations = om
	cmds = append(cmds, omCmd)

	return m, tea.Batch(cmds...)
}

func (m PodSelectionModel) View() string {
	return lipgloss.JoinVertical(
		0.0,
		m.operations.View(),
		m.pods.View(),
	)
}

func (o *PodSelectionModel) Reset() {
	o.operations.ResetSelected()
	o.pods.ResetSelected()
	o.operations.SetItems(make([]list.Item, 0))
	o.pods.SetItems(make([]list.Item, 0))
	o.messageHandler.StopWatching()
}

func (m *PodSelectionModel) UpdateSize() {
	n := len(m.operations.Items())
	emptyCorrection := 0
	if n != 0 {
		emptyCorrection = -1
	}
	opsHeight := min(8, 4+n+emptyCorrection)
	m.operations.SetHeight(opsHeight)
	m.pods.SetHeight(m.dimensions.height - opsHeight)
	m.operations.SetWidth(m.dimensions.width)
	m.pods.SetWidth(m.dimensions.width)
}
