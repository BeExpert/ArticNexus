package handler

import (
	"net/http"

	"articnexus/backend/internal/domain"
	"articnexus/backend/internal/service"
)

// RoleHandler handles /roles/* routes including module assignment sub-routes.
type RoleHandler struct {
	roleService service.RoleService
}

// NewRoleHandler returns a new RoleHandler.
func NewRoleHandler(roleService service.RoleService) *RoleHandler {
	return &RoleHandler{roleService: roleService}
}

// ─── Roles ────────────────────────────────────────────────────────────────────

// GetAll godoc
// GET /api/v1/roles
func (h *RoleHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	params := paginationFromQuery(r)
	roles, total, err := h.roleService.GetAll(params)
	if err != nil {
		renderInternalError(w)
		return
	}
	renderList(w, roles, "", buildPagination(params, total))
}

// GetByID godoc
// GET /api/v1/roles/{id}
func (h *RoleHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := urlParamInt64(w, r, "id")
	if !ok {
		return
	}
	role, err := h.roleService.GetByID(id)
	if err != nil {
		if isNotFound(err) {
			renderNotFound(w)
			return
		}
		renderInternalError(w)
		return
	}
	renderOK(w, role, "")
}

// Create godoc
// POST /api/v1/roles
func (h *RoleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateRoleRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	role, err := h.roleService.Create(req)
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
	renderCreated(w, role, "role created successfully")
}

// Update godoc
// PUT /api/v1/roles/{id}
func (h *RoleHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := urlParamInt64(w, r, "id")
	if !ok {
		return
	}
	var req domain.UpdateRoleRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	role, err := h.roleService.Update(id, req)
	if err != nil {
		if isNotFound(err) {
			renderNotFound(w)
			return
		}
		renderAppError(w, err)
		return
	}
	renderOK(w, role, "role updated successfully")
}

// Delete godoc
// DELETE /api/v1/roles/{id}
func (h *RoleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := urlParamInt64(w, r, "id")
	if !ok {
		return
	}
	if err := h.roleService.Delete(id); err != nil {
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
	renderOK(w, nil, "role deleted successfully")
}

// ─── Module assignments ───────────────────────────────────────────────────────

// GetModules godoc
// GET /api/v1/roles/{id}/modules
func (h *RoleHandler) GetModules(w http.ResponseWriter, r *http.Request) {
	roleID, ok := urlParamInt64(w, r, "id")
	if !ok {
		return
	}
	modules, err := h.roleService.GetModules(roleID)
	if err != nil {
		if isNotFound(err) {
			renderNotFound(w)
			return
		}
		renderInternalError(w)
		return
	}
	renderOK(w, modules, "")
}

// AssignModules godoc
// POST /api/v1/roles/{id}/modules
func (h *RoleHandler) AssignModules(w http.ResponseWriter, r *http.Request) {
	roleID, ok := urlParamInt64(w, r, "id")
	if !ok {
		return
	}
	var req domain.AssignModulesRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	if err := h.roleService.AssignModules(roleID, req); err != nil {
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
	renderOK(w, nil, "modules assigned successfully")
}

// RemoveModules godoc
// DELETE /api/v1/roles/{id}/modules
func (h *RoleHandler) RemoveModules(w http.ResponseWriter, r *http.Request) {
	roleID, ok := urlParamInt64(w, r, "id")
	if !ok {
		return
	}
	var req domain.AssignModulesRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	if err := h.roleService.RemoveModules(roleID, req); err != nil {
		if isNotFound(err) {
			renderNotFound(w)
			return
		}
		renderAppError(w, err)
		return
	}
	renderOK(w, nil, "modules removed successfully")
}
