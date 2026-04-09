package domain

import (
	"time"

	"gorm.io/gorm"
)

// AuditFields is embedded in every entity model to provide automatic
// created_at / updated_at timestamps and soft-delete via deleted_at.
// GORM fills CreatedAt on INSERT, UpdatedAt on every UPDATE, and
// uses DeletedAt as a soft-delete sentinel (NULL = alive).
type AuditFields struct {
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"          json:"-"`
}
