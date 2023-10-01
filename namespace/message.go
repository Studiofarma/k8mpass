package namespace

type NamespaceMessage interface {
	isNamespaceMessage()
}

type AddedNamespaceMsg struct {
	Namespace NamespaceItem
}

type RemovedNamespaceMsg struct {
	Namespace NamespaceItem
}

type NamespaceListMsg struct {
	Namespaces      []NamespaceItem
	ResourceVersion string
}

type NextEventMsg struct{}

type WatchingNamespacesMsg struct{}

type ErrorMsg struct {
	Err error
}

func (m AddedNamespaceMsg) isNamespaceMessage()     {}
func (m RemovedNamespaceMsg) isNamespaceMessage()   {}
func (m WatchingNamespacesMsg) isNamespaceMessage() {}
func (m NextEventMsg) isNamespaceMessage()          {}
func (m ErrorMsg) isNamespaceMessage()              {}
