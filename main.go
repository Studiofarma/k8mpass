package main

import (
	"context"
	"errors"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	lm "github.com/charmbracelet/wish/logging"
	"github.com/studiofarma/k8mpass/api"
	"github.com/studiofarma/k8mpass/config"
	"os"
	"os/signal"
	"plugin"
	"syscall"
	"time"
)

const (
	host = "localhost"
	port = 23234
)

func main() {
	config.LoadFlags()
	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", host, port)),
		wish.WithHostKeyPath(".ssh/term_info_ed25519"),
		wish.WithMiddleware(
			bm.Middleware(teaHandler),
			lm.Middleware(),
		),
	)
	if err != nil {
		log.Error("could not start server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("could not stop server", "error", err)
	}
}

// You can wire any Bubble Tea model up to the middleware with a function that
// handles the incoming ssh.Session. Here we just grab the terminal info and
// pass it to the new model. You can also return tea.ProgramOptions (such as
// tea.WithAltScreen) on a session by session basis.
func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	_, _, active := s.Pty()
	if !active {
		wish.Fatalln(s, "no active terminal, skipping")
		return nil, nil
	}
	m := model(s.User())
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

func model(user string) Model {
	log.Infof("Creating model for user " + user)
	return initialModel(p, user)
}

var p = loadPlugins()

func loadPlugins() api.IPlugins {
	if config.Plugin == "" {
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
