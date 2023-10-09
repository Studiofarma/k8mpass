package log

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
		return logsTitleStyle.Copy().BorderStyle(b)
	}()
)

func NewViewport() viewport.Model {
	v := viewport.New(0, 0)
	return v
}

func HeaderView(m viewport.Model, namespace string, pod string) string {
	title := logsTitleStyle.Render(fmt.Sprintf("%s : %s", namespace, pod))
	line := strings.Repeat("─", max(0, m.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func FooterView(m viewport.Model) string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}
