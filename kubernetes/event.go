package kubernetes

import v1 "k8s.io/api/core/v1"

type NamespaceEvent struct {
	Type      EventType
	Namespace *v1.Namespace
}

type PodEvent struct {
	Type EventType
	Pod  *v1.Pod
}

type EventType string

const (
	Added     EventType = "ADDED"
	Modified  EventType = "MODIFIED"
	Deleted   EventType = "DELETED"
	Unhandled EventType = "UNHANDLED"
	Error     EventType = "ERROR"
	Closed    EventType = "CLOSED"
)
