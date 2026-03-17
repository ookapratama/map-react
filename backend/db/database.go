package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"

	"github.com/ookapratama/go-todo-api/config"
)

// Connect establishes a connection to PostgreSQL (local or Supabase).
func Connect(cfg *config.Config) (*sql.DB, error) {
	dsn := cfg.DSN()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Error().Err(err).Msg("Failed to open database connection")
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		log.Error().Err(err).Msg("Failed to ping database")
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Adjust pool based on connection type
	if cfg.DatabaseURL != "" {
		// Supabase / cloud — lower limits
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(3)
		log.Info().Msg("Supabase database connection established successfully")
	} else {
		// Local PostgreSQL — higher limits
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(5)
		log.Info().Msg("Local PostgreSQL connection established successfully")
	}

	return db, nil
}
