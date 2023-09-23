package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type OperationItem struct {
	name string
}

type OperationItemDelegate struct {
	NamespaceOperation NamespaceOperation
}

func (o OperationItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(NamespaceOperation)
	if !ok {
		return
	}

	str := i.Name

	fn := lipgloss.NewStyle().PaddingLeft(4).Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#ffcb78")).Render("> " + strings.Join(s, " "))
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

func initializeOperationList(ops []NamespaceOperation) list.Model {
	var items []list.Item
	for _, op := range ops {
		items = append(items, op)
	}
	l := list.New(items, OperationItemDelegate{}, pageWidth, pageHeight)
	l.Title = "Select an operation on"
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(true)
	l.SetFilteringEnabled(false)
	l.SetShowFilter(false)
	l.Styles.Title = titleStyle
	additionalKeys := func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("backspace"),
				key.WithHelp("backspace", "back to namespaces"),
			),
		}
	}
	l.AdditionalFullHelpKeys = additionalKeys
	l.AdditionalShortHelpKeys = additionalKeys
	return l
}
