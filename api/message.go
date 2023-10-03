package api

type NoOutputResultMsg struct {
	Success bool
	Message string
}

type AvailableOperationsMsg struct {
	Operations []NamespaceOperation
}
