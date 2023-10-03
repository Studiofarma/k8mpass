package main

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	pageMaxHeight = 8
	pageWidth     = 80
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

type OperationItemDelegate struct {
	NamespaceOperation NamespaceOperation
}

func (o OperationItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(NamespaceOperation)
	if !ok {
		return
	}

	str := i.Name

	fn := lipgloss.NewStyle().PaddingLeft(2).Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return lipgloss.
				NewStyle().
				MarginLeft(2).
				Foreground(lipgloss.Color("#ffcb78")).
				Background(lipgloss.Color("#444852")).
				Render("" + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func (o OperationItemDelegate) Height() int {
	return 1
}

func (o OperationItemDelegate) Spacing() int {
	return 0
}

func (o OperationItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (no NamespaceOperation) FilterValue() string {
	return no.Name
}

func initializeOperationList() list.Model {
	var items []list.Item
	l := list.New(items, OperationItemDelegate{}, pageWidth, pageMaxHeight)
	l.Title = "Namespace operations"
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)
	l.SetShowFilter(false)
	l.Styles.Title = titleStyle
	l.Styles.NoItems.MarginLeft(2)
	l.SetStatusBarItemName("operation", "operations")
	l.StatusMessageLifetime = time.Second * 3
	additionalKeys := func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("backspace", "esc"),
				key.WithHelp("⌫/esc", "back"),
			),
		}
	}
	l.KeyMap.GoToEnd.Unbind()
	l.KeyMap.GoToStart.Unbind()
	l.KeyMap.Quit.SetKeys("ctrl+c")
	l.KeyMap.ShowFullHelp.Unbind()
	l.KeyMap.CloseFullHelp.Unbind()
	l.KeyMap.CursorUp.SetHelp("↑", "up")
	l.KeyMap.CursorDown.SetHelp("↓", "down")
	l.KeyMap.NextPage.SetHelp("→", "right")
	l.KeyMap.PrevPage.SetHelp("←", "left")
	l.AdditionalFullHelpKeys = additionalKeys
	l.AdditionalShortHelpKeys = additionalKeys
	return l
}
