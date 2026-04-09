package repository

import (
	"articnexus/backend/internal/domain"

	"gorm.io/gorm"
)

// CompanyUserRepository handles company-user membership and role assignment.
type CompanyUserRepository interface {
	// Membership.
	FindUsersByCompany(companyID int64) ([]domain.User, error)
	AddUserToCompany(companyID, userID int64) error
	RemoveUserFromCompany(companyID, userID int64) error
	IsUserInCompany(companyID, userID int64) (bool, error)

	// Role assignment.
	FindUserRolesByCompany(companyID int64) ([]domain.UserRole, error)
	FindUserRolesForUser(companyID, userID int64) ([]domain.UserRole, error)
	AssignUserRole(ur domain.UserRole) error
	RemoveUserRole(ur domain.UserRole) error
}

type companyUserRepository struct {
	db *gorm.DB
}

func NewCompanyUserRepository(db *gorm.DB) CompanyUserRepository {
	return &companyUserRepository{db: db}
}

// ─── Membership ──────────────────────────────────────────────────────────────

func (r *companyUserRepository) FindUsersByCompany(companyID int64) ([]domain.User, error) {
	var users []domain.User
	err := r.db.Preload("Person").
		Joins(`JOIN "tblUserCompanies_UCO" uco ON uco.usr_id = "tblUsers_USR".usr_id`).
		Where(`uco.com_id = ?`, companyID).
		Find(&users).Error
	return users, err
}

func (r *companyUserRepository) AddUserToCompany(companyID, userID int64) error {
	uc := domain.UserCompany{UserID: userID, CompanyID: companyID}
	return r.db.Create(&uc).Error
}

func (r *companyUserRepository) RemoveUserFromCompany(companyID, userID int64) error {
	return r.db.
		Where("usr_id = ? AND com_id = ?", userID, companyID).
		Delete(&domain.UserCompany{}).Error
}

func (r *companyUserRepository) IsUserInCompany(companyID, userID int64) (bool, error) {
	var count int64
	err := r.db.Model(&domain.UserCompany{}).
		Where("usr_id = ? AND com_id = ?", userID, companyID).
		Count(&count).Error
	return count > 0, err
}

// ─── Role assignment ─────────────────────────────────────────────────────────

func (r *companyUserRepository) FindUserRolesByCompany(companyID int64) ([]domain.UserRole, error) {
	var urs []domain.UserRole
	err := r.db.Where("com_id = ?", companyID).Find(&urs).Error
	return urs, err
}

func (r *companyUserRepository) FindUserRolesForUser(companyID, userID int64) ([]domain.UserRole, error) {
	var urs []domain.UserRole
	err := r.db.Where("com_id = ? AND usr_id = ?", companyID, userID).Find(&urs).Error
	return urs, err
}

func (r *companyUserRepository) AssignUserRole(ur domain.UserRole) error {
	return r.db.Create(&ur).Error
}

func (r *companyUserRepository) RemoveUserRole(ur domain.UserRole) error {
	return r.db.
		Where("usr_id = ? AND com_id = ? AND bra_id = ? AND rol_id = ?",
			ur.UserID, ur.CompanyID, ur.BranchID, ur.RoleID).
		Delete(&domain.UserRole{}).Error
}
