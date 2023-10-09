package log

type Message interface {
	isLogsMessage()
}
type UpdateLogsMsg struct{}

func (m UpdateLogsMsg) isLogsMessage() {}
