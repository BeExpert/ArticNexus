package handler

import (
	"net/http"

	"articnexus/backend/internal/domain"
	"articnexus/backend/internal/service"
)

// ApplicationHandler handles /applications/* routes including modules sub-routes.
type ApplicationHandler struct {
	appService service.ApplicationService
}

// NewApplicationHandler returns a new ApplicationHandler.
func NewApplicationHandler(appService service.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{appService: appService}
}

// ─── Applications ─────────────────────────────────────────────────────────────

// GetAll godoc
// GET /api/v1/applications
func (h *ApplicationHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	params := paginationFromQuery(r)
	apps, total, err := h.appService.GetAll(params)
	if err != nil {
		renderInternalError(w)
		return
	}
	renderList(w, apps, "", buildPagination(params, total))
}

// GetByID godoc
// GET /api/v1/applications/{id}
func (h *ApplicationHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := urlParamInt64(w, r, "id")
	if !ok {
		return
	}
	app, err := h.appService.GetByID(id)
	if err != nil {
		if isNotFound(err) {
			renderNotFound(w)
			return
		}
		renderInternalError(w)
		return
	}
	renderOK(w, app, "")
}

// Create godoc
// POST /api/v1/applications
func (h *ApplicationHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateApplicationRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	app, err := h.appService.Create(req)
	if err != nil {
		if isDuplicate(err) {
			renderErrorCode(w, http.StatusConflict, domain.ErrCodeConflictDuplicate, duplicateMessage(err), nil)
			return
		}
		renderAppError(w, err)
		return
	}
	renderCreated(w, app, "application created successfully")
}

// Update godoc
// PUT /api/v1/applications/{id}
func (h *ApplicationHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := urlParamInt64(w, r, "id")
	if !ok {
		return
	}
	var req domain.UpdateApplicationRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	app, err := h.appService.Update(id, req)
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
	renderOK(w, app, "application updated successfully")
}

// Delete godoc
// DELETE /api/v1/applications/{id}
func (h *ApplicationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := urlParamInt64(w, r, "id")
	if !ok {
		return
	}
	if err := h.appService.Delete(id); err != nil {
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
	renderOK(w, nil, "application deleted successfully")
}

// ─── Modules ──────────────────────────────────────────────────────────────────

// GetModules godoc
// GET /api/v1/applications/{id}/modules
func (h *ApplicationHandler) GetModules(w http.ResponseWriter, r *http.Request) {
	appID, ok := urlParamInt64(w, r, "id")
	if !ok {
		return
	}
	params := paginationFromQuery(r)
	modules, total, err := h.appService.GetModules(appID, params)
	if err != nil {
		renderInternalError(w)
		return
	}
	renderList(w, modules, "", buildPagination(params, total))
}

// GetModule godoc
// GET /api/v1/applications/{id}/modules/{moduleId}
func (h *ApplicationHandler) GetModule(w http.ResponseWriter, r *http.Request) {
	moduleID, ok := urlParamInt64(w, r, "moduleId")
	if !ok {
		return
	}
	module, err := h.appService.GetModule(moduleID)
	if err != nil {
		if isNotFound(err) {
			renderNotFound(w)
			return
		}
		renderInternalError(w)
		return
	}
	renderOK(w, module, "")
}

// CreateModule godoc
// POST /api/v1/applications/{id}/modules
func (h *ApplicationHandler) CreateModule(w http.ResponseWriter, r *http.Request) {
	appID, ok := urlParamInt64(w, r, "id")
	if !ok {
		return
	}
	var req domain.CreateModuleRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	req.ApplicationID = appID

	module, err := h.appService.CreateModule(req)
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
	renderCreated(w, module, "module created successfully")
}

// UpdateModule godoc
// PUT /api/v1/applications/{id}/modules/{moduleId}
func (h *ApplicationHandler) UpdateModule(w http.ResponseWriter, r *http.Request) {
	moduleID, ok := urlParamInt64(w, r, "moduleId")
	if !ok {
		return
	}
	var req domain.UpdateModuleRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	module, err := h.appService.UpdateModule(moduleID, req)
	if err != nil {
		if isNotFound(err) {
			renderNotFound(w)
			return
		}
		renderAppError(w, err)
		return
	}
	renderOK(w, module, "module updated successfully")
}

// DeleteModule godoc
// DELETE /api/v1/applications/{id}/modules/{moduleId}
func (h *ApplicationHandler) DeleteModule(w http.ResponseWriter, r *http.Request) {
	moduleID, ok := urlParamInt64(w, r, "moduleId")
	if !ok {
		return
	}
	if err := h.appService.DeleteModule(moduleID); err != nil {
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
	renderOK(w, nil, "module deleted successfully")
}
