package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
	NamespaceSelection key.Binding
	OperationSelection key.Binding
	Quit               key.Binding
}

var keys = keyMap{
	NamespaceSelection: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "back to namespaces"),
	),
	OperationSelection: key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "back to operations"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.NamespaceSelection, k.OperationSelection, k.Quit}
}
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.NamespaceSelection, k.OperationSelection, k.Quit},
	}
}

func initializeHelpFooter() help.Model {
	h := help.New()
	return h
}
