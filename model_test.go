package main

import (
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
