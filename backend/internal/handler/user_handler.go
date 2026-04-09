package handler

import (
	"net/http"

	"articnexus/backend/internal/domain"
	"articnexus/backend/internal/service"
)

// UserHandler handles /users/* routes.
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler returns a new UserHandler.
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetAll godoc
// GET /api/v1/users
func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	params := paginationFromQuery(r)

	users, total, err := h.userService.GetAll(params)
	if err != nil {
		renderInternalError(w)
		return
	}

	renderList(w, users, "", buildPagination(params, total))
}

// GetByID godoc
// GET /api/v1/users/{id}
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := urlParamInt64(w, r, "id")
	if !ok {
		return
	}

	user, err := h.userService.GetByID(id)
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

// Create godoc
// POST /api/v1/users
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateUserRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	user, err := h.userService.Create(req)
	if err != nil {
		if isDuplicate(err) {
			renderErrorCode(w, http.StatusConflict, domain.ErrCodeConflictDuplicate, duplicateMessage(err), nil)
			return
		}
		renderAppError(w, err)
		return
	}

	renderCreated(w, user, "user created successfully")
}

// Update godoc
// PUT /api/v1/users/{id}
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := urlParamInt64(w, r, "id")
	if !ok {
		return
	}

	var req domain.UpdateUserRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	user, err := h.userService.Update(id, req)
	if err != nil {
		if isNotFound(err) {
			renderNotFound(w)
			return
		}
		if isDuplicate(err) {
			renderErrorCode(w, http.StatusConflict, domain.ErrCodeConflictDuplicate, duplicateMessage(err), nil)
			return
		}
		renderAppError(w, err)
		return
	}

	renderOK(w, user, "user updated successfully")
}

// Delete godoc
// DELETE /api/v1/users/{id}
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := urlParamInt64(w, r, "id")
	if !ok {
		return
	}

	if err := h.userService.Delete(id); err != nil {
		if isNotFound(err) {
			renderNotFound(w)
			return
		}
		if isForeignKeyViolation(err) {
			renderErrorCode(w, http.StatusConflict, domain.ErrCodeConflictDependency, fkViolationMessage(), nil)
			return
		}
		renderInternalError(w)
		return
	}

	renderOK(w, nil, "user deleted successfully")
}

// ResetPassword godoc
// POST /api/v1/users/{id}/reset-password
func (h *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	id, ok := urlParamInt64(w, r, "id")
	if !ok {
		return
	}

	var req domain.AdminResetPasswordRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	resp, err := h.userService.ResetUserPassword(id, req)
	if err != nil {
		if isNotFound(err) {
			renderNotFound(w)
			return
		}
		renderAppError(w, err)
		return
	}

	renderOK(w, resp, "password reset successfully")
}
