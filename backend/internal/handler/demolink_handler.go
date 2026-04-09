package handler

import (
	"net/http"

	"articnexus/backend/internal/domain"
	"articnexus/backend/internal/middleware"
	"articnexus/backend/internal/service"
)

// DemoLinkHandler handles /api/v1/demo-links routes.
type DemoLinkHandler struct {
	svc service.DemoLinkService
}

func NewDemoLinkHandler(svc service.DemoLinkService) *DemoLinkHandler {
	return &DemoLinkHandler{svc: svc}
}

// List godoc
// GET /api/v1/demo-links
func (h *DemoLinkHandler) List(w http.ResponseWriter, r *http.Request) {
	links, err := h.svc.List()
	if err != nil {
		renderInternalError(w)
		return
	}
	renderOK(w, links, "")
}

// Create godoc
// POST /api/v1/demo-links
func (h *DemoLinkHandler) Create(w http.ResponseWriter, r *http.Request) {
	actorID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		renderErrorCode(w, http.StatusUnauthorized, domain.ErrCodeAuthUnauthenticated, "unauthorized", nil)
		return
	}

	var req domain.CreateDemoLinkRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	if req.AppCode == "" {
		renderBadRequest(w, "appCode es requerido")
		return
	}

	resp, err := h.svc.Create(actorID, req)
	if err != nil {
		renderAppError(w, err)
		return
	}
	renderCreated(w, resp, "link de demo creado")
}

// Revoke godoc
// DELETE /api/v1/demo-links/{id}
func (h *DemoLinkHandler) Revoke(w http.ResponseWriter, r *http.Request) {
	id, ok := urlParamInt64(w, r, "id")
	if !ok {
		return
	}

	if err := h.svc.Revoke(id); err != nil {
		if isNotFound(err) {
			renderNotFound(w)
			return
		}
		renderInternalError(w)
		return
	}
	renderOK(w, nil, "link de demo revocado")
}
