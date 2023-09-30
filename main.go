package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	godotenv.Load(".env")
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	_, err := p.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
