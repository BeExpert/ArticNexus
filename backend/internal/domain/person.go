package domain

import "time"

// Person maps to tblPersons_PER.
type Person struct {
	ID             int64      `gorm:"column:per_id;primaryKey;autoIncrement"`
	FirstName      string     `gorm:"column:per_firstname;not null"`
	FirstSurname   string     `gorm:"column:per_firstsurname;not null"`
	SecondSurname  *string    `gorm:"column:per_secondsurname"`
	NationalID     *string    `gorm:"column:per_nationalid"`
	Email          *string    `gorm:"column:per_email"`
	BirthDate      *time.Time `gorm:"column:per_birthdate"`
	PhoneAreaCode  *string    `gorm:"column:per_phoneareacode"`
	PrimaryPhone   *string    `gorm:"column:per_primaryphone"`
	SecondaryPhone *string    `gorm:"column:per_secondaryphone"`
	Address        *string    `gorm:"column:per_address"`
	Status         string     `gorm:"column:per_status;not null;default:active"`
	AuditFields
}

func (Person) TableName() string { return "tblPersons_PER" }

// ─── DTOs ────────────────────────────────────────────────────────────────────

// CreatePersonRequest is used when creating a new person record.
type CreatePersonRequest struct {
	FirstName      string  `json:"firstName"      validate:"required,max=100"`
	FirstSurname   string  `json:"firstSurname"   validate:"required,max=100"`
	SecondSurname  *string `json:"secondSurname"`
	NationalID     *string `json:"nationalId"`
	Email          *string `json:"email"          validate:"omitempty,email"`
	BirthDate      *string `json:"birthDate"`
	PhoneAreaCode  *string `json:"phoneAreaCode"`
	PrimaryPhone   *string `json:"primaryPhone"`
	SecondaryPhone *string `json:"secondaryPhone"`
	Address        *string `json:"address"`
}

// UpdatePersonRequest is used to update an existing person. All fields are optional.
type UpdatePersonRequest struct {
	FirstName      *string `json:"firstName"      validate:"omitempty,max=100"`
	FirstSurname   *string `json:"firstSurname"   validate:"omitempty,max=100"`
	SecondSurname  *string `json:"secondSurname"`
	NationalID     *string `json:"nationalId"`
	Email          *string `json:"email"          validate:"omitempty,email"`
	BirthDate      *string `json:"birthDate"`
	PhoneAreaCode  *string `json:"phoneAreaCode"`
	PrimaryPhone   *string `json:"primaryPhone"`
	SecondaryPhone *string `json:"secondaryPhone"`
	Address        *string `json:"address"`
	Status         *string `json:"status"`
}

// PersonResponse is the public representation returned by the API.
type PersonResponse struct {
	ID             int64     `json:"id"`
	FirstName      string    `json:"firstName"`
	FirstSurname   string    `json:"firstSurname"`
	SecondSurname  *string   `json:"secondSurname"`
	NationalID     *string   `json:"nationalId"`
	Email          *string   `json:"email"`
	BirthDate      *string   `json:"birthDate"`
	PhoneAreaCode  *string   `json:"phoneAreaCode"`
	PrimaryPhone   *string   `json:"primaryPhone"`
	SecondaryPhone *string   `json:"secondaryPhone"`
	Address        *string   `json:"address"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
