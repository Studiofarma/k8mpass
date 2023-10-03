package namespace

type Message interface {
	isNamespaceMessage()
}

type AddedMsg struct {
	Namespace Item
}

type RemovedMsg struct {
	Namespace Item
}

type ListMsg struct {
	Namespaces      []Item
	ResourceVersion string
}

type NextEventMsg struct{}

type WatchingMsg struct{}

type ErrorMsg struct {
	Err error
}

func (m AddedMsg) isNamespaceMessage()     {}
func (m RemovedMsg) isNamespaceMessage()   {}
func (m WatchingMsg) isNamespaceMessage()  {}
func (m NextEventMsg) isNamespaceMessage() {}
func (m ErrorMsg) isNamespaceMessage()     {}
