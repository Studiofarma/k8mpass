package namespace

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"time"
)

func New() list.Model {
	l := list.New([]list.Item{}, ItemDelegate{}, pageWidth, pageHeight)
	l.Title = "Loading namespaces..."
	l.SetShowStatusBar(true)
	l.SetShowHelp(true)
	l.SetFilteringEnabled(true)
	l.SetShowFilter(true)
	l.Styles.Title = titleStyle
	l.SetStatusBarItemName("namespace", "namespaces")
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
				key.WithKeys("p"),
				key.WithHelp("p", "pin"),
			),
			key.NewBinding(
				key.WithKeys("u"),
				key.WithHelp("u", "unpin"),
			),
		}
	}

	l.AdditionalFullHelpKeys = additionalKeys
	l.AdditionalShortHelpKeys = additionalKeys
	return l
}
