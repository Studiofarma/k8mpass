package main

type errMsg error

type clusterConnectedMsg struct{}

type startupMsg struct{}

type namespaceSelectedMsg struct {
	namespace string
}

type backToNamespaceSelectionMsg struct{}
