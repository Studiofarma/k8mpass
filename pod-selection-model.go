package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/studiofarma/k8mpass/api"
	"sort"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/studiofarma/k8mpass/pod"
)

type PodSelectionModel struct {
	messageHandler *pod.MessageHandler
	pods           list.Model
	operations     list.Model
	namespace      string
	focus          focus
	logs           viewport.Model
	dimensions     struct {
		width  int
		height int
	}
}

type focus int8

const (
	operations focus = 0
	pods       focus = 1
	logs       focus = 2
)

func (m PodSelectionModel) Init() tea.Cmd {
	return nil
}

func (m PodSelectionModel) Update(msg tea.Msg) (PodSelectionModel, tea.Cmd) {
	var cmds []tea.Cmd
	var routedCmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.dimensions = struct {
			width  int
			height int
		}{width: msg.Width, height: msg.Height}
		m.UpdateSize()
	case namespaceSelectedMsg:
		m.namespace = msg.namespace
		m.operations.Title = msg.namespace
		cmds = append(cmds, m.messageHandler.CheckConditionsThatApply(msg.namespace))
		cmds = append(cmds, m.messageHandler.GetPods(msg.namespace))
		routedCmds = append(cmds, m.operations.StartSpinner())
	case pod.WatchingPodsMsg:
		cmds = append(cmds, m.messageHandler.NextEvent)
	case pod.ListMsg:
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
		//cmds = append(cmds, m.pods.NewStatusMessage(fmt.Sprintf("ADDED: %s", msg.Pod.K8sPod.Name)))
	case pod.RemovedPodMsg:
		var idx = pod.FindPod(m.pods.Items(), msg.Pod)
		m.pods.RemoveItem(idx)
		cmds = append(cmds, m.messageHandler.NextEvent)
		//cmds = append(cmds, m.pods.NewStatusMessage(fmt.Sprintf("REMOVED: %s", msg.Pod.K8sPod.Name)))
	case pod.ModifiedPodMsg:
		var idx = pod.FindPod(m.pods.Items(), msg.Pod)
		if idx < 0 {
			cmds = append(cmds, m.messageHandler.NextEvent)
			break
		}
		cmds = append(cmds, m.pods.SetItem(idx, msg.Pod))
		cmds = append(cmds, m.messageHandler.NextEvent)
	case pod.NextEventMsg:
		cmds = append(cmds, m.messageHandler.NextEvent)
	case pod.ErrorMsg:
		m.pods.NewStatusMessage(msg.Err.Error())
	case api.AvailableOperationsMsg:
		var ops []list.Item
		for _, operation := range msg.Operations {
			ops = append(ops, operation)
		}
		cmd := m.operations.SetItems(ops)
		cmds = append(cmds, cmd)
		m.operations.StopSpinner()
		if len(msg.Operations) == 0 {
			m.focus = pods
			m.pods.SetDelegate(pod.ItemDelegate{IsFocused: true})
			m.operations.SetDelegate(OperationItemDelegate{IsFocused: false})
		}
	case api.NoOutputResultMsg:
		var style lipgloss.Style
		if msg.Success {
			style = statusMessageGreen
		} else {
			style = statusMessageRed
		}
		cmd := m.operations.NewStatusMessage(style.Render(msg.Message))
		cmds = append(cmds, cmd)
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			switch m.focus {
			case operations:
				i, ok := m.operations.SelectedItem().(api.NamespaceOperation)
				if ok {
					opCommand := m.messageHandler.RunComand(i, m.namespace)
					cmds = append(cmds, opCommand)
				} else {
					panic("Casting went wrong")
				}
			case pods:
				m.focus = logs
				m.logs.SetContent(m.messageHandler.GetLogs(m.namespace, m.pods.SelectedItem().FilterValue(), m.dimensions.width))
				routedCmds = append(routedCmds, viewport.Sync(m.logs))
			}
		case "backspace", "esc":
			switch m.focus {
			case logs:
				m.focus = pods
			default:
				m.pods.NewStatusMessage("")
				cmds = append(cmds, func() tea.Msg {
					return backToNamespaceSelectionMsg{}
				})
			}
		case "tab":
			switch m.focus {
			case operations:
				m.focus = pods
				m.pods.SetDelegate(pod.ItemDelegate{IsFocused: true})
				m.operations.SetDelegate(OperationItemDelegate{IsFocused: false})
			case pods:
				if len(m.operations.Items()) == 0 {
					break
				}
				m.focus = operations
				m.pods.SetDelegate(pod.ItemDelegate{IsFocused: false})
				m.operations.SetDelegate(OperationItemDelegate{IsFocused: true})
			}
		default:
			switch m.focus {
			case pods:
				lm, lmCmd := m.pods.Update(msg)
				m.pods = lm
				routedCmds = append(cmds, lmCmd)

				cmds = append(cmds, pod.Route(routedCmds...)...)
				return m, tea.Batch(cmds...)
			case operations:
				om, omCmd := m.operations.Update(msg)
				m.operations = om
				routedCmds = append(cmds, omCmd)

				cmds = append(cmds, pod.Route(routedCmds...)...)
				return m, tea.Batch(cmds...)
			case logs:
				logm, logCmd := m.logs.Update(msg)
				m.logs = logm
				routedCmds = append(cmds, logCmd)
				cmds = append(cmds, pod.Route(routedCmds...)...)
				return m, tea.Batch(cmds...)
			}
		}
	}

	lm, lmCmd := m.pods.Update(msg)
	m.pods = lm
	cmds = append(cmds, lmCmd)
	om, omCmd := m.operations.Update(msg)
	m.operations = om
	cmds = append(cmds, omCmd)
	cmds = append(cmds, pod.Route(routedCmds...)...)
	return m, tea.Batch(cmds...)
}

func (m PodSelectionModel) View() string {
	switch m.focus {
	case logs:
		return fmt.Sprintf("%s\n%s\n%s", m.headerView(m.pods.SelectedItem().FilterValue()), m.logs.View(), m.footerView())
	default:
		return lipgloss.JoinVertical(
			0.0,
			m.operations.View(),
			m.pods.View(),
		)
	}
}

func (m *PodSelectionModel) Reset() {
	m.operations.ResetSelected()
	m.pods.ResetSelected()
	m.operations = initializeOperationList()
	m.pods = pod.New()
	m.focus = operations
	m.UpdateSize()
	m.messageHandler.StopWatching()
}

func (m *PodSelectionModel) UpdateSize() {
	opsHeight := 8
	m.operations.SetHeight(opsHeight)
	m.pods.SetHeight(m.dimensions.height - opsHeight)
	m.operations.SetWidth(m.dimensions.width)
	m.pods.SetWidth(m.dimensions.width)
	m.logs.Width = m.dimensions.width
	m.logs.Height = m.dimensions.height - 6
}
