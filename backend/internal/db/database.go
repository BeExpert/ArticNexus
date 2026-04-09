package db

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// PoolConfig holds connection pool settings sourced from Config.
type PoolConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// New opens a GORM connection to PostgreSQL using the provided DSN and pool settings.
func New(dsn string, development bool, pool PoolConfig) (*gorm.DB, error) {
	logLevel := logger.Silent
	if development {
		// Warn shows slow queries and errors but not every statement.
		// Use Info only when you need full query tracing.
		logLevel = logger.Warn
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := database.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(pool.MaxOpenConns)
	sqlDB.SetMaxIdleConns(pool.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(pool.ConnMaxLifetime)

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	log.Println("database connection established")
	return database, nil
}
