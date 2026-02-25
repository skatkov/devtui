package cmd

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/logging"
	"github.com/skatkov/devtui/tui/root"
	"github.com/spf13/cobra"
)

var (
	serveSSHHost  string
	serveSSHPort  string
	serveHTTPHost string
	serveHTTPPort string
	serveHostKey  string
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start DevTUI as an SSH server",
	Long: `Start DevTUI as an SSH server that users can connect to via SSH.

This also starts a simple HTTP server to show a landing page for web browsers.

Examples:
  # Start with default settings (SSH on :2222, HTTP on :8080)
  devtui serve

  # Custom ports
  devtui serve --ssh-port 22 --http-port 80

  # Specify host key path
  devtui serve --host-key /path/to/host_key`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runServer()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringVar(&serveSSHHost, "ssh-host", "0.0.0.0", "SSH server host")
	serveCmd.Flags().StringVar(&serveSSHPort, "ssh-port", "2222", "SSH server port")
	serveCmd.Flags().StringVar(&serveHTTPHost, "http-host", "0.0.0.0", "HTTP server host")
	serveCmd.Flags().StringVar(&serveHTTPPort, "http-port", "8080", "HTTP server port")
	serveCmd.Flags().StringVar(&serveHostKey, "host-key", ".ssh/devtui_host_key", "Path to SSH host key")
}

func runServer() error {
	// Create SSH server
	sshAddr := net.JoinHostPort(serveSSHHost, serveSSHPort)
	s, err := wish.NewServer(
		wish.WithAddress(sshAddr),
		wish.WithHostKeyPath(serveHostKey),
		wish.WithMiddleware(
			bubbleteaMiddlewareV2(teaHandler),
			activeterm.Middleware(), // Bubble Tea apps require a PTY
			logging.Middleware(),
		),
	)
	if err != nil {
		return fmt.Errorf("could not create SSH server: %w", err)
	}

	// Create HTTP server for landing page
	httpAddr := net.JoinHostPort(serveHTTPHost, serveHTTPPort)
	httpServer := &http.Server{
		Addr:              httpAddr,
		Handler:           http.HandlerFunc(landingPageHandler),
		ReadHeaderTimeout: 10 * time.Second,
	}

	// Handle shutdown signals
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start SSH server
	log.Info("Starting SSH server", "address", sshAddr)
	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("SSH server error", "error", err)
			done <- nil
		}
	}()

	// Start HTTP server
	log.Info("Starting HTTP server", "address", httpAddr)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("HTTP server error", "error", err)
			done <- nil
		}
	}()

	log.Info("DevTUI is ready!",
		"ssh", "ssh -p "+serveSSHPort+" "+serveSSHHost,
		"web", "http://"+httpAddr)

	<-done
	log.Info("Shutting down servers...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown both servers
	var shutdownErr error
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("SSH shutdown error", "error", err)
		shutdownErr = err
	}
	if err := httpServer.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("HTTP shutdown error", "error", err)
		shutdownErr = err
	}

	return shutdownErr
}

// teaHandler creates a new Bubble Tea model for each SSH session.
func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	// Get PTY info for initial window size
	pty, _, _ := s.Pty()

	model := root.RootScreenWithSize(pty.Window.Width, pty.Window.Height)

	return model, nil
}

type teaHandlerFunc func(sess ssh.Session) (tea.Model, []tea.ProgramOption)

func bubbleteaMiddlewareV2(handler teaHandlerFunc) wish.Middleware {
	return func(next ssh.Handler) ssh.Handler {
		return func(sess ssh.Session) {
			model, opts := handler(sess)
			if model == nil {
				next(sess)
				return
			}

			pty, windowChanges, ok := sess.Pty()
			if !ok {
				wish.Fatalln(sess, "no active terminal, skipping")
				return
			}

			opts = append(opts, makeTeaOptions(sess)...)
			env := append(sess.Environ(), "TERM="+pty.Term)
			opts = append(opts, tea.WithEnvironment(env))

			program := tea.NewProgram(model, opts...)

			ctx, cancel := context.WithCancel(sess.Context())
			go func() {
				for {
					select {
					case <-ctx.Done():
						program.Quit()
						return
					case w := <-windowChanges:
						program.Send(tea.WindowSizeMsg{Width: w.Width, Height: w.Height})
					}
				}
			}()

			if _, err := program.Run(); err != nil {
				log.Error("app exit with error", "error", err)
			}

			program.Kill()
			cancel()
			next(sess)
		}
	}
}

func makeTeaOptions(sess ssh.Session) []tea.ProgramOption {
	pty, _, ok := sess.Pty()
	if !ok || sess.EmulatedPty() {
		return []tea.ProgramOption{
			tea.WithInput(sess),
			tea.WithOutput(sess),
		}
	}

	return []tea.ProgramOption{
		tea.WithInput(pty.Slave),
		tea.WithOutput(pty.Slave),
	}
}

// landingPageHandler serves a simple landing page for web browsers
func landingPageHandler(w http.ResponseWriter, r *http.Request) {
	// Detect if this is a browser (has Accept: text/html)
	acceptHeader := r.Header.Get("Accept")
	if acceptHeader == "" || r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	domain := r.Host
	if domain == "" {
		domain = "devtui.sh"
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>DevTUI - Developer Tools in Your Terminal</title>
    <style>
        :root {
            --bg: #1a1b26;
            --fg: #a9b1d6;
            --accent: #7aa2f7;
            --accent2: #bb9af7;
            --code-bg: #24283b;
        }
        * {
            box-sizing: border-box;
            margin: 0;
            padding: 0;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
            background: var(--bg);
            color: var(--fg);
            min-height: 100vh;
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            padding: 2rem;
            text-align: center;
        }
        h1 {
            font-size: 3.5rem;
            margin-bottom: 1rem;
            background: linear-gradient(135deg, var(--accent), var(--accent2));
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
        }
        .subtitle {
            font-size: 1.25rem;
            margin-bottom: 3rem;
            opacity: 0.8;
        }
        .connect-box {
            background: var(--code-bg);
            border-radius: 12px;
            padding: 2rem 3rem;
            margin-bottom: 2rem;
        }
        .connect-label {
            font-size: 0.9rem;
            text-transform: uppercase;
            letter-spacing: 0.1em;
            opacity: 0.6;
            margin-bottom: 1rem;
        }
        .connect-cmd {
            font-family: 'SF Mono', 'Fira Code', 'Consolas', monospace;
            font-size: 1.5rem;
            color: var(--accent);
            cursor: pointer;
            transition: opacity 0.2s;
        }
        .connect-cmd:hover {
            opacity: 0.8;
        }
        .features {
            display: flex;
            flex-wrap: wrap;
            gap: 1rem;
            justify-content: center;
            max-width: 600px;
            margin-bottom: 2rem;
        }
        .feature {
            background: var(--code-bg);
            padding: 0.5rem 1rem;
            border-radius: 6px;
            font-size: 0.9rem;
        }
        .footer {
            margin-top: 3rem;
            opacity: 0.5;
            font-size: 0.9rem;
        }
        .footer a {
            color: var(--accent);
            text-decoration: none;
        }
        .footer a:hover {
            text-decoration: underline;
        }
        .copy-hint {
            font-size: 0.8rem;
            opacity: 0.5;
            margin-top: 0.5rem;
        }
    </style>
</head>
<body>
    <h1>DevTUI</h1>
    <p class="subtitle">A Swiss-army knife terminal app for developers</p>
    
    <div class="connect-box">
        <div class="connect-label">Connect via SSH</div>
        <code class="connect-cmd" onclick="navigator.clipboard.writeText('ssh %s')">ssh %s</code>
        <div class="copy-hint">Click to copy</div>
    </div>
    
    <div class="features">
        <span class="feature">JSON Formatter</span>
        <span class="feature">Base64 Encode/Decode</span>
        <span class="feature">UUID Generator</span>
        <span class="feature">YAML Formatter</span>
        <span class="feature">TOML Converter</span>
        <span class="feature">XML Formatter</span>
        <span class="feature">Cron Parser</span>
        <span class="feature">URL Extractor</span>
        <span class="feature">+ Many More</span>
    </div>
    
    <div class="footer">
        <a href="https://github.com/skatkov/devtui" target="_blank">GitHub</a> · 
        <a href="https://devtui.com" target="_blank">Documentation</a>
    </div>

    <script>
        document.querySelector('.connect-cmd').addEventListener('click', function() {
            navigator.clipboard.writeText('ssh %s');
            this.textContent = 'Copied!';
            setTimeout(() => {
                this.textContent = 'ssh %s';
            }, 2000);
        });
    </script>
</body>
</html>`, domain, domain, domain, domain)

	_, _ = w.Write([]byte(html))
}
