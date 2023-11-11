package log

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	tealog "github.com/charmbracelet/log"
	"time"
)

var (
	timestampStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#949494"))
)

type LogJson struct {
	Timestamp time.Time `json:"@timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
}

func (l LogJson) Render() string {
	t := timestampStyle.Render(l.Timestamp.Format("01/02-15:04:05"))
	level := levelStyle(l.Level).MaxWidth(5).Render(l.Level)
	return fmt.Sprintf("%s %s: %s", t, level, l.Message)
}

func levelStyle(level string) lipgloss.Style {
	switch level {
	case "INFO", "DEBUG":
		return tealog.InfoLevelStyle
	case "WARN", "WARNING":
		return tealog.WarnLevelStyle
	case "ERROR":
		return tealog.ErrorLevelStyle
	default:
		return lipgloss.NewStyle()
	}
}

func JsonRendered(line string) string {
	var d LogJson
	err := json.Unmarshal([]byte(line), &d)
	if err != nil {
		return line
	}
	return d.Render()
}
