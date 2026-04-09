package domain

import "time"

// DemoLink represents a tblDemoLinks_DML row.
type DemoLink struct {
	ID             int64      `gorm:"column:dml_id;primaryKey"`
	TokenHash      string     `gorm:"column:dml_token_hash"`
	AppCode        string     `gorm:"column:dml_app_code"`
	DemoUserID     int64      `gorm:"column:dml_demo_user_id"`
	ExpiresAt      time.Time  `gorm:"column:dml_expires_at"`
	IsActive       bool       `gorm:"column:dml_is_active"`
	RecipientEmail *string    `gorm:"column:dml_recipient_email"`
	CreatedBy      int64      `gorm:"column:dml_created_by"`
	CreatedAt      time.Time  `gorm:"column:dml_created_at"`
}

func (DemoLink) TableName() string { return `tblDemoLinks_DML` }

// ─── Request / Response DTOs ───────────────────────────────────────────────────

// CreateDemoLinkRequest is the body for POST /api/v1/demo-links.
type CreateDemoLinkRequest struct {
	AppCode        string  `json:"appCode"`
	// DemoUserID is optional. When 0 the backend resolves the default demo
	// account for the given application (demo_[appcode]).
	DemoUserID     int64   `json:"demoUserId"`
	ExpiresInHours int     `json:"expiresInHours"`
	RecipientEmail *string `json:"recipientEmail,omitempty"`
	// GuestName is the display name used in the invitation email.
	// For Mode A (system user) it comes from the user's full name.
	// For Mode B (external guest) it is typed manually.
	GuestName      *string `json:"guestName,omitempty"`
}

// DemoLinkResponse is what the API returns.
type DemoLinkResponse struct {
	ID             int64     `json:"id"`
	AppCode        string    `json:"appCode"`
	DemoUserID     int64     `json:"demoUserId"`
	ExpiresAt      time.Time `json:"expiresAt"`
	IsActive       bool      `json:"isActive"`
	RecipientEmail *string   `json:"recipientEmail,omitempty"`
	CreatedBy      int64     `json:"createdBy"`
	CreatedAt      time.Time `json:"createdAt"`
	// Token is only populated on creation; never stored.
	Token *string `json:"token,omitempty"`
}
