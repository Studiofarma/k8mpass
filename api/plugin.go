package api

type IPlugins interface {
	GetNamespaceExtensions() []INamespaceExtension
	GetPodExtensions() []IPodExtension
	GetNamespaceOperations() []INamespaceOperation
}

type Plugins struct {
	NamespaceExtensions []INamespaceExtension
	PodExtensions       []IPodExtension
	NamespaceOperations []INamespaceOperation
}

func (p Plugins) GetNamespaceExtensions() []INamespaceExtension {
	return p.NamespaceExtensions
}

func (p Plugins) GetPodExtensions() []IPodExtension {
	return p.PodExtensions
}

func (p Plugins) GetNamespaceOperations() []INamespaceOperation {
	return p.NamespaceOperations
}
