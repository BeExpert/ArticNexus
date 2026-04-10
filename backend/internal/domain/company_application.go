package domain

import "time"

// CompanyApplication maps to tblCompanyApplications_CAP.
// It records which applications a given company has licensed.
type CompanyApplication struct {
	ID        int64     `gorm:"column:cap_id;primaryKey;autoIncrement"`
	CompanyID int64     `gorm:"column:com_id;not null"`
	AppID     int64     `gorm:"column:app_id;not null"`
	Status    string    `gorm:"column:cap_status;not null;default:active"`
	CreatedAt time.Time `gorm:"column:cap_created_at;autoCreateTime"`
}

func (CompanyApplication) TableName() string { return "tblCompanyApplications_CAP" }

// CompanyAppDetail is returned by GET /companies/{id}/applications.
// It includes the application metadata alongside the license status.
type CompanyAppDetail struct {
	AppID     int64     `json:"appId"     gorm:"column:app_id"`
	AppCode   string    `json:"appCode"   gorm:"column:app_code"`
	AppName   string    `json:"appName"   gorm:"column:app_name"`
	Status    string    `json:"status"    gorm:"column:status"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
}

// AddCompanyAppRequest is the body for POST /companies/{id}/applications.
type AddCompanyAppRequest struct {
	AppID int64 `json:"appId" validate:"required"`
}

// UpdateCompanyAppStatusRequest is the body for PATCH /companies/{id}/applications/{appId}.
type UpdateCompanyAppStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=active inactive"`
}
