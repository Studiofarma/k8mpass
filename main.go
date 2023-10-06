package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	"github.com/studiofarma/k8mpass/api"
	"log"
	"os"
)

func main() {
	_ = os.Remove("debug/k8s_debug.log")
	_, _ = tea.LogToFile("debug/k8s_debug.log", "DEBUG")
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Failed to load .env")
	} else {
		log.Println("Loaded .env correctly")
	}
	plugins := loadPlugins()

	p := tea.NewProgram(
		initialModel(plugins),
		tea.WithAltScreen(),
	)
	_, err := p.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func loadPlugins() api.IPlugins {
	return Plugins
}

//func loadPlugins() api.IPlugins {
//	args := os.Args
//	if len(args) < 3 {
//		return api.Plugins{}
//	}
//	pluginPath := os.Args[2]
//	plug, err := plugin.Open(pluginPath)
//	if err != nil {
//		fmt.Println(err)
//		os.Exit(1)
//	}
//	symPlugins, err := plug.Lookup("Plugins")
//	if err != nil {
//		fmt.Println(err)
//		os.Exit(1)
//	}
//	var plugins api.IPlugins
//	plugins, ok := symPlugins.(api.IPlugins)
//	if !ok {
//		fmt.Println("unexpected type from module symbol ", symPlugins)
//		os.Exit(1)
//	}
//	return plugins
//}
