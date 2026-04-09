package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all runtime configuration values loaded from environment variables.
// Variable names follow the .env convention used by this project:
//
//	APP_PORT, APP_ENV
//	DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SSL_MODE
//	DB_MAX_OPEN_CONNS, DB_MAX_IDLE_CONNS, DB_CONN_MAX_LIFETIME_MIN
//	JWT_SECRET, JWT_EXP_HOURS
//	ALLOWED_ORIGINS
//	MIGRATIONS_DIR
//	SUPER_ADMIN_USER, SUPER_ADMIN_PASS, SUPER_ADMIN_FORCE
//	BCRYPT_COST
type Config struct {
	// HTTP server.
	Port        string
	Environment string

	// Database DSN (built from DB_* vars or overridden by DATABASE_URL).
	DatabaseURL    string
	DBMaxOpenConns int
	DBMaxIdleConns int
	DBConnMaxLifetime time.Duration

	// JWT.
	JWTSecret  string
	JWTExpHours int

	// CORS.
	AllowedOrigins []string

	// Migrations.
	MigrationsDir string

	// Bootstrap super admin.
	SuperAdminUser  string
	SuperAdminPass  string
	SuperAdminForce bool

	// Bcrypt.
	BcryptCost int

	// Shared SMTP transport (both accounts use Gmail).
	SMTPHost string
	SMTPPort string

	// SMTP – soporte técnico: restablecimiento de contraseña y tickets.
	// Env vars: SUPPORT_SMTP_USER, SUPPORT_SMTP_PASSWORD, SUPPORT_SMTP_FROM
	SupportSMTPUser     string
	SupportSMTPPassword string
	SupportSMTPFrom     string

	// SMTP – empresa: formulario de contacto y propuestas.
	// Env vars: BUSINESS_SMTP_USER, BUSINESS_SMTP_PASSWORD, BUSINESS_SMTP_FROM
	BusinessSMTPUser     string
	BusinessSMTPPassword string
	BusinessSMTPFrom     string

	// Frontend URL for building reset links.
	FrontendURL string

	// Product URLs for building demo invitation links.
	// OFTADATA_URL — base URL of the OftaData app (e.g. https://oftadata.articdev.com)
	// VETDATA_URL  — base URL of the VetData app  (e.g. https://vetdata.articdev.com)
	OftaDataURL string
	VetDataURL  string

	// Destination email for contact-form submissions.
	ContactEmail string

	// Password reset token expiration (minutes).
	PasswordResetExpMin int

	// Logging configuration.
	// LOG_MAX_SIZE_MB: max megabytes per log file before rotation (default 10).
	// LOG_MAX_BACKUPS: number of rotated log files to retain (default 3).
	MaxLogSizeMB  int
	MaxLogBackups int

	// SessionEpoch is generated fresh on every boot via crypto/rand.
	// JWT tokens issued before the current boot are automatically rejected.
	// This field is NOT read from the environment — it is set in main.go.
	SessionEpoch string
}

// Load reads the .env file (if present) and populates a Config.
// It is safe to call in production where no .env file exists.
func Load() (*Config, error) {
	_ = godotenv.Load() // non-fatal: absent in production

	cfg := &Config{
		Port:        getEnv("APP_PORT", getEnv("PORT", "8080")),
		Environment: getEnv("APP_ENV", getEnv("ENVIRONMENT", "development")),

		DBMaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 10),
		DBConnMaxLifetime: time.Duration(getEnvInt("DB_CONN_MAX_LIFETIME_MIN", 60)) * time.Minute,

		JWTSecret:   os.Getenv("JWT_SECRET"),
		JWTExpHours: getEnvInt("JWT_EXP_HOURS", 24),

		MigrationsDir: getEnv("MIGRATIONS_DIR", "../database/migrations"),

		SuperAdminUser:  getEnv("SUPER_ADMIN_USER", ""),
		SuperAdminPass:  getEnv("SUPER_ADMIN_PASS", ""),
		SuperAdminForce: getEnv("SUPER_ADMIN_FORCE", "0") == "1",

		BcryptCost: getEnvInt("BCRYPT_COST", 12),

		SMTPHost: getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort: getEnv("SMTP_PORT", "587"),

		SupportSMTPUser:     getEnv("SUPPORT_SMTP_USER", ""),
		SupportSMTPPassword: getEnv("SUPPORT_SMTP_PASSWORD", ""),
		SupportSMTPFrom:     getEnv("SUPPORT_SMTP_FROM", ""),

		BusinessSMTPUser:     getEnv("BUSINESS_SMTP_USER", ""),
		BusinessSMTPPassword: getEnv("BUSINESS_SMTP_PASSWORD", ""),
		BusinessSMTPFrom:     getEnv("BUSINESS_SMTP_FROM", ""),

		ContactEmail: getEnv("CONTACT_EMAIL", ""),

		FrontendURL:         getEnv("FRONTEND_URL", "http://localhost:5173"),
		OftaDataURL:         getEnv("OFTADATA_URL", "http://localhost:5174"),
		VetDataURL:          getEnv("VETDATA_URL", "http://localhost:5175"),
		PasswordResetExpMin: getEnvInt("PASSWORD_RESET_EXP_MIN", 30),

		MaxLogSizeMB:  getEnvInt("LOG_MAX_SIZE_MB", 10),
		MaxLogBackups: getEnvInt("LOG_MAX_BACKUPS", 3),
	}

	// Prefer an explicit DATABASE_URL; otherwise build from individual DB_* vars.
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		cfg.DatabaseURL = dsn
	} else {
		cfg.DatabaseURL = buildDSN()
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf(
			"database connection is not configured: set DATABASE_URL or the DB_HOST/DB_PORT/DB_USER/DB_PASSWORD/DB_NAME variables",
		)
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	rawOrigins := getEnv("ALLOWED_ORIGINS", "http://localhost:5173")
	for _, o := range strings.Split(rawOrigins, ",") {
		if trimmed := strings.TrimSpace(o); trimmed != "" {
			cfg.AllowedOrigins = append(cfg.AllowedOrigins, trimmed)
		}
	}

	return cfg, nil
}

// IsDevelopment returns true when the service is running in development mode.
func (c *Config) IsDevelopment() bool {
	return strings.ToLower(c.Environment) == "development"
}

// AppBaseURL returns the base URL of the given application code.
// Used to build demo invitation links that point to the correct product.
func (c *Config) AppBaseURL(appCode string) string {
	switch strings.ToUpper(appCode) {
	case "OFTADATA":
		return c.OftaDataURL
	case "VETDATA":
		return c.VetDataURL
	default:
		return c.FrontendURL
	}
}

// buildDSN constructs a PostgreSQL DSN from individual DB_* environment variables.
func buildDSN() string {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := getEnv("DB_SSL_MODE", "disable")

	if user == "" || dbname == "" {
		return ""
	}

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}
