package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/ookapratama/go-todo-api/config"
	"github.com/ookapratama/go-todo-api/db"
	"github.com/ookapratama/go-todo-api/handler"
	"github.com/ookapratama/go-todo-api/middleware"
	"github.com/ookapratama/go-todo-api/repository"
	"github.com/ookapratama/go-todo-api/service"
)

func main() {
	// Setup zerolog - pretty print for development
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	// Load configuration
	cfg := config.LoadConfig()

	// Determine migrations directory (relative to executable)
	migrationsDir := getMigrationsDir()

	// Check if running a migration command
	if len(os.Args) > 1 {
		runMigrationCommand(cfg, migrationsDir, os.Args[1])
		return
	}

	// Normal server startup
	startServer(cfg, migrationsDir)
}

// getMigrationsDir finds the migrations directory.
func getMigrationsDir() string {
	// Try relative to working directory first
	candidates := []string{
		"db/migrations",
		"backend/db/migrations",
	}

	for _, dir := range candidates {
		absDir, _ := filepath.Abs(dir)
		if info, err := os.Stat(absDir); err == nil && info.IsDir() {
			return absDir
		}
	}

	// Default
	absDir, _ := filepath.Abs("db/migrations")
	return absDir
}

// runMigrationCommand handles CLI migration commands (like php artisan migrate).
func runMigrationCommand(cfg *config.Config, migrationsDir string, command string) {
	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer database.Close()

	migrator := db.NewMigrator(database, migrationsDir)

	switch command {
	case "migrate":
		// Similar to: php artisan migrate
		fmt.Println("📦 Running migrations...")
		if err := migrator.MigrateUp(); err != nil {
			log.Fatal().Err(err).Msg("Migration failed")
		}

	case "migrate:rollback":
		// Similar to: php artisan migrate:rollback
		fmt.Println("⏪ Rolling back last batch...")
		if err := migrator.MigrateDown(); err != nil {
			log.Fatal().Err(err).Msg("Rollback failed")
		}

	case "migrate:reset":
		// Similar to: php artisan migrate:reset
		fmt.Println("🔄 Resetting all migrations...")
		if err := migrator.MigrateReset(); err != nil {
			log.Fatal().Err(err).Msg("Reset failed")
		}

	case "migrate:fresh":
		// Similar to: php artisan migrate:fresh
		fmt.Println("🗑️  Dropping all tables and re-migrating...")
		if err := migrator.MigrateFresh(); err != nil {
			log.Fatal().Err(err).Msg("Fresh migration failed")
		}

	case "migrate:status":
		// Similar to: php artisan migrate:status
		statuses, err := migrator.Status()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to get migration status")
		}
		fmt.Println("")
		fmt.Println("┌─────────────────────────────────────────────┬──────────┐")
		fmt.Println("│ Migration                                   │ Status   │")
		fmt.Println("├─────────────────────────────────────────────┼──────────┤")
		for _, s := range statuses {
			icon := "⏳"
			if s["status"] == "Ran" {
				icon = "✅"
			}
			fmt.Printf("│ %-43s │ %s %-5s │\n", s["migration"], icon, s["status"])
		}
		fmt.Println("└─────────────────────────────────────────────┴──────────┘")
		fmt.Println("")

	case "make:migration":
		// Similar to: php artisan make:migration [name]
		if len(os.Args) < 3 {
			log.Fatal().Msg("Migration name is required. Usage: go run main.go make:migration [name]")
		}
		name := os.Args[2]
		timestamp := time.Now().Format("20060102150405")
		baseName := fmt.Sprintf("%s_%s", timestamp, name)

		upFile := filepath.Join(migrationsDir, baseName+".up.sql")
		downFile := filepath.Join(migrationsDir, baseName+".down.sql")

		err := os.WriteFile(upFile, []byte("-- Up Migration\n"), 0644)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create up migration file")
		}
		err = os.WriteFile(downFile, []byte("-- Down Migration\n"), 0644)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create down migration file")
		}

		fmt.Printf("✨ Created migration files:\n  - %s\n  - %s\n", upFile, downFile)

	default:
		fmt.Println("❌ Unknown command:", command)
		fmt.Println("")
		fmt.Println("Available commands:")
		fmt.Println("  migrate            Run all pending migrations")
		fmt.Println("  migrate:rollback   Rollback the last batch of migrations")
		fmt.Println("  migrate:reset      Rollback all migrations")
		fmt.Println("  migrate:fresh      Drop all tables and re-run all migrations")
		fmt.Println("  migrate:status     Show the status of each migration")
		fmt.Println("  make:migration     Create new up/down migration files")
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("  go run main.go make:migration create_users_table")
		os.Exit(1)
	}
}

// startServer starts the HTTP server with auto-migration.
func startServer(cfg *config.Config, migrationsDir string) {
	log.Info().Msg("Starting Go To-Do API Server...")

	// Connect to database
	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer database.Close()

	// Auto-run migrations on startup
	migrator := db.NewMigrator(database, migrationsDir)
	if err := migrator.MigrateUp(); err != nil {
		log.Fatal().Err(err).Msg("Auto-migration failed")
	}

	// Initialize layers
	taskRepo := repository.NewTaskRepository(database)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := handler.NewTaskHandler(taskService)

	// Setup router
	r := chi.NewRouter()

	// Apply middleware
	r.Use(middleware.Recovery)
	r.Use(middleware.CORS)
	r.Use(middleware.Logger)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status":"healthy","timestamp":"`+time.Now().Format(time.RFC3339)+`"}`)
	})

	// Register task routes
	taskHandler.RegisterRoutes(r)

	// Create HTTP server
	addr := ":" + cfg.ServerPort
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		log.Info().Str("address", addr).Msg("Server is running")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	// Graceful shutdown with 30-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exited gracefully")
}
