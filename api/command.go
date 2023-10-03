package api

import (
	tea "github.com/charmbracelet/bubbletea"
	"k8s.io/client-go/kubernetes"
)

type K8mpassCommand func(model *kubernetes.Clientset, namespace string) tea.Cmd

type K8mpassCondition func(cs *kubernetes.Clientset, namespace string) bool

type NamespaceOperation struct {
	Name      string
	Command   K8mpassCommand
	Condition K8mpassCondition
}

func (o NamespaceOperation) FilterValue() string {
	return o.Name
}
