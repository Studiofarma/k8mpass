package log

import (
	"github.com/muesli/reflow/truncate"
	"k8s.io/utils/strings/slices"
	"strings"
)

type Logs struct {
	lines        []string
	channel      chan string
	lineRenderer func(line string) string
}

func NewLogs() *Logs {
	return &Logs{
		lines:        make([]string, 0),
		channel:      make(chan string),
		lineRenderer: JsonRendered,
	}
}

func (l *Logs) AppendLines(line string) {
	if l.lineRenderer == nil {
		l.lines = append(l.lines, line)
	} else {
		renderedLine := l.lineRenderer(line)
		l.lines = append(l.lines, renderedLine)
	}
}

func (l Logs) TruncatedLines(length int) []string {
	return truncateList(l.lines, length)
}

func (l Logs) FilterAndTruncate(filterBy string, length int) []string {
	if filterBy == "" {
		return truncateList(l.lines, length)
	}
	var filteredList []string
	filteredList = slices.Filter(filteredList, l.lines, func(s string) bool {
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
