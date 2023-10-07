package main

import (
	"flag"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	"github.com/studiofarma/k8mpass/api"
	"log"
	"os"
	"plugin"
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
	log.Println("Loaded config correctly")

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

//func loadPlugins() api.IPlugins {
//	return Plugins
//}

func loadPlugins() api.IPlugins {
	pluginPath := flag.String("plugin", "", "path to plugin file")
	flag.Parse()
	if *pluginPath == "" {
		log.Println("No plugin to load")
		return api.Plugins{}
	}
	plug, err := plugin.Open(*pluginPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	symPlugins, err := plug.Lookup("Plugins")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var plugins api.IPlugins
	plugins, ok := symPlugins.(api.IPlugins)
	if !ok {
		fmt.Println("unexpected type from module symbol ", symPlugins)
		os.Exit(1)
	}
	return plugins
}
