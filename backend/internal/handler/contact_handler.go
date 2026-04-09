package handler

import (
	"net/http"
	"strings"
	"unicode/utf8"

	"articnexus/backend/internal/domain"
	"articnexus/backend/internal/service"
)

// ContactRequest holds the fields submitted via the public contact form.
type ContactRequest struct {
	Nombre      string `json:"nombre"`
	Email       string `json:"email"`
	Tipo        string `json:"tipo"`
	Descripcion string `json:"descripcion"`
	Honeypot    string `json:"honeypot"`
}

type ContactHandler struct {
	emailSvc service.EmailService
}

func NewContactHandler(emailSvc service.EmailService) *ContactHandler {
	return &ContactHandler{emailSvc: emailSvc}
}

func (h *ContactHandler) Submit(w http.ResponseWriter, r *http.Request) {
	var req ContactRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	// Honeypot check — bots fill this field
	if req.Honeypot != "" {
		// Silently succeed to not reveal detection
		renderOK(w, nil, "solicitud recibida")
		return
	}

	// Sanitise and validate
	req.Nombre = strings.TrimSpace(req.Nombre)
	req.Email = strings.TrimSpace(req.Email)
	req.Tipo = strings.TrimSpace(req.Tipo)
	req.Descripcion = strings.TrimSpace(req.Descripcion)

	if utf8.RuneCountInString(req.Nombre) < 2 {
		renderErrorCode(w, http.StatusUnprocessableEntity, domain.ErrCodeContactInvalidName, "nombre inválido", nil)
		return
	}
	if !strings.Contains(req.Email, "@") || len(req.Email) > 150 {
		renderErrorCode(w, http.StatusUnprocessableEntity, domain.ErrCodeContactInvalidEmail, "correo inválido", nil)
		return
	}
	validTipos := map[string]bool{
		"Proponer un nuevo proyecto":     true,
		"Unirme a un proyecto existente": true,
		"Solicitud de proyecto":          true,
	}
	if !validTipos[req.Tipo] {
		renderErrorCode(w, http.StatusUnprocessableEntity, domain.ErrCodeContactInvalidType, "tipo de solicitud inválido", nil)
		return
	}
	if utf8.RuneCountInString(req.Descripcion) < 20 {
		renderErrorCode(w, http.StatusUnprocessableEntity, domain.ErrCodeContactDescTooShort, "descripción demasiado corta", nil)
		return
	}

	if err := h.emailSvc.SendContactForm(req.Nombre, req.Email, req.Tipo, req.Descripcion); err != nil {
		renderErrorCode(w, http.StatusInternalServerError, domain.ErrCodeContactSendFailed, "no se pudo enviar la solicitud", nil)
		return
	}

	renderOK(w, nil, "solicitud recibida correctamente")
}
