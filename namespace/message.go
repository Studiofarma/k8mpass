package namespace

type Message interface {
	isNamespaceMessage()
}

type AddedNamespaceMsg struct {
	Namespace Item
}

type RemovedNamespaceMsg struct {
	Namespace Item
}

type NamespaceListMsg struct {
	Namespaces      []Item
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
