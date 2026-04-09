package db

import (
	"fmt"

	"gorm.io/gorm"
)

// ResetPublicSchema drops and recreates the public schema.
// Intended only for local development workflows where a clean database is required.
func ResetPublicSchema(database *gorm.DB) error {
	sqlDB, err := database.DB()
	if err != nil {
		return fmt.Errorf("could not obtain sql.DB: %w", err)
	}

	if _, err := sqlDB.Exec(`DROP SCHEMA IF EXISTS public CASCADE`); err != nil {
		return fmt.Errorf("drop schema public failed: %w", err)
	}

	if _, err := sqlDB.Exec(`CREATE SCHEMA public`); err != nil {
		return fmt.Errorf("create schema public failed: %w", err)
	}

	if _, err := sqlDB.Exec(`ALTER SCHEMA public OWNER TO CURRENT_USER`); err != nil {
		return fmt.Errorf("set public schema owner failed: %w", err)
	}

	if _, err := sqlDB.Exec(`GRANT ALL ON SCHEMA public TO CURRENT_USER`); err != nil {
		return fmt.Errorf("grant current user on public schema failed: %w", err)
	}

	if _, err := sqlDB.Exec(`GRANT USAGE ON SCHEMA public TO PUBLIC`); err != nil {
		return fmt.Errorf("grant usage on public schema failed: %w", err)
	}

	return nil
}
