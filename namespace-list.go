package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"io"
	"strings"
)

type NamespaceItem struct {
	name string
}

type NamespaceItemDelegate struct{}

func (n NamespaceItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(NamespaceItem)
	if !ok {
		return
	}

	str := i.name

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
	return n.name
}

func initializeList() list.Model {
	l := list.New([]list.Item{}, NamespaceItemDelegate{}, pageWidth, pageHeight)
	l.Title = "Loading namespaces..."
	l.SetShowStatusBar(true)
	l.SetShowHelp(true)
	l.SetShowTitle(false)
	l.SetFilteringEnabled(true)
	l.SetShowFilter(true)
	l.Styles.Title = titleStyle
	additionalKeys := func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("r"),
				key.WithHelp("r", "refresh namespaces"),
			),
		}
	}
	l.AdditionalFullHelpKeys = additionalKeys
	l.AdditionalShortHelpKeys = additionalKeys
	return l
}
