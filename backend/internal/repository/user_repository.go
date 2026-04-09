package repository

import (
	"articnexus/backend/internal/domain"

	"gorm.io/gorm"
)

// UserRepository defines the persistence contract for User entities.
type UserRepository interface {
	FindByID(id int64) (*domain.User, error)
	FindByUsername(username string) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	FindAll(params domain.PaginationParams) ([]domain.User, int64, error)
	Create(user *domain.User) error
	Update(user *domain.User) error
	Delete(id int64) error
	// FindUserPermissions returns the distinct module names a user has
	// access to for a given application code (via UserRole → Role → RoleModule → Module).
	FindUserPermissions(userID int64, appCode string) ([]string, error)
	// FindCompaniesByUserID returns the companies a user belongs to.
	FindCompaniesByUserID(userID int64) ([]domain.Company, error)
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository returns a GORM-backed UserRepository.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByID(id int64) (*domain.User, error) {
	var user domain.User
	err := r.db.Preload("Person").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.db.Preload("Person").
		Where("usr_username = ?", username).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Preload("Person").
		Where("usr_email = ?", email).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindAll(params domain.PaginationParams) ([]domain.User, int64, error) {
	var users []domain.User
	var total int64

	query := r.db.Model(&domain.User{})

	if params.Search != "" {
		like := "%" + params.Search + "%"
		query = query.Joins(`JOIN "tblPersons_PER" p ON p.per_id = "tblUsers_USR".per_id`).
			Where(`"tblUsers_USR".usr_username ILIKE ? OR "tblUsers_USR".usr_email ILIKE ? OR p.per_firstname ILIKE ? OR p.per_firstsurname ILIKE ?`,
				like, like, like, like)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	fetchQuery := r.db.Preload("Person")
	if params.Search != "" {
		like := "%" + params.Search + "%"
		fetchQuery = fetchQuery.Joins(`JOIN "tblPersons_PER" p ON p.per_id = "tblUsers_USR".per_id`).
			Where(`"tblUsers_USR".usr_username ILIKE ? OR "tblUsers_USR".usr_email ILIKE ? OR p.per_firstname ILIKE ? OR p.per_firstsurname ILIKE ?`,
				like, like, like, like)
	}

	err := fetchQuery.
		Limit(params.PageSize).
		Offset(params.Offset()).
		Find(&users).Error

	return users, total, err
}

func (r *userRepository) FindCompaniesByUserID(userID int64) ([]domain.Company, error) {
	var companies []domain.Company
	err := r.db.Raw(`
		SELECT c.com_id, c.com_name, c.com_status, c.created_at, c.updated_at
		FROM "tblCompanies_COM" c
		JOIN "tblUserCompanies_UCO" uc ON uc.com_id = c.com_id
		WHERE uc.usr_id = ?
		ORDER BY c.com_name
	`, userID).Scan(&companies).Error
	return companies, err
}

func (r *userRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id int64) error {
	return r.db.Delete(&domain.User{}, id).Error
}

func (r *userRepository) FindUserPermissions(userID int64, appCode string) ([]string, error) {
	var names []string
	err := r.db.Raw(`
		SELECT DISTINCT m.mod_name
		FROM "tblUserRoles_URO"   ur
		JOIN "tblRoles_ROL"       ro ON ro.rol_id  = ur.rol_id
		JOIN "tblApplications_APP" a ON a.app_id   = ro.app_id
		JOIN "tblRoleModules_RMO" rm ON rm.rol_id  = ro.rol_id
		JOIN "tblModules_MOD"      m ON m.mod_id   = rm.mod_id
		WHERE ur.usr_id = ? AND a.app_code = ? AND m.mod_status = 'active'
		ORDER BY m.mod_name
	`, userID, appCode).Scan(&names).Error
	if err != nil {
		return nil, err
	}
	if names == nil {
		names = []string{}
	}
	return names, nil
}
