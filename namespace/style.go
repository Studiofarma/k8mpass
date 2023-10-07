package namespace

import "github.com/charmbracelet/lipgloss"

var (
	pageHeight       = 20
	pageWidth        = 80
	nameMaxLength    = 30
	propertyMaxWidth = 15
)

var (
	titleStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("62")).
			Bold(true).
			Padding(0, 1)
	commonStyle = lipgloss.NewStyle().
			MarginLeft(2).
			Width(nameMaxLength)
	selectedItemStyle = lipgloss.NewStyle().
				Background(slightlyBrighterTerminalColor())
	pinnedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffe39c"))

	customPropertiesStyle = lipgloss.NewStyle().
				PaddingLeft(4).
				Foreground(lipgloss.Color("#707070")).
				Width(propertyMaxWidth)

	terminatingNamespaceStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color("#6b6b6b"))
)
