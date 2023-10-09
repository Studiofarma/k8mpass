package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/studiofarma/k8mpass/kubernetes"
	"github.com/studiofarma/k8mpass/log"
	"strings"
)

type LogsModel struct {
	namespace string
	pod       string
	handler   *log.MessageHandler
	view      viewport.Model
	logs      *log.Logs
	follow    bool
}

func (m LogsModel) Init() tea.Cmd {
	return tea.Sequence(
		func() tea.Msg {
			m.handler.FollowLogs(m.namespace, m.pod, m.logs)
			return nil
		},
		func() tea.Msg {
			return log.UpdateLogsMsg{}
		},
	)
}

func (m LogsModel) Update(msg tea.Msg) (LogsModel, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.view.Width = msg.Width
		m.view.Height = msg.Height - 6
	case log.UpdateLogsMsg:
		if m.logs == nil {
			break
		}
		m.view.SetContent(strings.Join(m.logs.TruncatedLines(m.view.Width), "\n"))
		if m.follow {
			m.view.GotoBottom()
		}
		cmds = append(cmds, m.handler.TickUpdateLogs)
	case tea.KeyMsg:
		if key := msg.String(); key == "f" {
			m.follow = !m.follow
			if m.follow {
				m.view.GotoBottom()
			}
		} else {
			newM, cmd := m.view.Update(msg)
			m.view = newM
			cmds = append(cmds, cmd)
		}
	default:
		newM, cmd := m.view.Update(msg)
		m.view = newM
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m LogsModel) View() string {
	return fmt.Sprintf("%s\n%s\n%s", log.HeaderView(m.view, m.namespace, m.pod), m.view.View(), log.FooterView(m.view))
}

func NewLogModel(service kubernetes.ILogService) LogsModel {
	return LogsModel{
		handler: log.NewHandler(service),
		view:    log.NewViewport(),
		follow:  true,
	}
}

func (m *LogsModel) Reset() {
	m.namespace = ""
	m.pod = ""
	m.handler.CloseLogs()
	m.logs = nil
}
