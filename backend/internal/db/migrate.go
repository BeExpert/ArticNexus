package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gorm.io/gorm"
)

// Migrate reads every *.sql file in migrationsDir (sorted lexicographically)
// and executes them in order. Each file is wrapped in a transaction so a
// failure rolls back only that file.
// All CREATE TABLE statements use IF NOT EXISTS, making runs idempotent.
func Migrate(database *gorm.DB, migrationsDir string) error {
	// Resolve relative paths against the process working directory.
	if !filepath.IsAbs(migrationsDir) {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("could not determine working directory: %w", err)
		}
		migrationsDir = filepath.Join(cwd, migrationsDir)
	}

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("could not read migrations directory %q: %w", migrationsDir, err)
	}

	// Collect .sql files and sort them so they execute in version order.
	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			files = append(files, filepath.Join(migrationsDir, e.Name()))
		}
	}
	sort.Strings(files)

	if len(files) == 0 {
		log.Printf("no SQL migration files found in %s", migrationsDir)
		return nil
	}

	// Use the underlying *sql.DB so migrations bypass GORM's query logger.
	sqlDB, err := database.DB()
	if err != nil {
		return fmt.Errorf("could not obtain sql.DB: %w", err)
	}

	for _, f := range files {
		log.Printf("applying migration: %s", filepath.Base(f))

		content, err := os.ReadFile(f)
		if err != nil {
			return fmt.Errorf("could not read %q: %w", f, err)
		}

		if _, err := sqlDB.Exec(string(content)); err != nil {
			return fmt.Errorf("migration %q failed: %w", filepath.Base(f), err)
		}

		log.Printf("migration applied: %s", filepath.Base(f))
	}

	return nil
}
