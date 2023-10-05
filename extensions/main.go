package main

import "github.com/studiofarma/k8mpass/api"

var Plugins = api.Plugins{
	NamespaceOperations: namespaceOperations,
	NamespaceExtensions: namespaceExtensions,
	PodExtensions:       podExtensions,
}
