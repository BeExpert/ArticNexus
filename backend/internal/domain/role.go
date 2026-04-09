package domain

import "time"

// Role maps to tblRoles_ROL.
type Role struct {
	ID            int64  `gorm:"column:rol_id;primaryKey;autoIncrement"`
	ApplicationID int64  `gorm:"column:app_id;not null"`
	Name          string `gorm:"column:rol_name;not null"`
	Status        string `gorm:"column:rol_status;not null;default:active"`
	AuditFields

	// Preloadable associations.
	Application Application `gorm:"foreignKey:app_id;references:app_id"`
	Modules     []Module    `gorm:"many2many:tblRoleModules_RMO;joinForeignKey:rol_id;joinReferences:mod_id"`
}

func (Role) TableName() string { return "tblRoles_ROL" }

// ─── DTOs ────────────────────────────────────────────────────────────────────

type CreateRoleRequest struct {
	ApplicationID int64  `json:"applicationId" validate:"required"`
	Name          string `json:"name"          validate:"required,max=100"`
}

type UpdateRoleRequest struct {
	Name   *string `json:"name"   validate:"omitempty,max=100"`
	Status *string `json:"status"`
}

type RoleResponse struct {
	ID            int64     `json:"id"`
	ApplicationID int64     `json:"applicationId"`
	AppName       string    `json:"appName"`
	Name          string    `json:"name"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// AssignModulesRequest is used to assign or replace the module list of a role.
type AssignModulesRequest struct {
	ModuleIDs []int64 `json:"moduleIds" validate:"required"`
}
