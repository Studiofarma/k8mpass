package namespace

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
)

func (n Item) IsReviewApp() bool {
	return strings.HasPrefix(n.K8sNamespace.Name, "review")
}

func New() list.Model {
	l := list.New([]list.Item{}, ItemDelegate{}, pageWidth, pageHeight)
	l.Title = "Loading namespaces..."
	l.SetShowStatusBar(true)
	l.SetShowHelp(true)
	l.SetFilteringEnabled(true)
	l.SetShowFilter(true)
	l.Styles.Title = titleStyle
	l.SetStatusBarItemName("namespace", "namespaces")
	l.Styles.NoItems.MarginLeft(2)
	l.KeyMap.GoToEnd.Unbind()
	l.KeyMap.GoToStart.Unbind()
	l.KeyMap.ShowFullHelp.Unbind()
	l.KeyMap.CloseFullHelp.Unbind()
	l.KeyMap.CursorUp.SetHelp("↑", "up")
	l.KeyMap.CursorDown.SetHelp("↓", "down")
	l.KeyMap.NextPage.SetHelp("→", "right")
	l.KeyMap.PrevPage.SetHelp("←", "left")
	return l
}
