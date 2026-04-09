package domain

import "time"

// Module maps to tblModules_MOD.
type Module struct {
	ID            int64   `gorm:"column:mod_id;primaryKey;autoIncrement"`
	ApplicationID int64   `gorm:"column:app_id;not null"`
	Name          string  `gorm:"column:mod_name;not null"`
	DisplayName   *string `gorm:"column:mod_display_name"`
	MenuOption    *string `gorm:"column:mod_menuoption"`
	SubFunction   *string `gorm:"column:mod_subfunction"`
	Description   *string `gorm:"column:mod_description"`
	Status        string  `gorm:"column:mod_status;not null;default:active"`
	AuditFields
}

func (Module) TableName() string { return "tblModules_MOD" }

// ─── DTOs ────────────────────────────────────────────────────────────────────

type CreateModuleRequest struct {
	ApplicationID int64   `json:"applicationId" validate:"required"`
	Name          string  `json:"name"          validate:"required,max=100"`
	DisplayName   *string `json:"displayName"`
	MenuOption    *string `json:"menuOption"`
	SubFunction   *string `json:"subFunction"`
	Description   *string `json:"description"`
}

type UpdateModuleRequest struct {
	Name        *string `json:"name"        validate:"omitempty,max=100"`
	DisplayName *string `json:"displayName"`
	MenuOption  *string `json:"menuOption"`
	SubFunction *string `json:"subFunction"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
}

type ModuleResponse struct {
	ID            int64     `json:"id"`
	ApplicationID int64     `json:"applicationId"`
	Name          string    `json:"name"`
	DisplayName   *string   `json:"displayName"`
	MenuOption    *string   `json:"menuOption"`
	SubFunction   *string   `json:"subFunction"`
	Description   *string   `json:"description"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
