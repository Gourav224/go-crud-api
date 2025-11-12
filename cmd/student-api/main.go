package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gourav224/student-api/internal/config"
)

func main() {
	// -------------------------------
	// 1️⃣ Load configuration
	// -------------------------------
	cfg := config.MustLoad() // load settings like port, DB, etc.

	// -------------------------------
	// 2️⃣ Initialize structured logger
	// -------------------------------
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo, // log level: Debug, Info, Warn, Error
	}))
	slog.SetDefault(logger) // make it global, accessible via slog.Info/Error etc.

	slog.Info("starting server", "address", cfg.HTTPServer.Addr)

	// -------------------------------
	// 3️⃣ Setup routes using ServeMux
	// -------------------------------
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome to api"))
	})

	// -------------------------------
	// 4️⃣ Setup server
	// -------------------------------
	server := &http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: router,
	}

	// -------------------------------
	// 5️⃣ Listen for OS signals
	// -------------------------------
	// Create a channel that receives OS signals (like Ctrl+C)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// -------------------------------
	// 6️⃣ Start the server in a goroutine
	// -------------------------------
	go func() {
		slog.Info("server is listening", "address", cfg.HTTPServer.Addr)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", slog.String("error", err.Error()))
		}
	}()

	// -------------------------------
	// 7️⃣ Block until signal is received
	// -------------------------------
	<-done
	slog.Warn("shutdown signal received")

	// -------------------------------
	// 8️⃣ Graceful shutdown with timeout
	// -------------------------------
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown gracefully", slog.String("error", err.Error()))
	} else {
		slog.Info("server stopped gracefully")
	}
}
