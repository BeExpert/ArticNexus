package repository

import (
	"articnexus/backend/internal/domain"

	"gorm.io/gorm"
)

// RoleRepository defines the persistence contract for Role entities.
type RoleRepository interface {
	FindByID(id int64) (*domain.Role, error)
	FindAll(params domain.PaginationParams) ([]domain.Role, int64, error)
	// FindAllByAppIDs filters roles to only those belonging to the given app IDs.
	FindAllByAppIDs(appIDs []int64, params domain.PaginationParams) ([]domain.Role, int64, error)
	FindByApplication(appID int64, params domain.PaginationParams) ([]domain.Role, int64, error)
	Create(role *domain.Role) error
	Update(role *domain.Role) error
	Delete(id int64) error
	// Module assignments.
	AssignModules(roleID int64, moduleIDs []int64) error
	RemoveModules(roleID int64, moduleIDs []int64) error
	FindModules(roleID int64) ([]domain.Module, error)
}

type roleRepository struct {
	db *gorm.DB
}

// NewRoleRepository returns a GORM-backed RoleRepository.
func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) FindByID(id int64) (*domain.Role, error) {
	var role domain.Role
	err := r.db.Preload("Modules").Preload("Application").First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) FindAll(params domain.PaginationParams) ([]domain.Role, int64, error) {
	var roles []domain.Role
	var total int64

	if err := r.db.Model(&domain.Role{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Preload("Application").Limit(params.PageSize).Offset(params.Offset()).Find(&roles).Error
	return roles, total, err
}

func (r *roleRepository) FindAllByAppIDs(appIDs []int64, params domain.PaginationParams) ([]domain.Role, int64, error) {
	var roles []domain.Role
	var total int64

	base := r.db.Model(&domain.Role{}).Where("app_id IN ?", appIDs)
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := base.Preload("Application").Limit(params.PageSize).Offset(params.Offset()).Find(&roles).Error
	return roles, total, err
}

func (r *roleRepository) FindByApplication(appID int64, params domain.PaginationParams) ([]domain.Role, int64, error) {
	var roles []domain.Role
	var total int64

	base := r.db.Model(&domain.Role{}).Where("app_id = ?", appID)

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := base.Limit(params.PageSize).Offset(params.Offset()).Find(&roles).Error
	return roles, total, err
}

func (r *roleRepository) Create(role *domain.Role) error {
	return r.db.Create(role).Error
}

func (r *roleRepository) Update(role *domain.Role) error {
	return r.db.Save(role).Error
}

func (r *roleRepository) Delete(id int64) error {
	return r.db.Delete(&domain.Role{}, id).Error
}

func (r *roleRepository) AssignModules(roleID int64, moduleIDs []int64) error {
	records := make([]domain.RoleModule, len(moduleIDs))
	for i, mid := range moduleIDs {
		records[i] = domain.RoleModule{RoleID: roleID, ModuleID: mid}
	}
	return r.db.CreateInBatches(records, 100).Error
}

func (r *roleRepository) RemoveModules(roleID int64, moduleIDs []int64) error {
	return r.db.
		Where("rol_id = ? AND mod_id IN ?", roleID, moduleIDs).
		Delete(&domain.RoleModule{}).Error
}

func (r *roleRepository) FindModules(roleID int64) ([]domain.Module, error) {
	var modules []domain.Module
	err := r.db.
		Joins("JOIN \"tblRoleModules_RMO\" rmo ON rmo.mod_id = \"tblModules_MOD\".mod_id").
		Where("rmo.rol_id = ?", roleID).
		Find(&modules).Error
	return modules, err
}
