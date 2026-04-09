package repository

import (
	"articnexus/backend/internal/domain"

	"gorm.io/gorm"
)

// PasswordResetRepository handles persistence for password reset tokens.
type PasswordResetRepository interface {
	Create(token *domain.PasswordResetToken) error
	FindByTokenHash(hash string) (*domain.PasswordResetToken, error)
	MarkUsed(id int64) error
	InvalidateAllForUser(userID int64) error
}

type passwordResetRepository struct {
	db *gorm.DB
}

func NewPasswordResetRepository(db *gorm.DB) PasswordResetRepository {
	return &passwordResetRepository{db: db}
}

func (r *passwordResetRepository) Create(token *domain.PasswordResetToken) error {
	return r.db.Create(token).Error
}

func (r *passwordResetRepository) FindByTokenHash(hash string) (*domain.PasswordResetToken, error) {
	var t domain.PasswordResetToken
	err := r.db.Where("prt_token_hash = ?", hash).First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *passwordResetRepository) MarkUsed(id int64) error {
	return r.db.Model(&domain.PasswordResetToken{}).Where("prt_id = ?", id).Update("prt_used", true).Error
}

func (r *passwordResetRepository) InvalidateAllForUser(userID int64) error {
	return r.db.Model(&domain.PasswordResetToken{}).
		Where("usr_id = ? AND prt_used = false", userID).
		Update("prt_used", true).Error
}
