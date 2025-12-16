package database

import (
	"embed"
	"fmt"
	"log"
	"sort"
	"strings"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// RunMigrations executes all pending migrations
func RunMigrations() error {
	if DB == nil {
		return fmt.Errorf("database connection not established")
	}

	// Create migrations tracking table
	err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`).Error
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get list of migration files
	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Filter and sort up migrations
	var upMigrations []string
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".up.sql") {
			upMigrations = append(upMigrations, entry.Name())
		}
	}
	sort.Strings(upMigrations)

	// Run each migration
	for _, filename := range upMigrations {
		version := strings.TrimSuffix(filename, ".up.sql")

		// Check if already applied
		var count int64
		DB.Raw("SELECT COUNT(*) FROM schema_migrations WHERE version = ?", version).Scan(&count)
		if count > 0 {
			log.Printf("Migration %s already applied, skipping", version)
			continue
		}

		// Read and execute migration
		content, err := migrationsFS.ReadFile("migrations/" + filename)
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", filename, err)
		}

		log.Printf("Running migration: %s", filename)
		if err := DB.Exec(string(content)).Error; err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", filename, err)
		}

		// Record migration
		if err := DB.Exec("INSERT INTO schema_migrations (version) VALUES (?)", version).Error; err != nil {
			return fmt.Errorf("failed to record migration %s: %w", filename, err)
		}

		log.Printf("Migration %s completed", filename)
	}

	log.Println("All migrations completed successfully")
	return nil
}

// RollbackMigration rolls back the last migration
func RollbackMigration() error {
	if DB == nil {
		return fmt.Errorf("database connection not established")
	}

	// Get last applied migration
	var version string
	err := DB.Raw("SELECT version FROM schema_migrations ORDER BY applied_at DESC LIMIT 1").Scan(&version).Error
	if err != nil {
		return fmt.Errorf("failed to get last migration: %w", err)
	}

	if version == "" {
		log.Println("No migrations to rollback")
		return nil
	}

	// Read and execute down migration
	downFile := version + ".down.sql"
	content, err := migrationsFS.ReadFile("migrations/" + downFile)
	if err != nil {
		return fmt.Errorf("failed to read rollback migration %s: %w", downFile, err)
	}

	log.Printf("Rolling back migration: %s", version)
	if err := DB.Exec(string(content)).Error; err != nil {
		return fmt.Errorf("failed to execute rollback %s: %w", downFile, err)
	}

	// Remove migration record
	if err := DB.Exec("DELETE FROM schema_migrations WHERE version = ?", version).Error; err != nil {
		return fmt.Errorf("failed to remove migration record %s: %w", version, err)
	}

	log.Printf("Rollback %s completed", version)
	return nil
}
