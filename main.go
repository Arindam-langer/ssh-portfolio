package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/log/v2"
	"charm.land/ssh"
	"charm.land/wish/v2"
	"charm.land/wish/v2/activeterm"
	"charm.land/wish/v2/bubbletea"
	"charm.land/wish/v2/logging"
	"github.com/Arindam-Langer/ssh-portfolio/internal/data"
	"github.com/Arindam-Langer/ssh-portfolio/internal/tui"
)

const (
	host = "0.0.0.0"
	port = "2222"
)

func main() {
	if err := data.LoadConfig("config.yaml"); err != nil {
		log.Error("Could not load config.yaml", "error", err)
		os.Exit(1)
	}

	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
		os.Exit(1)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println()
	fmt.Println("  ┌─────────────────────────────────────────────────┐")
	fmt.Println("  │                                                 │")

	text := fmt.Sprintf("🚀 %s's SSH Portfolio is running!", data.AppConfig.Profile.Name)
	// Calculate padding to center or simply fit in the box. The box is 47 chars wide inside.
	padLen := 47 - len([]rune(text))
	if padLen < 0 {
		padLen = 0
	}
	padding := ""
	for i := 0; i < padLen; i++ {
		padding += " "
	}

	fmt.Printf("  │   %s%s│\n", text, padding)

	fmt.Println("  │                                                 │")
	fmt.Printf("  │   Connect: ssh -p %s localhost              │\n", port)
	fmt.Println("  │                                                 │")
	fmt.Println("  │   Press Ctrl+C to stop the server               │")
	fmt.Println("  │                                                 │")
	fmt.Println("  └─────────────────────────────────────────────────┘")
	fmt.Println()

	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := s.Pty()
	m := tui.NewModel(pty.Term, pty.Window.Width, pty.Window.Height)
	return m, []tea.ProgramOption{}
}
