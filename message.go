package main

type errMsg error

type clusterConnectedMsg struct {
	context string
}

type startupMsg struct{}

type namespaceSelectedMsg struct {
	namespace string
}

type backToNamespaceSelectionMsg struct{}
