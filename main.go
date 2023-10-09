package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	"github.com/studiofarma/k8mpass/api"
	"github.com/studiofarma/k8mpass/config"
	api2 "github.com/studiofarma/k8mpass/extensions/api"
	"log"
	"os"
	"plugin"
	"runtime"
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
	config.LoadFlags()
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
	switch runtime.GOOS {
	case "windows":
		return loadPluginsWindows()
	default:
		return loadPluginsLinux()
	}
}

func loadPluginsWindows() api.IPlugins {
	return api2.SharedPlugins
}

func loadPluginsLinux() api.IPlugins {
	if config.Plugin == "" {
		log.Println("No plugin to load")
		return api.Plugins{}
	}
	plug, err := plugin.Open(config.Plugin)
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
