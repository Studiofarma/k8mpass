package kubernetes

import (
	"bytes"
	"context"
	"io"
	v1 "k8s.io/api/core/v1"
	"strings"
	"time"
)

type ILogService interface {
	GetLogReader(nameSpace string, podName string, maxWidth int) (string, error)
	GetNextLog() string
}

func (c *Cluster) GetLogReader(nameSpace string, podName string, maxWidth int) (string, error) {
	var tailLines int64 = 100
	podLogOpts := v1.PodLogOptions{
		TailLines: &tailLines,
	}
	req := c.cs.CoreV1().Pods(nameSpace).GetLogs(podName, &podLogOpts)
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		return "", err
	}
	c.logReader = &podLogs

	buf := bytes.Buffer{}
	_, err = io.Copy(&buf, podLogs)
	return truncateLines(buf.String(), maxWidth), nil
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

func (c *Cluster) GetNextLog() string {
	t := time.NewTimer(10 * time.Second)

	buf := make([]byte, 256)
	_, err := (*c.logReader).Read(buf)
	if err == io.EOF {
		return ""
	}
	<-t.C
	return ""
}
