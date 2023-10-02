package pod

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
)

func New() list.Model {
	l := list.New([]list.Item{}, ItemDelegate{}, pageWidth, pageHeight)
	l.Title = "Pods"
	l.SetShowStatusBar(true)
	l.SetShowHelp(true)
	l.Styles.Title = titleStyle
	l.SetStatusBarItemName("pod", "pods")
	l.Styles.NoItems.MarginLeft(2)
	l.KeyMap.GoToEnd.Unbind()
	l.KeyMap.Quit.SetKeys("q", "ctrl+c")
	l.KeyMap.GoToStart.Unbind()
	l.KeyMap.ShowFullHelp.Unbind()
	l.KeyMap.CloseFullHelp.Unbind()
	l.KeyMap.CursorUp.SetHelp("↑", "up")
	l.KeyMap.CursorDown.SetHelp("↓", "down")
	l.KeyMap.NextPage.SetHelp("→", "right")
	l.KeyMap.PrevPage.SetHelp("←", "left")
	additionalKeys := func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("backspace", "esc"),
				key.WithHelp("⌫/esc", "back"),
			),
		}
	}

	l.AdditionalFullHelpKeys = additionalKeys
	l.AdditionalShortHelpKeys = additionalKeys
	return l
}
