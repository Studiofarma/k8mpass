package namespace

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	v1 "k8s.io/api/core/v1"
)

type NamespaceItem struct {
	K8sNamespace v1.Namespace
	IsAwake      bool
}

func (n NamespaceItem) IsReviewApp() bool {
	return strings.HasPrefix(n.K8sNamespace.Name, "review")
}

type NamespaceItemDelegate struct{}

func (n NamespaceItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(NamespaceItem)
	if !ok {
		return
	}

	str := i.K8sNamespace.Name

	fn := lipgloss.NewStyle().PaddingLeft(4).Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170")).Render("> " + strings.Join(s, " "))
		}
	}
	fmt.Fprint(w, fn(str))
	//fmt.Fprint(w, "\n")
	if !i.IsReviewApp() {
		// do nothing
	} else if i.IsAwake {
		fmt.Fprint(w, lipgloss.NewStyle().PaddingLeft(4).Foreground(lipgloss.Color("#7d7d7d")).Render("Wide awake!"))
	} else {
		fmt.Fprint(w, lipgloss.NewStyle().PaddingLeft(4).Foreground(lipgloss.Color("#7d7d7d")).Render("Sleeping..."))
	}
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
	return n.K8sNamespace.Name
}

func New() list.Model {
	l := list.New([]list.Item{}, NamespaceItemDelegate{}, 80, 20)
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
