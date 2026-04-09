package handler

import (
	"net/http"

	"gorm.io/gorm"
)

// StatsHandler returns aggregate counts for the dashboard.
type StatsHandler struct {
	db *gorm.DB
}

func NewStatsHandler(db *gorm.DB) *StatsHandler {
	return &StatsHandler{db: db}
}

type statsResponse struct {
	Users        int64 `json:"users"`
	Companies    int64 `json:"companies"`
	Applications int64 `json:"applications"`
	Roles        int64 `json:"roles"`
}

// GetStats returns entity counts.
// GET /api/v1/stats
func (h *StatsHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	var s statsResponse
	h.db.Raw(`SELECT COUNT(*) FROM "tblUsers_USR" WHERE deleted_at IS NULL`).Scan(&s.Users)
	h.db.Raw(`SELECT COUNT(*) FROM "tblCompanies_COM" WHERE deleted_at IS NULL`).Scan(&s.Companies)
	h.db.Raw(`SELECT COUNT(*) FROM "tblApplications_APP" WHERE deleted_at IS NULL`).Scan(&s.Applications)
	h.db.Raw(`SELECT COUNT(*) FROM "tblRoles_ROL" WHERE deleted_at IS NULL`).Scan(&s.Roles)
	renderOK(w, s, "")
}
