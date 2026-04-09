package repository

import (
	"articnexus/backend/internal/domain"

	"gorm.io/gorm"
)

// CompanyRepository defines the persistence contract for Company entities.
type CompanyRepository interface {
	FindByID(id int64) (*domain.Company, error)
	FindAll(params domain.PaginationParams) ([]domain.Company, int64, error)
	FindByUserID(userID int64) ([]domain.Company, error)
	Create(company *domain.Company) error
	Update(company *domain.Company) error
	Delete(id int64) error
}

type companyRepository struct {
	db *gorm.DB
}

// NewCompanyRepository returns a GORM-backed CompanyRepository.
func NewCompanyRepository(db *gorm.DB) CompanyRepository {
	return &companyRepository{db: db}
}

func (r *companyRepository) FindByID(id int64) (*domain.Company, error) {
	var company domain.Company
	if err := r.db.First(&company, id).Error; err != nil {
		return nil, err
	}
	return &company, nil
}

func (r *companyRepository) FindAll(params domain.PaginationParams) ([]domain.Company, int64, error) {
	var companies []domain.Company
	var total int64

	if err := r.db.Model(&domain.Company{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Limit(params.PageSize).Offset(params.Offset()).Find(&companies).Error
	return companies, total, err
}

func (r *companyRepository) Create(company *domain.Company) error {
	return r.db.Create(company).Error
}

func (r *companyRepository) Update(company *domain.Company) error {
	return r.db.Save(company).Error
}

func (r *companyRepository) Delete(id int64) error {
	return r.db.Delete(&domain.Company{}, id).Error
}

// FindByUserID returns all companies linked to the given user via tblUserCompanies_UCO.
func (r *companyRepository) FindByUserID(userID int64) ([]domain.Company, error) {
	var companies []domain.Company
	err := r.db.
		Joins(`JOIN "tblUserCompanies_UCO" uco ON uco.com_id = "tblCompanies_COM".com_id`).
		Where(`uco.usr_id = ?`, userID).
		Find(&companies).Error
	return companies, err
}

// BranchRepository defines the persistence contract for Branch entities.
type BranchRepository interface {
	FindByID(id int64) (*domain.Branch, error)
	FindByCompany(companyID int64, params domain.PaginationParams) ([]domain.Branch, int64, error)
	Create(branch *domain.Branch) error
	Update(branch *domain.Branch) error
	Delete(id int64) error
}

type branchRepository struct {
	db *gorm.DB
}

// NewBranchRepository returns a GORM-backed BranchRepository.
func NewBranchRepository(db *gorm.DB) BranchRepository {
	return &branchRepository{db: db}
}

func (r *branchRepository) FindByID(id int64) (*domain.Branch, error) {
	var branch domain.Branch
	if err := r.db.First(&branch, id).Error; err != nil {
		return nil, err
	}
	return &branch, nil
}

func (r *branchRepository) FindByCompany(companyID int64, params domain.PaginationParams) ([]domain.Branch, int64, error) {
	var branches []domain.Branch
	var total int64

	base := r.db.Model(&domain.Branch{}).Where("com_id = ?", companyID)

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := base.Limit(params.PageSize).Offset(params.Offset()).Find(&branches).Error
	return branches, total, err
}

func (r *branchRepository) Create(branch *domain.Branch) error {
	return r.db.Create(branch).Error
}

func (r *branchRepository) Update(branch *domain.Branch) error {
	return r.db.Save(branch).Error
}

func (r *branchRepository) Delete(id int64) error {
	return r.db.Delete(&domain.Branch{}, id).Error
}
