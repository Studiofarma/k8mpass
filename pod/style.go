package pod

import (
	"github.com/charmbracelet/lipgloss"
	v1 "k8s.io/api/core/v1"
)

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
				PaddingLeft(2).
				Foreground(lipgloss.Color("#7d7d7d"))
	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("170"))
	unselectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(4)
)

func styleString(s string, style lipgloss.Style) lipgloss.Style {
	return style.SetString(s)
}

func podStyle(status v1.PodStatus) lipgloss.Style {

	switch status.Phase {
	case v1.PodRunning:
		var ready = true
		for _, c := range status.ContainerStatuses {
			ready = ready && c.Ready
		}
		if !ready {
			return lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ff6666"))
		}
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#66ffc2"))
	case v1.PodFailed:
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff6666"))
	case v1.PodPending:
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#fcaf49"))
	default:
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#a6a6a6"))
	}

}
