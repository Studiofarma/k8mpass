package main

import "github.com/charmbracelet/lipgloss"

var (
	pageHeight = 20
	pageWidth  = 80
)

var (
	titleStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("62")).
		Bold(true).
		Padding(0, 1)
	statusMessageGreen = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#66ffc2"))
	statusMessageRed = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ff6666"))
)
