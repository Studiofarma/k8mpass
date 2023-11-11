package log

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	k8mpasskube "github.com/studiofarma/k8mpass/kubernetes"
	"time"
)

type MessageHandler struct {
	service k8mpasskube.ILogService
	newlogs chan interface{}
	cancel  context.CancelFunc
}

func NewHandler(service k8mpasskube.ILogService) *MessageHandler {
	return &MessageHandler{
		service: service,
		newlogs: make(chan interface{}, 1),
	}
}

func (handler *MessageHandler) FollowLogs(namespace string, pod string, logs *Logs) {
	ctx, cancel := context.WithCancel(context.Background())
	handler.cancel = cancel
	err := handler.service.SendLogsToChannel(ctx, namespace, pod, logs.channel)
	go CacheLogs(logs, handler.newlogs)
	if err != nil {
		return
	}
}

func (handler *MessageHandler) CloseLogs() {
	handler.cancel()
}

func CacheLogs(logs *Logs, newlogs chan interface{}) {
	for line := range logs.channel {
		//Add message to queue if not full
		select {
		case newlogs <- "something":
		default:
		}
		logs.AppendLines(line)
	}
}

func (handler *MessageHandler) TickUpdateLogs() tea.Msg {
	time.Sleep(time.Millisecond)
	<-handler.newlogs
	return UpdateLogsMsg{}
}

func Route(cmds ...tea.Cmd) []tea.Cmd {
	var ret []tea.Cmd
	for _, cmd := range cmds {
		if cmd == nil {
			continue
		}
		ret = append(ret, func() tea.Msg {
			return RoutedMsg{Embedded: cmd()}
		})
	}
	return ret
}
