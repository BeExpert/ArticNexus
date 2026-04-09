package handler

import (
	"net/http"

	"articnexus/backend/internal/domain"
	"articnexus/backend/internal/middleware"
	"articnexus/backend/internal/service"
)

// AuthHandler handles /auth/* routes.
type AuthHandler struct {
	authService    service.AuthService
	companyService service.CompanyService
	userService    service.UserService
}

// NewAuthHandler returns a new AuthHandler.
func NewAuthHandler(authService service.AuthService, companyService service.CompanyService, userService service.UserService) *AuthHandler {
	return &AuthHandler{authService: authService, companyService: companyService, userService: userService}
}

// Login godoc
// POST /api/v1/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	resp, err := h.authService.Login(req)
	if err != nil {
		renderAppError(w, err)
		return
	}

	renderOK(w, resp, "login successful")
}

// Logout godoc
// POST /api/v1/auth/logout
// Since we use stateless JWT, logout is client-side (discard the token).
// The endpoint exists for API completeness and future token-blocklist support.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	renderOK(w, nil, "logged out successfully")
}

// Me godoc
// GET /api/v1/auth/me
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		renderErrorCode(w, http.StatusUnauthorized, domain.ErrCodeAuthUnauthenticated, "unauthenticated", nil)
		return
	}

	user, err := h.authService.GetCurrentUser(userID)
	if err != nil {
		if isNotFound(err) {
			renderNotFound(w)
			return
		}
		renderInternalError(w)
		return
	}

	renderOK(w, user, "")
}

// MyCompanies godoc
// GET /api/v1/auth/me/companies
// Returns all companies the authenticated user is assigned to.
// Super-admins receive the full company list.
func (h *AuthHandler) MyCompanies(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		renderErrorCode(w, http.StatusUnauthorized, domain.ErrCodeAuthUnauthenticated, "unauthenticated", nil)
		return
	}

	user, err := h.authService.GetCurrentUser(userID)
	if err != nil {
		if isNotFound(err) {
			renderNotFound(w)
			return
		}
		renderInternalError(w)
		return
	}

	if user.IsSuperAdmin {
		// Super admin sees everything.
		params := domain.PaginationParams{Page: 1, PageSize: 1000}
		companies, _, err := h.companyService.GetAll(params)
		if err != nil {
			renderInternalError(w)
			return
		}
		renderOK(w, companies, "")
		return
	}

	// Regular user: return only assigned companies.
	companies, err := h.companyService.GetUserCompanies(userID)
	if err != nil {
		renderInternalError(w)
		return
	}
	renderOK(w, companies, "")
}

// UpdateMe godoc
// PUT /api/v1/auth/me
func (h *AuthHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		renderErrorCode(w, http.StatusUnauthorized, domain.ErrCodeAuthUnauthenticated, "unauthenticated", nil)
		return
	}

	var req domain.UpdateProfileRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	user, err := h.authService.UpdateProfile(userID, req)
	if err != nil {
		if isDuplicate(err) {
			renderErrorCode(w, http.StatusConflict, domain.ErrCodeConflictDuplicate, duplicateMessage(err), nil)
			return
		}
		renderInternalError(w)
		return
	}

	renderOK(w, user, "perfil actualizado")
}

// ForgotPassword godoc
// POST /api/v1/auth/forgot-password  (public — no JWT required)
func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req domain.ForgotPasswordRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	// Always return 200 to avoid leaking whether the username exists.
	_ = h.authService.ForgotPassword(req.Username)

	renderOK(w, nil, "Si el usuario existe y tiene un correo asociado, recibirás un enlace para restablecer tu contraseña.")
}

// ResetPassword godoc
// POST /api/v1/auth/reset-password  (public — no JWT required)
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req domain.ResetPasswordRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	if err := h.authService.ResetPassword(req.Token, req.NewPassword); err != nil {
		renderAppError(w, err)
		return
	}

	renderOK(w, nil, "Contraseña actualizada correctamente. Ya puedes iniciar sesión.")
}
