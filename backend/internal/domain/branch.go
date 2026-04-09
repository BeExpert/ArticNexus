package domain

import "time"

// Branch maps to tblBranches_BRA.
type Branch struct {
	ID          int64   `gorm:"column:bra_id;primaryKey;autoIncrement"`
	CompanyID   int64   `gorm:"column:com_id;not null"`
	Code        string  `gorm:"column:bra_code;not null"`
	Name        string  `gorm:"column:bra_name;not null"`
	Address     *string `gorm:"column:bra_address"`
	PhoneNumber *string `gorm:"column:bra_phonenumber"`
	Email       *string `gorm:"column:bra_email"`
	Status      string  `gorm:"column:bra_status;not null;default:active"`
	AuditFields
}

func (Branch) TableName() string { return "tblBranches_BRA" }

// ─── DTOs ────────────────────────────────────────────────────────────────────

type CreateBranchRequest struct {
	Code        string  `json:"code"    validate:"required,max=50"`
	Name        string  `json:"name"    validate:"required,max=150"`
	Address     *string `json:"address"`
	PhoneNumber *string `json:"phoneNumber"`
	Email       *string `json:"email"   validate:"omitempty,email"`
}

type UpdateBranchRequest struct {
	Name        *string `json:"name"    validate:"omitempty,max=150"`
	Address     *string `json:"address"`
	PhoneNumber *string `json:"phoneNumber"`
	Email       *string `json:"email"   validate:"omitempty,email"`
	Status      *string `json:"status"`
}

type BranchResponse struct {
	ID          int64     `json:"id"`
	CompanyID   int64     `json:"companyId"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Address     *string   `json:"address"`
	PhoneNumber *string   `json:"phoneNumber"`
	Email       *string   `json:"email"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
