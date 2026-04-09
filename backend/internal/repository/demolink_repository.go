package repository

import (
	"fmt"

	"articnexus/backend/internal/domain"

	"gorm.io/gorm"
)

// DemoLinkRepository handles persistence for demo links.
type DemoLinkRepository interface {
	Create(link *domain.DemoLink) error
	FindByTokenHash(hash string) (*domain.DemoLink, error)
	ListByApp(appCode string) ([]domain.DemoLink, error)
	ListAll() ([]domain.DemoLink, error)
	Revoke(id int64) error
	// UpdateDemoUserPassword replaces the bcrypt password hash for a demo user.
	// Called each time a demo link is created so the password rotates with every invite.
	UpdateDemoUserPassword(userID int64, passwordHash string) error
	// FindDefaultDemoUserByApp returns the usr_id of the generic demo account
	// for the given application (username = "demo_" + lower(appCode)).
	FindDefaultDemoUserByApp(appCode string) (int64, error)
}

type demoLinkRepository struct {
	db *gorm.DB
}

func NewDemoLinkRepository(db *gorm.DB) DemoLinkRepository {
	return &demoLinkRepository{db: db}
}

func (r *demoLinkRepository) Create(link *domain.DemoLink) error {
	return r.db.Create(link).Error
}

func (r *demoLinkRepository) FindByTokenHash(hash string) (*domain.DemoLink, error) {
	var l domain.DemoLink
	err := r.db.Where("dml_token_hash = ?", hash).First(&l).Error
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *demoLinkRepository) ListByApp(appCode string) ([]domain.DemoLink, error) {
	var links []domain.DemoLink
	err := r.db.Where("dml_app_code = ?", appCode).
		Order("dml_created_at DESC").
		Find(&links).Error
	return links, err
}

func (r *demoLinkRepository) ListAll() ([]domain.DemoLink, error) {
	var links []domain.DemoLink
	err := r.db.Order("dml_created_at DESC").Find(&links).Error
	return links, err
}

func (r *demoLinkRepository) Revoke(id int64) error {
	return r.db.Model(&domain.DemoLink{}).
		Where("dml_id = ?", id).
		Update("dml_is_active", false).Error
}

func (r *demoLinkRepository) UpdateDemoUserPassword(userID int64, passwordHash string) error {
	return r.db.Exec(
		`UPDATE "tblUsers_USR" SET usr_password = ?, updated_at = now() WHERE usr_id = ?`,
		passwordHash, userID,
	).Error
}

func (r *demoLinkRepository) FindDefaultDemoUserByApp(appCode string) (int64, error) {
	var userID int64
	err := r.db.Raw(
		`SELECT usr_id FROM "tblUsers_USR"
		 WHERE usr_username = 'demo_' || lower(?) AND deleted_at IS NULL
		 LIMIT 1`,
		appCode,
	).Scan(&userID).Error
	if err != nil {
		return 0, err
	}
	if userID == 0 {
		return 0, domain.ErrValidation(domain.ErrCodeDemoNoDemoUser,
			fmt.Sprintf("no demo user found for app %q — run seed first", appCode))
	}
	return userID, nil
}
