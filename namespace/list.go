package namespace

import (
	"github.com/charmbracelet/bubbles/list"
)

func New(pinned []string) list.Model {
	l := list.New([]list.Item{}, ItemDelegate{Pinned: pinned}, 0, 0)
	l.Title = "Loading namespaces..."
	l.SetShowStatusBar(true)
	l.SetShowHelp(true)
	l.SetFilteringEnabled(true)
	l.SetShowFilter(true)
	l.Styles.Title = titleStyle
	l.SetStatusBarItemName("namespace", "namespaces")
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
	return l
}
