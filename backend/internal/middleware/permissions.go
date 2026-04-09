package middleware

import (
	"context"
	"fmt"
	"net/http"

	"articnexus/backend/internal/domain"
	"articnexus/backend/internal/repository"
	"articnexus/backend/pkg/logger"
)

// Context keys for the permission layer.
const (
	permissionsKey  contextKey = "permissions"
	isSuperAdminKey contextKey = "isSuperAdmin"
)

// LoadPermissions is a Chi middleware that, after JWT authentication has
// already set the userID in the context, loads the list of module names the
// user is allowed to access. Super-admins (identified by username matching
// superAdminUser) receive all active module names.
func LoadPermissions(
	userRepo repository.UserRepository,
	moduleRepo repository.ModuleRepository,
	superAdminUser string,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := UserIDFromContext(r.Context())
			if !ok {
				writeUnauthorized(w, domain.ErrCodeAuthUnauthenticated, "missing user identity")
				return
			}

			user, err := userRepo.FindByID(userID)
			if err != nil {
				writeUnauthorized(w, domain.ErrCodeAuthUnauthenticated, "user not found")
				return
			}

			ctx := r.Context()

			if user.Username == superAdminUser {
				// Super-admin gets every active module.
				names, err := moduleRepo.FindNamesByAppCode("ARTICNEXUS")
				if err != nil {
					logger.Error(logger.Security, fmt.Sprintf("error loading modules for super-admin: %v", err))
					names = []string{}
				}
				ctx = context.WithValue(ctx, permissionsKey, names)
				ctx = context.WithValue(ctx, isSuperAdminKey, true)
			} else {
				names, err := userRepo.FindUserPermissions(userID, "ARTICNEXUS")
				if err != nil {
					logger.Error(logger.Security, fmt.Sprintf("error loading permissions for user %d: %v", userID, err))
					names = []string{}
				}
				ctx = context.WithValue(ctx, permissionsKey, names)
				ctx = context.WithValue(ctx, isSuperAdminKey, false)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireModule returns a Chi middleware that checks whether the current
// user has at least ONE of the given module names in their permissions
// (OR logic). Super-admins always pass. Regular users without any of the
// required modules receive a 403 Forbidden.
func RequireModule(modules ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if IsSuperAdminFromContext(r.Context()) {
				next.ServeHTTP(w, r)
				return
			}

			perms := PermissionsFromContext(r.Context())
			for _, required := range modules {
				for _, p := range perms {
					if p == required {
						next.ServeHTTP(w, r)
						return
					}
				}
			}

			logger.Warn(logger.Security, fmt.Sprintf("access denied: modules required=%v path=%s", modules, r.URL.Path))
			writeForbidden(w, domain.ErrCodeForbidden, "no tienes acceso al módulo requerido")
		})
	}
}

// PermissionsFromContext returns the module names stored in the request context.
func PermissionsFromContext(ctx context.Context) []string {
	perms, _ := ctx.Value(permissionsKey).([]string)
	return perms
}

// IsSuperAdminFromContext returns true when the current user is the super-admin.
func IsSuperAdminFromContext(ctx context.Context) bool {
	v, _ := ctx.Value(isSuperAdminKey).(bool)
	return v
}

func writeForbidden(w http.ResponseWriter, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	resp := domain.NewErrorResponseCode(message, code, nil)
	encodeJSON(w, resp)
}
