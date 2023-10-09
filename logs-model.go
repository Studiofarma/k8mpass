package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	filter    textinput.Model
	filtering bool
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
	var routedCmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.view.Width = msg.Width
		m.view.Height = msg.Height - 6
	case log.UpdateLogsMsg:
		if m.logs == nil {
			break
		}
		m.view.SetContent(strings.Join(m.logs.FilterAndTruncate(m.filter.Value(), m.view.Width), "\n"))
		if m.follow {
			m.view.GotoBottom()
		}
		cmds = append(cmds, m.handler.TickUpdateLogs)
	case log.RoutedMsg:
		if batchMsg, ok := msg.Embedded.(tea.BatchMsg); ok {
			cmds = append(routedCmds, func() tea.Msg {
				return batchMsg
			})
			break
		}
		lm, lmCmd := m.view.Update(msg.Embedded)
		m.view = lm
		routedCmds = append(routedCmds, lmCmd)
		fm, fmCmd := m.filter.Update(msg.Embedded)
		m.filter = fm
		routedCmds = append(routedCmds, fmCmd)

	case tea.KeyMsg:
		if key := msg.String(); key == "f" {
			m.filtering = !m.filtering
			cmds = append(cmds, textinput.Blink)
		} else if m.filtering {
			newM, cmd := m.filter.Update(msg)
			m.filter = newM
			cmds = append(cmds, cmd)
		} else {
			newM, cmd := m.view.Update(msg)
			m.view = newM
			cmds = append(cmds, cmd)
		}
	default:
		newM, cmd := m.view.Update(msg)
		m.view = newM
		cmds = append(cmds, cmd)

		newF, cmd := m.filter.Update(msg)
		m.filter = newF
		cmds = append(cmds, cmd)
	}
	cmds = append(cmds, log.Route(routedCmds...)...)
	return m, tea.Batch(cmds...)
}

func (m LogsModel) View() string {
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(m.namespace, m.pod), m.view.View(), m.footerView())
}

func NewLogModel(service kubernetes.ILogService) LogsModel {
	text := textinput.New()
	text.PromptStyle.Width(30)
	text.Placeholder = "filtering"
	return LogsModel{
		handler:   log.NewHandler(service),
		view:      log.NewViewport(),
		follow:    true,
		filter:    text,
		filtering: false,
	}
}

func (m *LogsModel) Reset() {
	m.namespace = ""
	m.pod = ""
	m.handler.CloseLogs()
	m.logs = nil
}

func (m LogsModel) headerView(namespace string, pod string) string {
	title := log.LogsTitleStyle.Render(fmt.Sprintf("%s : %s", namespace, pod))
	filter := log.LogsTitleStyle.Render(m.filter.View())
	line := strings.Repeat("─", max(0, m.view.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, filter, line)
}

func (m LogsModel) footerView() string {
	info := log.InfoStyle.Render(fmt.Sprintf("%3.f%%", m.view.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.view.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}
