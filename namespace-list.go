package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"io"
	v1 "k8s.io/api/core/v1"
	"strings"
)

type NamespaceItem v1.Namespace

type NamespaceItemDelegate struct{}

func (n NamespaceItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(NamespaceItem)
	if !ok {
		return
	}

	str := i.Name

	fn := lipgloss.NewStyle().PaddingLeft(4).Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170")).Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func (n NamespaceItemDelegate) Height() int {
	return 1
}

func (n NamespaceItemDelegate) Spacing() int {
	return 0
}

func (n NamespaceItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}
func (n NamespaceItem) FilterValue() string {
	return n.Name
}

func initializeList() list.Model {
	l := list.New([]list.Item{}, NamespaceItemDelegate{}, pageWidth, pageHeight)
	l.Title = "Loading namespaces..."
	l.SetShowStatusBar(true)
	l.SetShowHelp(true)
	l.SetFilteringEnabled(true)
	l.SetShowFilter(true)
	l.Styles.Title = titleStyle
	l.SetStatusBarItemName("namespace", "namespaces")
	additionalKeys := func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("r"),
				key.WithHelp("r", "reload"),
			),
		}
	}
	l.Styles.NoItems.MarginLeft(2)
	//	l.Styles.StatusBar.MarginLeft(2)
	//	l.Styles.Title.MarginLeft(2)
	//	l.Styles.StatusBarActiveFilter.MarginLeft(2)
	//	l.Styles.StatusBarFilterCount.MarginLeft(2)
	l.KeyMap.GoToEnd.Unbind()
	l.KeyMap.GoToStart.Unbind()
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
