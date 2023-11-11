package log

import tea "github.com/charmbracelet/bubbletea"

type Message interface {
	isLogsMessage()
}
type UpdateLogsMsg struct{}

type RoutedMsg struct {
	Embedded tea.Msg
}

func (m UpdateLogsMsg) isLogsMessage() {}
func (m RoutedMsg) isLogsMessage()     {}
