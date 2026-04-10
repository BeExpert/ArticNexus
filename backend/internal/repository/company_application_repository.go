package repository

import (
	"articnexus/backend/internal/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CompanyApplicationRepository manages the company ↔ application licensing table.
type CompanyApplicationRepository interface {
	// GetLicensedAppIDs returns the app_ids the company is licensed for (status=active).
	GetLicensedAppIDs(companyID int64) ([]int64, error)
	// IsLicensed returns true if the company holds an active license for the app.
	IsLicensed(companyID, appID int64) (bool, error)
	// Create records a new license entry.
	Create(ca *domain.CompanyApplication) error
	// BulkCreate inserts multiple license entries (ignores duplicates).
	BulkCreate(companyID int64, appIDs []int64) error
	// GetByCompanyID returns all licensed apps for a company with app metadata.
	GetByCompanyID(companyID int64) ([]domain.CompanyAppDetail, error)
	// UpdateStatus sets cap_status for a specific company+app row.
	UpdateStatus(companyID, appID int64, status string) error
	// Delete removes a license entry for a specific company+app.
	Delete(companyID, appID int64) error
	// GetAppCodeByID returns the app_code for the given app_id.
	GetAppCodeByID(appID int64) (string, error)
}

type companyApplicationRepository struct {
	db *gorm.DB
}

func NewCompanyApplicationRepository(db *gorm.DB) CompanyApplicationRepository {
	return &companyApplicationRepository{db: db}
}

func (r *companyApplicationRepository) GetLicensedAppIDs(companyID int64) ([]int64, error) {
	var ids []int64
	err := r.db.
		Model(&domain.CompanyApplication{}).
		Where("com_id = ? AND cap_status = 'active'", companyID).
		Pluck("app_id", &ids).Error
	return ids, err
}

func (r *companyApplicationRepository) IsLicensed(companyID, appID int64) (bool, error) {
	var count int64
	err := r.db.
		Model(&domain.CompanyApplication{}).
		Where("com_id = ? AND app_id = ? AND cap_status = 'active'", companyID, appID).
		Count(&count).Error
	return count > 0, err
}

func (r *companyApplicationRepository) Create(ca *domain.CompanyApplication) error {
	return r.db.Create(ca).Error
}

func (r *companyApplicationRepository) BulkCreate(companyID int64, appIDs []int64) error {
	if len(appIDs) == 0 {
		return nil
	}
	entries := make([]domain.CompanyApplication, len(appIDs))
	for i, aid := range appIDs {
		entries[i] = domain.CompanyApplication{
			CompanyID: companyID,
			AppID:     aid,
			Status:    "active",
		}
	}
	// ON CONFLICT DO NOTHING — safe to call multiple times
	return r.db.
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&entries).Error
}

func (r *companyApplicationRepository) GetByCompanyID(companyID int64) ([]domain.CompanyAppDetail, error) {
	var results []domain.CompanyAppDetail
	err := r.db.
		Table(`"tblCompanyApplications_CAP" cap`).
		Select(`a.app_id, a.app_code, a.app_name, cap.cap_status AS status, cap.cap_created_at AS created_at`).
		Joins(`JOIN "tblApplications_APP" a ON a.app_id = cap.app_id`).
		Where("cap.com_id = ?", companyID).
		Order("a.app_name ASC").
		Scan(&results).Error
	return results, err
}

func (r *companyApplicationRepository) UpdateStatus(companyID, appID int64, status string) error {
	res := r.db.
		Model(&domain.CompanyApplication{}).
		Where("com_id = ? AND app_id = ?", companyID, appID).
		Update("cap_status", status)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *companyApplicationRepository) Delete(companyID, appID int64) error {
	res := r.db.
		Where("com_id = ? AND app_id = ?", companyID, appID).
		Delete(&domain.CompanyApplication{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *companyApplicationRepository) GetAppCodeByID(appID int64) (string, error) {
	var code string
	err := r.db.
		Table(`"tblApplications_APP"`).
		Select("app_code").
		Where("app_id = ?", appID).
		Scan(&code).Error
	return code, err
}
