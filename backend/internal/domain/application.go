package domain

import "time"

// Application maps to tblApplications_APP.
type Application struct {
	ID     int64  `gorm:"column:app_id;primaryKey;autoIncrement"`
	Code   string `gorm:"column:app_code;not null;uniqueIndex"`
	Name   string `gorm:"column:app_name;not null"`
	Status string `gorm:"column:app_status;not null;default:active"`
	AuditFields
}

func (Application) TableName() string { return "tblApplications_APP" }

// ─── DTOs ────────────────────────────────────────────────────────────────────

type CreateApplicationRequest struct {
	Code string `json:"code" validate:"required,max=50"`
	Name string `json:"name" validate:"required,max=255"`
}

type UpdateApplicationRequest struct {
	Name   *string `json:"name"   validate:"omitempty,max=255"`
	Status *string `json:"status"`
}

type ApplicationResponse struct {
	ID        int64     `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
