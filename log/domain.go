package log

import (
	"github.com/muesli/reflow/truncate"
	"k8s.io/utils/strings/slices"
	"strings"
)

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
	return truncateList(l.Lines, length)
}

func (l Logs) FilterAndTruncate(filterBy string, length int) []string {
	if filterBy == "" {
		return truncateList(l.Lines, length)
	}
	var filteredList []string
	filteredList = slices.Filter(filteredList, l.Lines, func(s string) bool {
		return strings.Contains(s, filterBy)
	})
	return truncateList(filteredList, length)
}

func ellipsis(s string, length int) string {
	if len(s) > length {
		return truncate.StringWithTail(s, uint(length), "..")
	}
	return s
}

func truncateList(list []string, length int) []string {
	truncatedLines := make([]string, len(list))
	for idx, line := range list {
		truncatedLines[idx] = ellipsis(line, length)
	}
	return truncatedLines
}
