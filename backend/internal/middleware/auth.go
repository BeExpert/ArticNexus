package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"articnexus/backend/internal/domain"
	"articnexus/backend/pkg/logger"

	"github.com/golang-jwt/jwt/v5"
)

// contextKey is an unexported type for context keys in this package.
type contextKey string

const UserIDKey contextKey = "userID"
const companyIDKey contextKey = "companyID"

// Authenticate is a Chi middleware that validates the Bearer JWT token in the
// Authorization header, checks the session epoch to detect stale tokens from
// before the last backend restart, and injects the userID into the request context.
func Authenticate(jwtSecret, sessionEpoch string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.Warn(logger.Security, "auth rejected: missing authorization header")
				writeUnauthorized(w, domain.ErrCodeAuthMissingHeader, "missing authorization header")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				logger.Warn(logger.Security, "auth rejected: invalid authorization header format")
				writeUnauthorized(w, domain.ErrCodeAuthInvalidHeaderFmt, "invalid authorization header format")
				return
			}

			tokenStr := parts[1]

			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return []byte(jwtSecret), nil
			})
			if err != nil || !token.Valid {
				logger.Warn(logger.Security, "auth rejected: invalid or expired token")
				writeUnauthorized(w, domain.ErrCodeAuthInvalidToken, "invalid or expired token")
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				logger.Warn(logger.Security, "auth rejected: invalid token claims")
				writeUnauthorized(w, domain.ErrCodeAuthInvalidToken, "invalid token claims")
				return
			}

			// Validate session epoch: tokens from before the last backend restart
			// will carry a different epoch and must be rejected.
			if sessionEpoch != "" {
				epochClaim, _ := claims["epoch"].(string)
				if epochClaim != sessionEpoch {
					logger.Warn(logger.Security, "auth rejected: session epoch mismatch (backend restarted)")
					writeUnauthorized(w, domain.ErrCodeAuthSessionExpired, "session expired — please log in again")
					return
				}
			}

			// "sub" is stored as float64 in JWT MapClaims.
			subRaw, ok := claims["sub"]
			if !ok {
				logger.Warn(logger.Security, "auth rejected: missing subject in token")
				writeUnauthorized(w, domain.ErrCodeAuthInvalidToken, "missing subject in token")
				return
			}
			var userID int64
			switch v := subRaw.(type) {
			case float64:
				userID = int64(v)
			case int64:
				userID = v
			default:
				logger.Warn(logger.Security, "auth rejected: invalid subject type in token")
				writeUnauthorized(w, domain.ErrCodeAuthInvalidToken, "invalid subject type in token")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)

			// Extract optional company_id claim; defaults to 0 (super-admin / not-yet-selected).
			var companyID int64
			if comRaw, ok2 := claims["com_id"]; ok2 {
				switch v := comRaw.(type) {
				case float64:
					companyID = int64(v)
				case int64:
					companyID = v
				}
			}
			ctx = context.WithValue(ctx, companyIDKey, companyID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// UserIDFromContext retrieves the authenticated user's ID from the request context.
func UserIDFromContext(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(UserIDKey).(int64)
	return id, ok
}

// CompanyIDFromContext retrieves the company ID embedded in the JWT from the
// request context. Returns 0 (and false) when the token carries no company
// scope (super-admin or multi-company user that hasn't selected a company yet).
func CompanyIDFromContext(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(companyIDKey).(int64)
	return id, ok
}

func writeUnauthorized(w http.ResponseWriter, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	resp := domain.NewErrorResponseCode(message, code, nil)
	encodeJSON(w, resp)
}
