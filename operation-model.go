package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"k8s.io/client-go/kubernetes"
)

type OperationModel struct {
	namespace   string
	operations  list.Model
	clientset   *kubernetes.Clientset
	isCompleted bool
	helpFooter  help.Model
}

func (o OperationModel) Init() tea.Cmd {
	return nil
}

func (o OperationModel) Update(msg tea.Msg) (OperationModel, tea.Cmd) {
	var cmds []tea.Cmd

	if o.isCompleted {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "n":
				f := func() tea.Msg {
					return backToNamespaceSelectionMsg{}
				}
				cmds = append(cmds, f)
			case "o":
				f := func() tea.Msg {
					return backToOperationSelectionMsg{}
				}
				cmds = append(cmds, f)
			case "q":
				return o, tea.Quit

			}
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			i, ok := o.operations.SelectedItem().(NamespaceOperation)
			if ok {
				opCommand := i.Command(o.clientset, o.namespace)
				cmds = append(cmds, opCommand)
			} else {
				panic("Casting went wrong")
			}

		}
		om, omCmd := o.operations.Update(msg)
		o.operations = om
		cmds = append(cmds, omCmd)
	case wakeUpReviewMsg:
		o.isCompleted = true
	}
	return o, tea.Batch(cmds...)
}

func (o OperationModel) View() string {
	gap := "  "
	header := gap + titleStyle.Render("Output of operation")
	styledOperation := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("170")).Render(o.operations.SelectedItem().FilterValue())
	s := ""
	if !o.isCompleted {
		s += o.operations.View()
	} else {
		s += header + " " + styledOperation + "\n\n"
		s += gap + "Okay I'm done\n"
		for i := 1; i <= pageHeight-4; i++ {
			s += "\n"
		}
		s += "  " + o.helpFooter.View(keys)
	}
	return s
}

func (o *OperationModel) Reset() {
	o.isCompleted = false
}
