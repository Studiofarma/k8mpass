package kubernetes

import (
	"bufio"
	"context"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
	"strings"
	"time"
)

type ILogService interface {
	GetLogReader(nameSpace string, podName string, maxWidth int) error
	GetNextLogs() ([]string, bool)
}

func (c *Cluster) GetLogReader(nameSpace string, podName string, maxWidth int) error {
	podLogOpts := v1.PodLogOptions{
		Follow: true,
	}
	req := c.cs.CoreV1().Pods(nameSpace).GetLogs(podName, &podLogOpts)

	c.logLines = make(chan string)
	go LogToChannel(c.logLines, req, maxWidth)
	return nil
}

func truncateLines(s string, maxWidth int) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		if len(line) <= maxWidth {
			continue
		}
		lines[i] = line[:maxWidth-1]
	}
	return strings.Join(lines, "\n")
}

func (c *Cluster) GetNextLogs() ([]string, bool) {
	timer := time.NewTimer(1 * time.Second)
	var lines []string
	closed := false
out:
	for i := 0; i < 1000; i++ {
		select {
		case line, open := <-c.logLines:
			lines = append(lines, line)
			closed = closed || !open
			i++
		case <-timer.C:
			break out
		}
	}
	return lines, closed
}

func LogToChannel(ch chan string, req *rest.Request, maxWidth int) {
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		return
	}
	scanner := bufio.NewScanner(podLogs)
	for scanner.Scan() {
		truncatedLine := truncateLines(scanner.Text(), maxWidth)
		ch <- truncatedLine
	}
	close(ch)
}
