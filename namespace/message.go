package namespace

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Message interface {
	isNamespaceMessage()
}

type AddedMsg struct {
	Namespace Item
}

type ModifiedMsg struct {
	Namespace Item
}

type RemovedMsg struct {
	Namespace Item
}

type ListMsg struct {
	Namespaces      []Item
	ResourceVersion string
}

type NextEventMsg struct{}

type WatchingMsg struct{}

type ReloadTick struct{}

type ReloadExtensionsMsg struct {
	Properties map[string][]Property
}

type RoutedMsg struct {
	Embedded tea.Msg
}

type ErrorMsg struct {
	Err error
}

func (m AddedMsg) isNamespaceMessage()            {}
func (m RemovedMsg) isNamespaceMessage()          {}
func (m WatchingMsg) isNamespaceMessage()         {}
func (m NextEventMsg) isNamespaceMessage()        {}
func (m ReloadTick) isNamespaceMessage()          {}
func (m ReloadExtensionsMsg) isNamespaceMessage() {}
func (m RoutedMsg) isNamespaceMessage()           {}
func (m ErrorMsg) isNamespaceMessage()            {}
