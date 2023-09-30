package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {

	}
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	_, err = p.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
