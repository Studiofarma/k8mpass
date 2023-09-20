package main

import (
	"github.com/charmbracelet/bubbles/list"
	"k8s.io/client-go/kubernetes"
	"testing"
)

func TestModelConnection(t *testing.T) {
	model := initialModel()

	if model.isConnected {
		t.Fatalf("Before a connection the model should not be connected")
	}
	updatedModel, cmd := model.Update(clusterConnectedMsg{clientset: nil})
	if !updatedModel.(K8mpassModel).isConnected {
		t.Fatalf("After a connection the model should be connected but its not")
	}

	msg, ok := cmd().(fetchNamespacesMsg)
	if !ok {
		t.Fatalf("After a connection it should send fetchNamespacesMsg but it is sending %T", msg)
	}
}

type fakeCluster struct {
}

func (f fakeCluster) FetchNamespaces() ([]string, error) {
	return []string{"X", "Y"}, nil
}

func (f fakeCluster) SetClientset(clientset *kubernetes.Clientset) {
}

func (f fakeCluster) GetClientset() *kubernetes.Clientset {
	return nil
}

func TestFetchNamespace(t *testing.T) {
	var items = []list.Item{
		item("X"),
		item("Y"),
	}

	l := list.New(items, list.NewDefaultDelegate(), 2, 5)

	model := initialModel()

	model.cluster = fakeCluster{}

	updatedModel, _ := model.Update(fetchNamespacesMsg{})
	if updatedModel.(K8mpassModel).namespacesList.Items()[0] != l.Items()[0] {
		t.Fatalf("The list of namespace is not that expected")
	}
}
