package namespace

type NamespaceMessage interface {
	isNamespaceMessage()
}

type AddedNamespaceMsg struct {
	Namespace NamespaceItem
}

func (m AddedNamespaceMsg) isNamespaceMessage() {}

type RemovedNamespaceMsg struct {
	Namespace NamespaceItem
}

func (m RemovedNamespaceMsg) isNamespaceMessage() {}
