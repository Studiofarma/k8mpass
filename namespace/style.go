package namespace

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
	customPropertiesStyle = lipgloss.NewStyle().
				PaddingLeft(4).
				Foreground(lipgloss.Color("#7d7d7d"))
	selectedItemStyle = lipgloss.NewStyle().
				MarginLeft(2)
	unselectedItemStyle = lipgloss.NewStyle().
				MarginLeft(2)
	pinnedStyle = lipgloss.NewStyle().
			MarginLeft(2).
			Foreground(lipgloss.Color("#fac65f"))

	terminatingNamespace = lipgloss.NewStyle().
				MarginLeft(2).
				Foreground(lipgloss.Color("#6b6b6b"))
)
