package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type OperationModel struct {
	namespace   string
	operations  list.Model
	isCompleted bool
	output      string
	helpFooter  help.Model
}

func (o OperationModel) Init() tea.Cmd {
	return o.operations.StartSpinner()
}

func (o OperationModel) Update(msg tea.Msg) (OperationModel, tea.Cmd) {
	var cmds []tea.Cmd

	if o.isCompleted {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			/*			case "n":
						f := func() tea.Msg {
							return backToNamespaceSelectionMsg{}
						}
						cmds = append(cmds, f)*/
			case "backspace":
				f := func() tea.Msg {
					return backToOperationSelectionMsg{}
				}
				cmds = append(cmds, f)
			case "q":
				return o, tea.Quit

			}
		}
	} else {
		switch msg := msg.(type) {
		case namespaceSelectedMsg:
			o.namespace = msg.namespace
			styledNamespace := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("170")).Render(msg.namespace)
			o.operations.NewStatusMessage(styledNamespace)
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "backspace":
				f := func() tea.Msg {
					return backToNamespaceSelectionMsg{}
				}
				cmds = append(cmds, f)
			case "enter":
				i, ok := o.operations.SelectedItem().(NamespaceOperation)
				if ok {
					opCommand := i.Command(K8sCluster.kubernetes, o.namespace)
					f := func() tea.Msg {
						return operationSelectedMsg{operation: i.Name}
					}
					cmds = append(cmds, opCommand, f)
				} else {
					panic("Casting went wrong")
				}

			}
			om, omCmd := o.operations.Update(msg)
			o.operations = om
			cmds = append(cmds, omCmd)
		case operationResultMsg:
			o.isCompleted = true
			o.output = msg.body
		}
	}
	return o, tea.Batch(cmds...)
}

func (o OperationModel) View() string {
	/*gap := "  "
	header := gap + titleStyle.Render("Output of operation")
	styledOperation := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("170")).Render(o.operations.SelectedItem().FilterValue())
	*/s := ""
	if !o.isCompleted {
		s += o.operations.View()
	} else {
		s += o.output + "\n\n"
		s += "  " + o.helpFooter.View(keys)
	}
	return s
}

func (o *OperationModel) Reset() {
	o.isCompleted = false
	o.operations.ResetSelected()
}
