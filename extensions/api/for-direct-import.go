package api

import "github.com/studiofarma/k8mpass/api"

var SharedPlugins = api.Plugins{
	NamespaceOperations: NamespaceOperations,
	NamespaceExtensions: NamespaceExtensions,
	PodExtensions:       PodExtensions,
}
