package main

import (
	"fmt"
	"sort"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/studiofarma/k8mpass/pod"
)

type PodSelectionModel struct {
	messageHandler *pod.MessageHandler
	pods           list.Model
}

func (m PodSelectionModel) Init() tea.Cmd {
	return nil
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
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
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

	return m, tea.Batch(cmds...)
}

func (m PodSelectionModel) View() string {
	return m.pods.View()
}
