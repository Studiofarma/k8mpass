package log

import (
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

var (
	LogsTitleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	InfoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return LogsTitleStyle.Copy().BorderStyle(b)
	}()
)

func NewViewport() viewport.Model {
	v := viewport.New(0, 0)
	return v
}
