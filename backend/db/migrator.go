package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// Migrator handles database migrations similar to Laravel's migration system.
// It tracks which migrations have been run in a `migrations` table.
type Migrator struct {
	db             *sql.DB
	migrationsDir  string
}

// MigrationRecord represents a row in the migrations tracking table.
type MigrationRecord struct {
	ID        int
	Migration string
	Batch     int
	RanAt     time.Time
}

// NewMigrator creates a new Migrator instance.
func NewMigrator(db *sql.DB, migrationsDir string) *Migrator {
	return &Migrator{
		db:            db,
		migrationsDir: migrationsDir,
	}
}

// Initialize creates the migrations tracking table if it doesn't exist.
// This is equivalent to Laravel's migrations table.
func (m *Migrator) Initialize() error {
	query := `
	CREATE TABLE IF NOT EXISTS migrations (
		id SERIAL PRIMARY KEY,
		migration VARCHAR(255) NOT NULL UNIQUE,
		batch INT NOT NULL,
		ran_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`

	_, err := m.db.Exec(query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create migrations table")
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	log.Info().Msg("Migrations table ready")
	return nil
}

// MigrateUp runs all pending "up" migrations that haven't been executed yet.
// Similar to: php artisan migrate
func (m *Migrator) MigrateUp() error {
	if err := m.Initialize(); err != nil {
		return err
	}

	// Get list of already-run migrations
	ran, err := m.getRanMigrations()
	if err != nil {
		return err
	}

	// Get all migration files
	files, err := m.getMigrationFiles("up")
	if err != nil {
		return err
	}

	if len(files) == 0 {
		log.Info().Msg("No migration files found")
		return nil
	}

	// Filter out already-run migrations
	pending := []string{}
	for _, f := range files {
		name := migrationName(f)
		if !ran[name] {
			pending = append(pending, f)
		}
	}

	if len(pending) == 0 {
		log.Info().Msg("Nothing to migrate. All migrations are up to date")
		return nil
	}

	// Get the next batch number
	batch, err := m.getNextBatch()
	if err != nil {
		return err
	}

	// Run each pending migration
	for _, file := range pending {
		name := migrationName(file)
		log.Info().Str("migration", name).Int("batch", batch).Msg("Migrating...")

		sqlContent, err := os.ReadFile(file)
		if err != nil {
			log.Error().Err(err).Str("file", file).Msg("Failed to read migration file")
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		// Execute migration in a transaction
		tx, err := m.db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		if _, err := tx.Exec(string(sqlContent)); err != nil {
			tx.Rollback()
			log.Error().Err(err).Str("migration", name).Msg("Migration failed")
			return fmt.Errorf("migration %s failed: %w", name, err)
		}

		// Record the migration
		if _, err := tx.Exec(
			"INSERT INTO migrations (migration, batch) VALUES ($1, $2)",
			name, batch,
		); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", name, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", name, err)
		}

		log.Info().Str("migration", name).Msg("Migrated successfully ✓")
	}

	log.Info().Int("count", len(pending)).Int("batch", batch).Msg("All migrations completed")
	return nil
}

// MigrateDown rolls back the last batch of migrations.
// Similar to: php artisan migrate:rollback
func (m *Migrator) MigrateDown() error {
	if err := m.Initialize(); err != nil {
		return err
	}

	// Get the last batch number
	var lastBatch int
	err := m.db.QueryRow("SELECT COALESCE(MAX(batch), 0) FROM migrations").Scan(&lastBatch)
	if err != nil {
		return fmt.Errorf("failed to get last batch: %w", err)
	}

	if lastBatch == 0 {
		log.Info().Msg("Nothing to rollback")
		return nil
	}

	// Get migrations from the last batch (in reverse order)
	rows, err := m.db.Query(
		"SELECT migration FROM migrations WHERE batch = $1 ORDER BY id DESC",
		lastBatch,
	)
	if err != nil {
		return fmt.Errorf("failed to get last batch migrations: %w", err)
	}
	defer rows.Close()

	var migrations []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("failed to scan migration name: %w", err)
		}
		migrations = append(migrations, name)
	}

	// Run each "down" migration
	for _, name := range migrations {
		downFile := filepath.Join(m.migrationsDir, name+".down.sql")

		sqlContent, err := os.ReadFile(downFile)
		if err != nil {
			log.Error().Err(err).Str("file", downFile).Msg("Down migration file not found")
			return fmt.Errorf("down migration file not found for %s: %w", name, err)
		}

		log.Info().Str("migration", name).Msg("Rolling back...")

		tx, err := m.db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		if _, err := tx.Exec(string(sqlContent)); err != nil {
			tx.Rollback()
			log.Error().Err(err).Str("migration", name).Msg("Rollback failed")
			return fmt.Errorf("rollback %s failed: %w", name, err)
		}

		// Remove migration record
		if _, err := tx.Exec("DELETE FROM migrations WHERE migration = $1", name); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to remove migration record %s: %w", name, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit rollback %s: %w", name, err)
		}

		log.Info().Str("migration", name).Msg("Rolled back successfully ✓")
	}

	log.Info().Int("count", len(migrations)).Int("batch", lastBatch).Msg("Rollback completed")
	return nil
}

// MigrateReset rolls back ALL migrations.
// Similar to: php artisan migrate:reset
func (m *Migrator) MigrateReset() error {
	if err := m.Initialize(); err != nil {
		return err
	}

	// Get all migrations in reverse order
	rows, err := m.db.Query("SELECT migration FROM migrations ORDER BY id DESC")
	if err != nil {
		return fmt.Errorf("failed to get migrations: %w", err)
	}
	defer rows.Close()

	var migrations []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("failed to scan migration: %w", err)
		}
		migrations = append(migrations, name)
	}

	if len(migrations) == 0 {
		log.Info().Msg("Nothing to reset")
		return nil
	}

	for _, name := range migrations {
		downFile := filepath.Join(m.migrationsDir, name+".down.sql")

		sqlContent, err := os.ReadFile(downFile)
		if err != nil {
			log.Warn().Str("migration", name).Msg("Down file not found, skipping")
			continue
		}

		log.Info().Str("migration", name).Msg("Resetting...")

		if _, err := m.db.Exec(string(sqlContent)); err != nil {
			log.Error().Err(err).Str("migration", name).Msg("Reset failed")
			return fmt.Errorf("reset %s failed: %w", name, err)
		}
	}

	// Clear all migration records
	if _, err := m.db.Exec("DELETE FROM migrations"); err != nil {
		return fmt.Errorf("failed to clear migrations table: %w", err)
	}

	log.Info().Int("count", len(migrations)).Msg("Reset completed")
	return nil
}

// MigrateFresh drops all tables and re-runs all migrations.
// Similar to: php artisan migrate:fresh
func (m *Migrator) MigrateFresh() error {
	log.Warn().Msg("Running fresh migration — dropping all tables...")

	// Drop all tables
	_, err := m.db.Exec(`
		DO $$ DECLARE
			r RECORD;
		BEGIN
			FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') LOOP
				EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(r.tablename) || ' CASCADE';
			END LOOP;
		END $$;
	`)
	if err != nil {
		return fmt.Errorf("failed to drop tables: %w", err)
	}

	log.Info().Msg("All tables dropped")

	// Re-run all migrations
	return m.MigrateUp()
}

// Status prints the status of all migrations.
// Similar to: php artisan migrate:status
func (m *Migrator) Status() ([]map[string]string, error) {
	if err := m.Initialize(); err != nil {
		return nil, err
	}

	ran, err := m.getRanMigrations()
	if err != nil {
		return nil, err
	}

	files, err := m.getMigrationFiles("up")
	if err != nil {
		return nil, err
	}

	var statuses []map[string]string
	for _, f := range files {
		name := migrationName(f)
		status := "Pending"
		if ran[name] {
			status = "Ran"
		}
		statuses = append(statuses, map[string]string{
			"migration": name,
			"status":    status,
		})
	}

	return statuses, nil
}

// getRanMigrations returns a set of migration names that have already been run.
func (m *Migrator) getRanMigrations() (map[string]bool, error) {
	ran := make(map[string]bool)

	rows, err := m.db.Query("SELECT migration FROM migrations ORDER BY id")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get ran migrations")
		return nil, fmt.Errorf("failed to get ran migrations: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("failed to scan migration: %w", err)
		}
		ran[name] = true
	}

	return ran, nil
}

// getNextBatch returns the next batch number.
func (m *Migrator) getNextBatch() (int, error) {
	var batch int
	err := m.db.QueryRow("SELECT COALESCE(MAX(batch), 0) + 1 FROM migrations").Scan(&batch)
	if err != nil {
		return 0, fmt.Errorf("failed to get next batch: %w", err)
	}
	return batch, nil
}

// getMigrationFiles returns sorted migration files of the given type (up/down).
func (m *Migrator) getMigrationFiles(direction string) ([]string, error) {
	pattern := filepath.Join(m.migrationsDir, fmt.Sprintf("*.%s.sql", direction))
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to glob migration files: %w", err)
	}
	sort.Strings(files)
	return files, nil
}

// migrationName extracts the migration name from a file path.
// e.g., "db/migrations/000001_create_tasks_table.up.sql" → "000001_create_tasks_table"
func migrationName(filePath string) string {
	base := filepath.Base(filePath)
	// Remove .up.sql or .down.sql
	name := strings.TrimSuffix(base, ".up.sql")
	name = strings.TrimSuffix(name, ".down.sql")
	return name
}
