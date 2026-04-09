package repository

import (
	"articnexus/backend/internal/domain"

	"gorm.io/gorm"
)

// ApplicationRepository defines the persistence contract for Application entities.
type ApplicationRepository interface {
	FindByID(id int64) (*domain.Application, error)
	FindAll(params domain.PaginationParams) ([]domain.Application, int64, error)
	Create(app *domain.Application) error
	Update(app *domain.Application) error
	Delete(id int64) error
}

type applicationRepository struct {
	db *gorm.DB
}

// NewApplicationRepository returns a GORM-backed ApplicationRepository.
func NewApplicationRepository(db *gorm.DB) ApplicationRepository {
	return &applicationRepository{db: db}
}

func (r *applicationRepository) FindByID(id int64) (*domain.Application, error) {
	var app domain.Application
	if err := r.db.First(&app, id).Error; err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *applicationRepository) FindAll(params domain.PaginationParams) ([]domain.Application, int64, error) {
	var apps []domain.Application
	var total int64

	if err := r.db.Model(&domain.Application{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Limit(params.PageSize).Offset(params.Offset()).Find(&apps).Error
	return apps, total, err
}

func (r *applicationRepository) Create(app *domain.Application) error {
	return r.db.Create(app).Error
}

func (r *applicationRepository) Update(app *domain.Application) error {
	return r.db.Save(app).Error
}

func (r *applicationRepository) Delete(id int64) error {
	return r.db.Delete(&domain.Application{}, id).Error
}

// ─────────────────────────────────────────────────────────────────────────────

// ModuleRepository defines the persistence contract for Module entities.
type ModuleRepository interface {
	FindByID(id int64) (*domain.Module, error)
	FindByApplication(appID int64, params domain.PaginationParams) ([]domain.Module, int64, error)
	FindByIDs(ids []int64) ([]domain.Module, error)
	Create(module *domain.Module) error
	Update(module *domain.Module) error
	Delete(id int64) error
	// FindNamesByAppCode returns all active module names for a given application code.
	FindNamesByAppCode(appCode string) ([]string, error)
}

type moduleRepository struct {
	db *gorm.DB
}

// NewModuleRepository returns a GORM-backed ModuleRepository.
func NewModuleRepository(db *gorm.DB) ModuleRepository {
	return &moduleRepository{db: db}
}

func (r *moduleRepository) FindByID(id int64) (*domain.Module, error) {
	var module domain.Module
	if err := r.db.First(&module, id).Error; err != nil {
		return nil, err
	}
	return &module, nil
}

func (r *moduleRepository) FindByApplication(appID int64, params domain.PaginationParams) ([]domain.Module, int64, error) {
	var modules []domain.Module
	var total int64

	base := r.db.Model(&domain.Module{}).Where("app_id = ?", appID)

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := base.Limit(params.PageSize).Offset(params.Offset()).Find(&modules).Error
	return modules, total, err
}

func (r *moduleRepository) FindByIDs(ids []int64) ([]domain.Module, error) {
	var modules []domain.Module
	err := r.db.Where("mod_id IN ?", ids).Find(&modules).Error
	return modules, err
}

func (r *moduleRepository) Create(module *domain.Module) error {
	return r.db.Create(module).Error
}

func (r *moduleRepository) Update(module *domain.Module) error {
	return r.db.Save(module).Error
}

func (r *moduleRepository) Delete(id int64) error {
	return r.db.Delete(&domain.Module{}, id).Error
}

func (r *moduleRepository) FindNamesByAppCode(appCode string) ([]string, error) {
	var names []string
	err := r.db.Raw(`
		SELECT m.mod_name
		FROM "tblModules_MOD" m
		JOIN "tblApplications_APP" a ON a.app_id = m.app_id
		WHERE a.app_code = ? AND m.mod_status = 'active'
		ORDER BY m.mod_name
	`, appCode).Scan(&names).Error
	if err != nil {
		return nil, err
	}
	if names == nil {
		names = []string{}
	}
	return names, nil
}
