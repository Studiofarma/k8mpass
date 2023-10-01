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

type NextEventMsg struct{}

type WatchingNamespacesMsg struct{}

type ErrorMsg struct {
	err error
}

func (m AddedNamespaceMsg) isNamespaceMessage()     {}
func (m RemovedNamespaceMsg) isNamespaceMessage()   {}
func (m WatchingNamespacesMsg) isNamespaceMessage() {}
func (m NextEventMsg) isNamespaceMessage()          {}
func (m ErrorMsg) isNamespaceMessage()              {}
