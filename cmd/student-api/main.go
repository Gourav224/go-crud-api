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
	"github.com/gourav224/student-api/internal/http/handlers/student"
	"github.com/gourav224/student-api/internal/storage/sqlite"
)

func main() {
	// -------------------------------
	// 1️⃣ Load configuration
	// -------------------------------
	cfg := config.MustLoad()

	// -------------------------------
	// 2️⃣ Setup structured logger
	// -------------------------------
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("initializing server", "address", cfg.HTTPServer.Addr)

	// -------------------------------
	// 3️⃣ Initialize Database
	// -------------------------------
	db, err := sqlite.New(cfg)
	if err != nil {
		slog.Error("failed to initialize database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer func() {
		if cerr := db.Db.Close(); cerr != nil {
			slog.Warn("failed to close database", slog.String("error", cerr.Error()))
		}
	}()
	slog.Info("connected to sqlite database", "path", cfg.StoragePath)

	// -------------------------------
	// 4️⃣ Setup HTTP Router
	// -------------------------------
	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", student.New(db))
	router.HandleFunc("GET /api/students/{id}", student.GetById(db))
	router.HandleFunc("GET /api/students/", student.GetList(db))

	// -------------------------------
	// 5️⃣ Create HTTP Server
	// -------------------------------
	server := &http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: router,
	}

	// -------------------------------
	// 6️⃣ Graceful Shutdown Setup
	// -------------------------------
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// -------------------------------
	// 7️⃣ Run Server in Goroutine
	// -------------------------------
	go func() {
		slog.Info("server is listening", "address", cfg.HTTPServer.Addr)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", slog.String("error", err.Error()))
		}
	}()

	// -------------------------------
	// 8️⃣ Wait for Interrupt Signal
	// -------------------------------
	<-done
	slog.Warn("shutdown signal received")

	// -------------------------------
	// 9️⃣ Graceful Shutdown with Timeout
	// -------------------------------
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown gracefully", slog.String("error", err.Error()))
	} else {
		slog.Info("server stopped gracefully")
	}
}
