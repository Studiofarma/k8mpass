package kubernetes

import (
	"bufio"
	"context"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
	"log"
)

type ILogService interface {
	SendLogsToChannel(ctx context.Context, nameSpace string, podName string, ch chan string) error
}

func (c *Cluster) SendLogsToChannel(ctx context.Context, nameSpace string, podName string, ch chan string) error {
	tailLines := int64(1000)
	podLogOpts := v1.PodLogOptions{
		Follow:    true,
		TailLines: &tailLines,
	}
	req := c.cs.CoreV1().Pods(nameSpace).GetLogs(podName, &podLogOpts)
	go LogToChannel(ctx, ch, req)
	return nil
}
func LogToChannel(ctx context.Context, ch chan string, req *rest.Request) {
	podLogs, err := req.Stream(ctx)
	if err != nil {
		return
	}
	scanner := bufio.NewScanner(podLogs)
	for scanner.Scan() {
		ch <- scanner.Text()
	}
	close(ch)
	log.Println("Closing log channel")
}
