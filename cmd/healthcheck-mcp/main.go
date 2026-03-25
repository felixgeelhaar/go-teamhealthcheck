package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	bolt "github.com/felixgeelhaar/bolt"
	"github.com/felixgeelhaar/mcp-go"
	"github.com/felixgeelhaar/mcp-go/middleware"

	"github.com/felixgeelhaar/go-teamhealthcheck/internal/auth"
	mcptools "github.com/felixgeelhaar/go-teamhealthcheck/internal/mcp"
	"github.com/felixgeelhaar/go-teamhealthcheck/internal/storage"
)

func main() {
	home, _ := os.UserHomeDir()
	defaultDB := filepath.Join(home, ".healthcheck-mcp", "data.db")
	defaultAuth := filepath.Join(home, ".healthcheck-mcp", "auth.json")

	dbPath := flag.String("db", defaultDB, "Path to SQLite database file")
	mode := flag.String("mode", "stdio", "Transport mode: stdio or http")
	addr := flag.String("addr", ":8080", "HTTP listen address (only used with --mode http)")
	authConfig := flag.String("auth", defaultAuth, "Path to auth config file (only used with --mode http)")
	dev := flag.Bool("dev", false, "Development mode (colored console logging)")
	flag.Parse()

	// Initialize logger
	var handler bolt.Handler
	if *dev {
		handler = bolt.NewConsoleHandler(os.Stderr)
	} else {
		handler = bolt.NewJSONHandler(os.Stderr)
	}
	logger := bolt.New(handler)

	store, err := storage.New(*dbPath, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize storage")
	}
	defer store.Close()

	srv := mcptools.NewServer(store, logger)

	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	switch *mode {
	case "http":
		// Build middleware chain for HTTP mode
		var middlewares []middleware.Middleware
		middlewares = append(middlewares,
			middleware.Recover(),
			middleware.Timeout(30*time.Second),
		)

		// Load auth config if it exists
		if tokenValidator, err := auth.LoadConfig(*authConfig); err == nil {
			logger.Info().Str("config", *authConfig).Msg("auth enabled")
			middlewares = append(middlewares,
				middleware.Auth(
					middleware.BearerTokenAuthenticator(tokenValidator),
					middleware.WithAuthSkipMethods("initialize", "ping"),
				),
			)
		} else if !os.IsNotExist(err) {
			logger.Warn().Err(err).Msg("failed to load auth config, running without auth")
		} else {
			logger.Info().Msg("no auth config found, running without auth")
		}

		logger.Info().Str("addr", *addr).Msg("starting HTTP/SSE transport")
		if err := mcp.ServeHTTPWithMiddleware(ctx, srv, *addr,
			[]mcp.HTTPOption{
				mcp.WithReadTimeout(30 * time.Second),
				mcp.WithWriteTimeout(30 * time.Second),
			},
			mcp.WithMiddleware(middlewares...),
		); err != nil {
			logger.Fatal().Err(err).Msg("HTTP server error")
		}
	default:
		logger.Debug().Msg("starting stdio transport")
		if err := mcp.ServeStdio(ctx, srv); err != nil {
			logger.Fatal().Err(err).Msg("stdio server error")
		}
	}
}
