package middleware

import (
	"fmt"
	"net/http"
	"time"

	"articnexus/backend/pkg/logger"
)

// responseWriter wraps http.ResponseWriter to capture the status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Logger is a Chi-compatible middleware that logs each request's method,
// path, status code, and elapsed time to the App logger.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := newResponseWriter(w)

		next.ServeHTTP(wrapped, r)

		dur := time.Since(start).Round(time.Millisecond)
		logger.Info(logger.App, fmt.Sprintf(
			"method=%s path=%s status=%d dur=%s ip=%s",
			r.Method,
			r.RequestURI,
			wrapped.statusCode,
			dur,
			r.RemoteAddr,
		))
	})
}
