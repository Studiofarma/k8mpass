package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

var (
	logsTitleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.Copy().BorderStyle(b)
	}()
)

func NewViewport() viewport.Model {
	v := viewport.New(0, 0)
	v.Style = v.Style.BorderStyle(lipgloss.NormalBorder())
	return v
}

func (m PodSelectionModel) headerView(pod string) string {
	title := logsTitleStyle.Render(pod)
	line := strings.Repeat("─", max(0, m.logs.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m PodSelectionModel) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.logs.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.logs.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}
