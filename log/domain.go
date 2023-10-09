package log

import "github.com/muesli/reflow/truncate"

type Logs struct {
	Lines   []string
	channel chan string
}

func NewLogs() *Logs {
	return &Logs{
		Lines:   make([]string, 0),
		channel: make(chan string),
	}
}

func (l Logs) TruncatedLines(length int) []string {
	truncatedLines := make([]string, len(l.Lines))
	for idx, line := range l.Lines {
		truncatedLines[idx] = ellipsis(line, length)
	}
	return truncatedLines
}

func ellipsis(s string, length int) string {
	if len(s) > length {
		return truncate.StringWithTail(s, uint(length), "..")
	}
	return s
}
