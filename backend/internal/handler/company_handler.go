package handler

import (
	"net/http"

	"articnexus/backend/internal/domain"
	"articnexus/backend/internal/service"
)

// CompanyHandler handles /companies/* routes including branches sub-routes.
type CompanyHandler struct {
	companyService service.CompanyService
}

// NewCompanyHandler returns a new CompanyHandler.
func NewCompanyHandler(companyService service.CompanyService) *CompanyHandler {
	return &CompanyHandler{companyService: companyService}
}

// ─── Companies ────────────────────────────────────────────────────────────────

// GetAll godoc
// GET /api/v1/companies
func (h *CompanyHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	params := paginationFromQuery(r)
	companies, total, err := h.companyService.GetAll(params)
	if err != nil {
		renderInternalError(w)
		return
	}
	renderList(w, companies, "", buildPagination(params, total))
}

// GetByID godoc
// GET /api/v1/companies/{id}
func (h *CompanyHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := urlParamInt64(w, r, "id")
	if !ok {
		return
	}
	company, err := h.companyService.GetByID(id)
	if err != nil {
		if isNotFound(err) {
			renderNotFound(w)
			return
		}
		renderInternalError(w)
		return
	}
	renderOK(w, company, "")
}

// Create godoc
// POST /api/v1/companies
func (h *CompanyHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateCompanyRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	company, err := h.companyService.Create(req)
	if err != nil {
		if isDuplicate(err) {
			renderErrorCode(w, http.StatusConflict, domain.ErrCodeConflictDuplicate, duplicateMessage(err), nil)
			return
		}
		renderAppError(w, err)
		return
	}
	renderCreated(w, company, "company created successfully")
}

// Update godoc
// PUT /api/v1/companies/{id}
func (h *CompanyHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := urlParamInt64(w, r, "id")
	if !ok {
		return
	}
	var req domain.UpdateCompanyRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	company, err := h.companyService.Update(id, req)
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
	renderOK(w, company, "company updated successfully")
}

// Delete godoc
// DELETE /api/v1/companies/{id}
func (h *CompanyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := urlParamInt64(w, r, "id")
	if !ok {
		return
	}
	if err := h.companyService.Delete(id); err != nil {
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
	renderOK(w, nil, "company deleted successfully")
}

// ─── Branches ─────────────────────────────────────────────────────────────────

// GetBranches godoc
// GET /api/v1/companies/{id}/branches
func (h *CompanyHandler) GetBranches(w http.ResponseWriter, r *http.Request) {
	companyID, ok := urlParamInt64(w, r, "id")
	if !ok {
		return
	}
	params := paginationFromQuery(r)
	branches, total, err := h.companyService.GetBranches(companyID, params)
	if err != nil {
		renderInternalError(w)
		return
	}
	renderList(w, branches, "", buildPagination(params, total))
}

// GetBranch godoc
// GET /api/v1/companies/{id}/branches/{branchId}
func (h *CompanyHandler) GetBranch(w http.ResponseWriter, r *http.Request) {
	branchID, ok := urlParamInt64(w, r, "branchId")
	if !ok {
		return
	}
	branch, err := h.companyService.GetBranch(branchID)
	if err != nil {
		if isNotFound(err) {
			renderNotFound(w)
			return
		}
		renderInternalError(w)
		return
	}
	renderOK(w, branch, "")
}

// CreateBranch godoc
// POST /api/v1/companies/{id}/branches
func (h *CompanyHandler) CreateBranch(w http.ResponseWriter, r *http.Request) {
	companyID, ok := urlParamInt64(w, r, "id")
	if !ok {
		return
	}
	var req domain.CreateBranchRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	branch, err := h.companyService.CreateBranch(companyID, req)
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
	renderCreated(w, branch, "branch created successfully")
}

// UpdateBranch godoc
// PUT /api/v1/companies/{id}/branches/{branchId}
func (h *CompanyHandler) UpdateBranch(w http.ResponseWriter, r *http.Request) {
	branchID, ok := urlParamInt64(w, r, "branchId")
	if !ok {
		return
	}
	var req domain.UpdateBranchRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	branch, err := h.companyService.UpdateBranch(branchID, req)
	if err != nil {
		if isNotFound(err) {
			renderNotFound(w)
			return
		}
		renderAppError(w, err)
		return
	}
	renderOK(w, branch, "branch updated successfully")
}

// DeleteBranch godoc
// DELETE /api/v1/companies/{id}/branches/{branchId}
func (h *CompanyHandler) DeleteBranch(w http.ResponseWriter, r *http.Request) {
	branchID, ok := urlParamInt64(w, r, "branchId")
	if !ok {
		return
	}
	if err := h.companyService.DeleteBranch(branchID); err != nil {
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
	renderOK(w, nil, "branch deleted successfully")
}

// ─── Company Users ────────────────────────────────────────────────────────────

// GetCompanyUsers godoc
// GET /api/v1/companies/{id}/users
func (h *CompanyHandler) GetCompanyUsers(w http.ResponseWriter, r *http.Request) {
	companyID, ok := urlParamInt64(w, r, "id")
	if !ok {
		return
	}
	users, err := h.companyService.GetCompanyUsers(companyID)
	if err != nil {
		if isNotFound(err) {
			renderNotFound(w)
			return
		}
		renderInternalError(w)
		return
	}
	renderOK(w, users, "")
}

// AddUserToCompany godoc
// POST /api/v1/companies/{id}/users
func (h *CompanyHandler) AddUserToCompany(w http.ResponseWriter, r *http.Request) {
	companyID, ok := urlParamInt64(w, r, "id")
	if !ok {
		return
	}
	var req domain.AddUserToCompanyRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	if err := h.companyService.AddUserToCompany(companyID, req); err != nil {
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
	renderCreated(w, nil, "user added to company")
}

// RemoveUserFromCompany godoc
// DELETE /api/v1/companies/{id}/users/{userId}
func (h *CompanyHandler) RemoveUserFromCompany(w http.ResponseWriter, r *http.Request) {
	companyID, ok := urlParamInt64(w, r, "id")
	if !ok {
		return
	}
	userID, ok := urlParamInt64(w, r, "userId")
	if !ok {
		return
	}
	if err := h.companyService.RemoveUserFromCompany(companyID, userID); err != nil {
		renderInternalError(w)
		return
	}
	renderOK(w, nil, "user removed from company")
}

// ─── User Roles (within company) ─────────────────────────────────────────────

// AssignUserRole godoc
// POST /api/v1/companies/{id}/users/{userId}/roles
func (h *CompanyHandler) AssignUserRole(w http.ResponseWriter, r *http.Request) {
	companyID, ok := urlParamInt64(w, r, "id")
	if !ok {
		return
	}
	userID, ok := urlParamInt64(w, r, "userId")
	if !ok {
		return
	}
	var req domain.AssignUserRoleRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	if err := h.companyService.AssignUserRole(companyID, userID, req); err != nil {
		if isDuplicate(err) {
			renderErrorCode(w, http.StatusConflict, domain.ErrCodeConflictDuplicate, duplicateMessage(err), nil)
			return
		}
		renderAppError(w, err)
		return
	}
	renderCreated(w, nil, "role assigned to user")
}

// RemoveUserRole godoc
// DELETE /api/v1/companies/{id}/users/{userId}/roles
func (h *CompanyHandler) RemoveUserRole(w http.ResponseWriter, r *http.Request) {
	companyID, ok := urlParamInt64(w, r, "id")
	if !ok {
		return
	}
	userID, ok := urlParamInt64(w, r, "userId")
	if !ok {
		return
	}
	var req domain.AssignUserRoleRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	if err := h.companyService.RemoveUserRole(companyID, userID, req); err != nil {
		renderInternalError(w)
		return
	}
	renderOK(w, nil, "role removed from user")
}
