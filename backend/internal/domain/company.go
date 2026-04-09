package domain

import "time"

// Company maps to tblCompanies_COM.
type Company struct {
	ID     int64  `gorm:"column:com_id;primaryKey;autoIncrement"`
	Name   string `gorm:"column:com_name;not null"`
	Status string `gorm:"column:com_status;not null;default:active"`
	AuditFields
}

func (Company) TableName() string { return "tblCompanies_COM" }

// ─── DTOs ────────────────────────────────────────────────────────────────────

type CreateCompanyRequest struct {
	Name   string  `json:"name"   validate:"required,max=255"`
	Status *string `json:"status"`

	// Optional: create an admin user for the company in a single step.
	Admin *CreateCompanyAdminRequest `json:"admin,omitempty"`
}

// CreateCompanyAdminRequest holds the data needed to bootstrap an admin user
// when creating a new company.
type CreateCompanyAdminRequest struct {
	FirstName    string `json:"firstName"    validate:"required,max=100"`
	FirstSurname string `json:"firstSurname" validate:"required,max=100"`
	Username     string `json:"username"     validate:"required,max=100"`
	Email        string `json:"email"        validate:"required,email"`
	Password     string `json:"password"     validate:"required,min=8"`
}

type UpdateCompanyRequest struct {
	Name   *string `json:"name"   validate:"omitempty,max=255"`
	Status *string `json:"status"`
}

type CompanyResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
