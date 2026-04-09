package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"articnexus/backend/internal/domain"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// ─── Response helpers ─────────────────────────────────────────────────────────

func renderJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func renderOK(w http.ResponseWriter, data interface{}, message string) {
	renderJSON(w, http.StatusOK, domain.NewResponse(data, message))
}

func renderCreated(w http.ResponseWriter, data interface{}, message string) {
	renderJSON(w, http.StatusCreated, domain.NewResponse(data, message))
}

func renderList(w http.ResponseWriter, data interface{}, message string, p domain.Pagination) {
	renderJSON(w, http.StatusOK, domain.NewResponseWithPagination(data, message, p))
}

func renderError(w http.ResponseWriter, status int, message string, errs interface{}) {
	renderJSON(w, status, domain.NewErrorResponse(message, errs))
}

// renderErrorCode writes a JSON error response that includes a stable error code.
// Prefer this over renderError for all user-facing errors.
func renderErrorCode(w http.ResponseWriter, status int, code, message string, errs interface{}) {
	renderJSON(w, status, domain.NewErrorResponseCode(message, code, errs))
}

// renderAppError inspects err: if it is a *domain.AppError it uses its Status
// and Code; otherwise it falls back to 500 Internal Server Error.
func renderAppError(w http.ResponseWriter, err error) {
	var appErr *domain.AppError
	if errors.As(err, &appErr) {
		renderErrorCode(w, appErr.Status, appErr.Code, appErr.Message, nil)
		return
	}
	renderInternalError(w)
}

func renderNotFound(w http.ResponseWriter) {
	renderErrorCode(w, http.StatusNotFound, domain.ErrCodeNotFound, "resource not found", nil)
}

func renderBadRequest(w http.ResponseWriter, message string) {
	renderErrorCode(w, http.StatusBadRequest, domain.ErrCodeBadRequest, message, nil)
}

func renderInternalError(w http.ResponseWriter) {
	renderErrorCode(w, http.StatusInternalServerError, domain.ErrCodeInternal, "internal server error", nil)
}

// isNotFound returns true when err is a GORM record-not-found error.
func isNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// ─── Request helpers ──────────────────────────────────────────────────────────

// decodeJSON decodes the JSON request body into v. Returns false and writes a
// 400 error to w if decoding fails.
func decodeJSON(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		renderBadRequest(w, "invalid request body: "+err.Error())
		return false
	}
	return true
}

// urlParamInt64 extracts a named Chi URL parameter and converts it to int64.
// Returns false and writes a 400 error on failure.
func urlParamInt64(w http.ResponseWriter, r *http.Request, name string) (int64, bool) {
	raw := chi.URLParam(r, name)
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		renderBadRequest(w, "invalid "+name+" parameter")
		return 0, false
	}
	return id, true
}

// paginationFromQuery extracts page/pageSize from query string with safe defaults.
func paginationFromQuery(r *http.Request) domain.PaginationParams {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	search := r.URL.Query().Get("search")
	return domain.PaginationParams{Page: page, PageSize: pageSize, Search: search}
}

// buildPagination constructs the Pagination meta struct for list responses.
func buildPagination(params domain.PaginationParams, total int64) domain.Pagination {
	totalPages := int(total) / params.PageSize
	if int(total)%params.PageSize != 0 {
		totalPages++
	}
	return domain.Pagination{
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalItems: total,
		TotalPages: totalPages,
	}
}

// ─── Postgres error helpers ───────────────────────────────────────────────────

// isDuplicate returns true when err wraps a Postgres unique-violation (23505).
func isDuplicate(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}

// isForeignKeyViolation returns true when err wraps a Postgres FK violation (23503).
func isForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23503"
	}
	return false
}

// constraintMessages maps known Postgres constraint names to user-friendly messages.
var constraintMessages = map[string]string{
	"tblusers_usr_usr_username_key":     "Ya existe un usuario con ese nombre de usuario.",
	"tblusers_usr_usr_email_key":        "Ya existe un usuario con ese correo electrónico.",
	"tblcompanies_com_com_name_key":     "Ya existe una empresa con ese nombre.",
	"tblbranches_bra_bra_code_key":      "Ya existe una sucursal con ese código.",
	"tblapplications_app_app_code_key":  "Ya existe una aplicación con ese código.",
	"tblroles_rol_rol_name_key":         "Ya existe un rol con ese nombre.",
	"tbluserroles_uro_pkey":             "Este usuario ya tiene ese rol asignado en esa sucursal.",
	"tblusercompanies_uco_pkey":         "Este usuario ya está asignado a esta empresa.",
	"tblrolemodules_rmo_pkey":           "Este módulo ya está asignado al rol.",
}

// duplicateMessage returns a user-friendly message for a unique-violation error.
func duplicateMessage(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		key := strings.ToLower(pgErr.ConstraintName)
		if msg, ok := constraintMessages[key]; ok {
			return msg
		}
	}
	return "El registro ya existe (duplicado)."
}

// fkViolationMessage returns a friendly message for a foreign-key violation.
func fkViolationMessage() string {
	return "No se puede eliminar: existen registros que dependen de este."
}
