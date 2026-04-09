package logger

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

// App logs general application events (HTTP requests, startup, seeds, cleanup).
// Security logs authentication and authorization events.
// DB logs database connectivity and query errors.
var (
	App      *log.Logger
	Security *log.Logger
	DB       *log.Logger
)

// Init initialises the three global loggers.  Each logger writes to both a
// rotating file under storage/logs/ and to stderr simultaneously.  It is safe
// to call Init multiple times; subsequent calls replace the previous loggers.
func Init(maxSizeMB, maxBackups int) error {
	baseDir := filepath.Join("storage", "logs")
	// When the binary is run from the project root the path resolves relative
	// to the working directory.  The backend/ sub-directory variant is checked
	// as a fallback so the logger works regardless of cwd.
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		baseDir = filepath.Join("backend", "storage", "logs")
	}
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return err
	}

	App = newLogger(filepath.Join(baseDir, "app.log"), maxSizeMB, maxBackups, "APP")
	Security = newLogger(filepath.Join(baseDir, "security.log"), maxSizeMB, maxBackups, "SEC")
	DB = newLogger(filepath.Join(baseDir, "db.log"), maxSizeMB, maxBackups, "DB")
	return nil
}

// newLogger creates a log.Logger that fans out to a rotating file and stderr.
func newLogger(path string, maxSizeMB, maxBackups int, tag string) *log.Logger {
	w := &lumberjack.Logger{
		Filename:   path,
		MaxSize:    maxSizeMB,
		MaxBackups: maxBackups,
		MaxAge:     0,       // no age-based deletion
		Compress:   false,
	}
	mw := io.MultiWriter(w, os.Stderr)
	prefix := "[" + tag + "] "
	return log.New(mw, prefix, log.Ldate|log.Ltime|log.Lmicroseconds)
}

// Info logs an informational message to the given logger.
func Info(l *log.Logger, msg string) {
	if l != nil {
		l.Printf("[INFO] %s", msg)
	}
}

// Warn logs a warning message to the given logger.
func Warn(l *log.Logger, msg string) {
	if l != nil {
		l.Printf("[WARN] %s", msg)
	}
}

// Error logs an error message to the given logger.
func Error(l *log.Logger, msg string) {
	if l != nil {
		l.Printf("[ERROR] %s", msg)
	}
}
