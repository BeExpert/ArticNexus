package handler

import (
	"net/http"
	"strconv"

	"articnexus/backend/internal/domain"
	"articnexus/backend/internal/middleware"
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
// Accepts optional ?company_id=X query param to filter roles by the apps
// licensed for that company. When absent, falls back to the JWT com_id.
func (h *RoleHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	params := paginationFromQuery(r)
	companyID, _ := middleware.CompanyIDFromContext(r.Context())

	// Allow an explicit company_id query param to override the JWT value.
	// This is used when the super-admin views a specific company's roles.
	if qv := r.URL.Query().Get("company_id"); qv != "" {
		if parsed, err := strconv.ParseInt(qv, 10, 64); err == nil && parsed > 0 {
			companyID = parsed
		}
	}

	roles, total, err := h.roleService.GetAll(companyID, params)
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
	// Inject company scope from JWT so the service can enforce app licensing.
	if comID, ok := middleware.CompanyIDFromContext(r.Context()); ok {
		req.CompanyID = comID
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
