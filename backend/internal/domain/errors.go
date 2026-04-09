package domain

import "net/http"

// ─── Error codes ─────────────────────────────────────────────────────────────
//
// These stable string constants are the contract between the backend and
// frontend. The frontend translates them via its i18n dictionary.
// Never rename a constant that is already live — add new ones instead.

const (
	// HTTP generic
	ErrCodeBadRequest          = "http.bad_request"
	ErrCodeUnauthorized        = "http.unauthorized"
	ErrCodeForbidden           = "http.forbidden"
	ErrCodeNotFound            = "http.not_found"
	ErrCodeConflictDuplicate   = "http.conflict_duplicate"
	ErrCodeConflictDependency  = "http.conflict_dependency"
	ErrCodeValidation          = "http.validation"
	ErrCodeTooManyRequests     = "http.too_many_requests"
	ErrCodeInternal            = "http.internal"

	// Authentication — JWT & session
	ErrCodeAuthInvalidCredentials = "auth.invalid_credentials"
	ErrCodeAuthAccountInactive    = "auth.account_inactive"
	ErrCodeAuthDemoExpired        = "auth.demo_expired"
	ErrCodeAuthUnauthenticated    = "auth.unauthenticated"
	ErrCodeAuthMissingHeader      = "auth.missing_header"
	ErrCodeAuthInvalidHeaderFmt   = "auth.invalid_header_format"
	ErrCodeAuthInvalidToken       = "auth.invalid_token"
	ErrCodeAuthSessionExpired     = "auth.session_expired"

	// Password reset
	ErrCodeAuthResetLinkInvalid = "auth.reset_link_invalid"
	ErrCodeAuthResetLinkUsed    = "auth.reset_link_used"
	ErrCodeAuthResetLinkExpired = "auth.reset_link_expired"

	// Business domain — users
	ErrCodeUserAlreadyInCompany = "user.already_in_company"
	ErrCodeUserNotInCompany     = "user.not_in_company"

	// Business domain — demo links
	ErrCodeDemoNoDemoUser = "demolink.no_demo_user"

	// Contact form
	ErrCodeContactInvalidName  = "contact.invalid_name"
	ErrCodeContactInvalidEmail = "contact.invalid_email"
	ErrCodeContactInvalidType  = "contact.invalid_type"
	ErrCodeContactDescTooShort = "contact.desc_too_short"
	ErrCodeContactSendFailed   = "contact.send_failed"
)

// ─── AppError ─────────────────────────────────────────────────────────────────

// AppError is a structured, user-facing error that carries an HTTP status code
// and a stable error code in addition to a human-readable message. Services
// return *AppError for errors that should be translated on the frontend; they
// use plain fmt.Errorf for internal/unexpected errors (which handlers convert to
// 500 Internal Server Error).
type AppError struct {
	Status  int
	Code    string
	Message string
}

// Error implements the error interface.
func (e *AppError) Error() string {
	return e.Message
}

// NewAppError creates a new AppError.
func NewAppError(status int, code, message string) *AppError {
	return &AppError{Status: status, Code: code, Message: message}
}

// ─── Convenience constructors ─────────────────────────────────────────────────

func ErrUnauthorized(code, message string) *AppError {
	return NewAppError(http.StatusUnauthorized, code, message)
}

func ErrForbidden(message string) *AppError {
	return NewAppError(http.StatusForbidden, ErrCodeForbidden, message)
}

func ErrBadRequest(code, message string) *AppError {
	return NewAppError(http.StatusBadRequest, code, message)
}

func ErrValidation(code, message string) *AppError {
	return NewAppError(http.StatusUnprocessableEntity, code, message)
}
