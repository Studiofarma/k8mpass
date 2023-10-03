package main

import (
	"fmt"
	"github.com/studiofarma/k8mpass/api"
	"log"
	"os"
	"plugin"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
)

func main() {
	_ = os.Remove("debug/k8s_debug.log")
	_, _ = tea.LogToFile("debug/k8s_debug.log", "DEBUG")
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Failed to load .env")
	} else {
		log.Println("Loaded .env correctly")
	}
	extensions, operations := loadPlugins()

	p := tea.NewProgram(
		initialModel(extensions, operations),
		tea.WithAltScreen(),
	)
	_, err := p.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func loadPlugins() ([]api.IExtension, []api.INamespaceOperation) {
	// load module
	// 1. open the so file to load the symbols
	plug, err := plugin.Open("plugins/extensions.so")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	symExtensions, err := plug.Lookup("GetNamespaceExtensions")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var extensions func() []api.IExtension
	extensions, ok := symExtensions.(func() []api.IExtension)
	if !ok {
		fmt.Println("unexpected type from module symbol ", symExtensions)
		os.Exit(1)
	}

	symOperations, err := plug.Lookup("GetNamespaceOperations")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var operations func() []api.INamespaceOperation
	operations, ok = symOperations.(func() []api.INamespaceOperation)
	if !ok {
		fmt.Println("unexpected type from module symbol ", symOperations)
		os.Exit(1)
	}
	return extensions(), operations()
}
