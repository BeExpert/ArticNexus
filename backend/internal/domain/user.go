package domain

import "time"

// User maps to tblUsers_USR.
// Note: super-admin status is NOT stored in the database — it is derived at
// runtime by comparing the username to the SUPER_ADMIN_USER env variable.
type User struct {
	ID                int64      `gorm:"column:usr_id;primaryKey;autoIncrement"`
	PersonID          int64      `gorm:"column:per_id;not null"`
	Username          string     `gorm:"column:usr_username;not null;uniqueIndex"`
	Email             string     `gorm:"column:usr_email;not null;uniqueIndex"`
	Password          string     `gorm:"column:usr_password;not null"`
	Status            string     `gorm:"column:usr_status;not null;default:active"`
	PasswordExpiresAt *time.Time `gorm:"column:usr_password_expires_at"`
	AuditFields

	// Preloadable associations.
	Person Person `gorm:"foreignKey:per_id;references:per_id"`
}

func (User) TableName() string { return "tblUsers_USR" }

// ─── DTOs ────────────────────────────────────────────────────────────────────

// CreateUserRequest bundles person + user fields for account creation.
type CreateUserRequest struct {
	// Person fields.
	FirstName      string  `json:"firstName"    validate:"required,max=100"`
	FirstSurname   string  `json:"firstSurname" validate:"required,max=100"`
	SecondSurname  *string `json:"secondSurname"`
	NationalID     *string `json:"nationalId"`
	BirthDate      *string `json:"birthDate"`
	PhoneAreaCode  *string `json:"phoneAreaCode"`
	PrimaryPhone   *string `json:"primaryPhone"`
	SecondaryPhone *string `json:"secondaryPhone"`
	Address        *string `json:"address"`

	// User fields.
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"omitempty,min=8"`

	// When true the backend generates a random password and e-mails the
	// credentials to the user. The Password field is ignored in this case.
	SendCredentials bool `json:"sendCredentials"`

	// Optional: assign user to a company on creation.
	CompanyID *int64 `json:"companyId"`
}

// UpdateUserRequest allows partial updates of a user account.
type UpdateUserRequest struct {
	// User fields.
	Username          *string    `json:"username" validate:"omitempty,max=100"`
	Email             *string    `json:"email"    validate:"omitempty,email"`
	Password          *string    `json:"password" validate:"omitempty,min=8"`
	Status            *string    `json:"status"`
	PasswordExpiresAt *time.Time `json:"passwordExpiresAt"`
	ClearExpiry       bool       `json:"clearExpiry"`

	// Person fields (optional — only sent fields are updated).
	FirstName      *string `json:"firstName"      validate:"omitempty,max=100"`
	FirstSurname   *string `json:"firstSurname"   validate:"omitempty,max=100"`
	SecondSurname  *string `json:"secondSurname"`
	NationalID     *string `json:"nationalId"`
	BirthDate      *string `json:"birthDate"`
	PhoneAreaCode  *string `json:"phoneAreaCode"`
	PrimaryPhone   *string `json:"primaryPhone"`
	SecondaryPhone *string `json:"secondaryPhone"`
	Address        *string `json:"address"`
}

// UserResponse is the public user representation — password is never included.
type UserResponse struct {
	ID                int64      `json:"id"`
	Username          string     `json:"username"`
	Email             string     `json:"email"`
	Status            string     `json:"status"`
	PasswordExpiresAt *time.Time `json:"passwordExpiresAt"`
	IsSuperAdmin      bool       `json:"isSuperAdmin"`
	Permissions       []string          `json:"permissions"`
	CreatedAt         time.Time         `json:"createdAt"`
	UpdatedAt         time.Time         `json:"updatedAt"`
	Person            PersonResponse    `json:"person"`
	Companies         []CompanyResponse `json:"companies"`
}

// LoginRequest holds credentials for the /auth/login endpoint.
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse is the payload returned on successful authentication.
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// UpdateProfileRequest is used by PUT /auth/me — no password or status changes.
type UpdateProfileRequest struct {
	Email          *string `json:"email"    validate:"omitempty,email"`
	FirstName      *string `json:"firstName"      validate:"omitempty,max=100"`
	FirstSurname   *string `json:"firstSurname"   validate:"omitempty,max=100"`
	SecondSurname  *string `json:"secondSurname"`
	NationalID     *string `json:"nationalId"`
	BirthDate      *string `json:"birthDate"`
	PhoneAreaCode  *string `json:"phoneAreaCode"`
	PrimaryPhone   *string `json:"primaryPhone"`
	SecondaryPhone *string `json:"secondaryPhone"`
	Address        *string `json:"address"`
}

// ─── Password Reset ──────────────────────────────────────────────────────────

// PasswordResetToken maps to tblPasswordResetTokens_PRT.
type PasswordResetToken struct {
	ID        int64     `gorm:"column:prt_id;primaryKey;autoIncrement"`
	UserID    int64     `gorm:"column:usr_id;not null"`
	TokenHash string    `gorm:"column:prt_token_hash;not null"`
	ExpiresAt time.Time `gorm:"column:prt_expires_at;not null"`
	Used      bool      `gorm:"column:prt_used;default:false"`
	CreatedAt time.Time `gorm:"column:prt_created_at;autoCreateTime"`
}

func (PasswordResetToken) TableName() string { return "tblPasswordResetTokens_PRT" }

// ForgotPasswordRequest is the payload for POST /auth/forgot-password.
type ForgotPasswordRequest struct {
	Username string `json:"username" validate:"required"`
}

// ResetPasswordRequest is the payload for POST /auth/reset-password.
type ResetPasswordRequest struct {
	Token       string `json:"token"       validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,min=8"`
}

// AdminResetPasswordRequest is the payload for POST /users/{id}/reset-password.
// Either GenerateRandom must be true OR Password must be provided.
type AdminResetPasswordRequest struct {
	Password       *string `json:"password" validate:"omitempty,min=8"`
	GenerateRandom bool    `json:"generateRandom"`
}

// AdminResetPasswordResponse is returned after a successful admin password reset.
type AdminResetPasswordResponse struct {
	// GeneratedPassword is non-empty only when GenerateRandom was true.
	// It is the plain-text password — shown once, never stored.
	GeneratedPassword string `json:"generatedPassword,omitempty"`
}
