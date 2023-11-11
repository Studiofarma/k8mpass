package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
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
	help      help.Model
	filtering bool
	wrap      bool
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
		m.view.SetContent(m.content())
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
		if key := msg.String(); key == "esc" {
			m.filtering = false
			m.filter.Blur()
			m.filter.Reset()
			m.view.SetContent(m.content())
		} else if key == "enter" {
			m.filtering = false
			m.filter.Blur()
		} else if m.filtering {
			newM, cmd := m.filter.Update(msg)
			m.filter = newM
			routedCmds = append(routedCmds, cmd)
			m.view.SetContent(m.content())
		} else if key := msg.String(); key == "/" {
			m.filtering = true
			routedCmds = append(routedCmds, m.filter.Focus())
		} else if key == "f" {
			m.follow = !m.follow
		} else if key == "w" {
			m.wrap = !m.wrap
			m.view.SetContent(m.content())
		} else {
			newM, cmd := m.view.Update(msg)
			m.view = newM
			routedCmds = append(routedCmds, cmd)
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

func (m LogsModel) content() string {
	if m.wrap {
		return strings.Join(m.logs.FilterAndWrap(m.filter.Value(), m.view.Width), "\n")
	} else {
		return strings.Join(m.logs.FilterAndTruncate(m.filter.Value(), m.view.Width), "\n")
	}
}

func (m LogsModel) View() string {
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(m.namespace, m.pod), m.view.View(), m.footerView())
}

func NewLogModel(service kubernetes.ILogService) LogsModel {
	text := textinput.New()
	text.PromptStyle.Width(30)
	text.Placeholder = "Filter..."
	h := help.New()
	return LogsModel{
		handler:   log.NewHandler(service),
		view:      log.NewViewport(),
		follow:    true,
		filter:    text,
		filtering: false,
		help:      h,
		wrap:      false,
	}
}

func (m *LogsModel) Reset() {
	m.namespace = ""
	m.pod = ""
	m.handler.CloseLogs()
	m.logs = nil
	m.filter.Reset()
}

func (m LogsModel) headerView(namespace string, pod string) string {
	title := log.LogsTitleStyle.Render(fmt.Sprintf("%s : %s", namespace, pod))
	filter := log.LogsTitleStyle.Render(m.filter.View())
	line := strings.Repeat("─", max(0, m.view.Width-lipgloss.Width(title+filter)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, filter, line)
}

func (m LogsModel) footerView() string {
	var following string
	if m.follow {
		following = "Y"
	} else {
		following = "N"
	}
	keys := []key.Binding{
		key.NewBinding(
			key.WithKeys("f"),
			key.WithHelp("f", "auto follow"),
		),
		key.NewBinding(
			key.WithKeys("w"),
			key.WithHelp("w", "wrap"),
		),
		key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "filter"),
		),
		key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "confirm"),
		),
	}
	info := log.InfoStyle.Render(fmt.Sprintf("Following:%v %3.f%%", following, m.view.ScrollPercent()*100))
	h := log.InfoStyle.Render(m.help.ShortHelpView(keys))
	line := strings.Repeat("─", max(0, m.view.Width-lipgloss.Width(info))-lipgloss.Width(h))
	return lipgloss.JoinHorizontal(lipgloss.Center, h, line, info)
}
