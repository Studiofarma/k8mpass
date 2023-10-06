package pod

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"time"
)

func New() list.Model {
	l := list.New([]list.Item{}, ItemDelegate{}, 0, 0)
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(true)
	l.Title = "Pods"
	l.SetShowTitle(true)
	l.Styles.Title = titleStyle
	l.SetStatusBarItemName("pod", "pods")
	l.StatusMessageLifetime = 3 * time.Second
	l.Styles.NoItems.MarginLeft(2)
	l.KeyMap.GoToEnd.Unbind()
	l.KeyMap.Quit.SetKeys("ctrl+c")
	l.KeyMap.Quit.SetHelp("ctrl+c", "quit")
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
				key.WithHelp("backspace/esc", "back"),
			),
			key.NewBinding(
				key.WithKeys("tab"),
				key.WithHelp("tab", "switch focus"),
			),
		}
	}

	l.AdditionalFullHelpKeys = additionalKeys
	l.AdditionalShortHelpKeys = additionalKeys
	return l
}
