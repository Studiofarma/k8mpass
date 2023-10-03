package api

import (
	tea "github.com/charmbracelet/bubbletea"
	"k8s.io/client-go/kubernetes"
)

type INamespaceOperation interface {
	GetName() string
	GetCommand() K8mpassCommand
	GetCondition() K8mpassCondition

	FilterValue() string
}

type K8mpassCommand func(model *kubernetes.Clientset, namespace string) tea.Cmd

type K8mpassCondition func(cs *kubernetes.Clientset, namespace string) bool

type NamespaceOperation struct {
	Name      string
	Command   K8mpassCommand
	Condition K8mpassCondition
}

func (o NamespaceOperation) GetName() string {
	return o.Name
}

func (o NamespaceOperation) GetCommand() K8mpassCommand {
	return o.Command
}

func (o NamespaceOperation) GetCondition() K8mpassCondition {
	return o.Condition
}

func (o NamespaceOperation) FilterValue() string {
	return o.Name
}
