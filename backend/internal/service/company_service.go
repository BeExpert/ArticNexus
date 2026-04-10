package service

import (
	"fmt"

	"articnexus/backend/internal/domain"
	"articnexus/backend/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// CompanyService handles business logic for companies and their branches.
type CompanyService interface {
	GetByID(id int64) (*domain.CompanyResponse, error)
	GetAll(params domain.PaginationParams) ([]domain.CompanyResponse, int64, error)
	GetUserCompanies(userID int64) ([]domain.CompanyResponse, error)
	Create(req domain.CreateCompanyRequest) (*domain.CompanyResponse, error)
	Update(id int64, req domain.UpdateCompanyRequest) (*domain.CompanyResponse, error)
	Delete(id int64) error

	// Branch operations.
	GetBranch(branchID int64) (*domain.BranchResponse, error)
	GetBranches(companyID int64, params domain.PaginationParams) ([]domain.BranchResponse, int64, error)
	CreateBranch(companyID int64, req domain.CreateBranchRequest) (*domain.BranchResponse, error)
	UpdateBranch(branchID int64, req domain.UpdateBranchRequest) (*domain.BranchResponse, error)
	DeleteBranch(branchID int64) error

	// Company-user operations.
	GetCompanyUsers(companyID int64) ([]domain.CompanyUserResponse, error)
	AddUserToCompany(companyID int64, req domain.AddUserToCompanyRequest) error
	RemoveUserFromCompany(companyID, userID int64) error

	// User-role operations.
	AssignUserRole(companyID, userID int64, req domain.AssignUserRoleRequest) error
	RemoveUserRole(companyID, userID int64, req domain.AssignUserRoleRequest) error
}

type companyService struct {
	db              *gorm.DB
	companyRepo     repository.CompanyRepository
	branchRepo      repository.BranchRepository
	companyUserRepo repository.CompanyUserRepository
	roleRepo        repository.RoleRepository
	personRepo      repository.PersonRepository
	userRepo        repository.UserRepository
	companyAppRepo  repository.CompanyApplicationRepository
	superAdminUser  string
}

// NewCompanyService returns a CompanyService implementation.
func NewCompanyService(
	db *gorm.DB,
	companyRepo repository.CompanyRepository,
	branchRepo repository.BranchRepository,
	companyUserRepo repository.CompanyUserRepository,
	roleRepo repository.RoleRepository,
	personRepo repository.PersonRepository,
	userRepo repository.UserRepository,
	companyAppRepo repository.CompanyApplicationRepository,
	superAdminUser string,
) CompanyService {
	return &companyService{
		db:              db,
		companyRepo:     companyRepo,
		branchRepo:      branchRepo,
		companyUserRepo: companyUserRepo,
		roleRepo:        roleRepo,
		personRepo:      personRepo,
		userRepo:        userRepo,
		companyAppRepo:  companyAppRepo,
		superAdminUser:  superAdminUser,
	}
}

func (s *companyService) GetByID(id int64) (*domain.CompanyResponse, error) {
	company, err := s.companyRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	resp := mapCompanyToResponse(company)
	return &resp, nil
}

func (s *companyService) GetAll(params domain.PaginationParams) ([]domain.CompanyResponse, int64, error) {
	companies, total, err := s.companyRepo.FindAll(params)
	if err != nil {
		return nil, 0, err
	}
	responses := make([]domain.CompanyResponse, len(companies))
	for i := range companies {
		responses[i] = mapCompanyToResponse(&companies[i])
	}
	return responses, total, nil
}

func (s *companyService) GetUserCompanies(userID int64) ([]domain.CompanyResponse, error) {
	companies, err := s.companyRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	responses := make([]domain.CompanyResponse, len(companies))
	for i := range companies {
		responses[i] = mapCompanyToResponse(&companies[i])
	}
	return responses, nil
}

func (s *companyService) Create(req domain.CreateCompanyRequest) (*domain.CompanyResponse, error) {
	status := "active"
	if req.Status != nil && *req.Status != "" {
		status = *req.Status
	}

	var result *domain.Company

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		companyRepo := repository.NewCompanyRepository(tx)
		branchRepo := repository.NewBranchRepository(tx)

		// 1. Create the company.
		company := &domain.Company{Name: req.Name, Status: status}
		if err := companyRepo.Create(company); err != nil {
			return fmt.Errorf("could not create company: %w", err)
		}

		// 2. Create a default headquarters branch.
		branchCode := fmt.Sprintf("%s-001", abbreviate(req.Name))
		branch := &domain.Branch{
			CompanyID: company.ID,
			Code:      branchCode,
			Name:      "Casa Matriz",
			Status:    "active",
		}
		if err := branchRepo.Create(branch); err != nil {
			return fmt.Errorf("could not create default branch: %w", err)
		}

		// 3. Optionally bootstrap an admin user.
		if req.Admin != nil {
			personRepo := repository.NewPersonRepository(tx)
			userRepo := repository.NewUserRepository(tx)
			companyUserRepo := repository.NewCompanyUserRepository(tx)

			hash, err := bcrypt.GenerateFromPassword([]byte(req.Admin.Password), bcrypt.DefaultCost)
			if err != nil {
				return fmt.Errorf("could not hash password: %w", err)
			}

			person := &domain.Person{
				FirstName:    req.Admin.FirstName,
				FirstSurname: req.Admin.FirstSurname,
				Status:       "active",
			}
			if err := personRepo.Create(person); err != nil {
				return fmt.Errorf("could not create person: %w", err)
			}

			user := &domain.User{
				PersonID: person.ID,
				Username: req.Admin.Username,
				Email:    req.Admin.Email,
				Password: string(hash),
				Status:   "active",
			}
			if err := userRepo.Create(user); err != nil {
				return fmt.Errorf("could not create user: %w", err)
			}

			// Link user to company.
			if err := companyUserRepo.AddUserToCompany(company.ID, user.ID); err != nil {
				return fmt.Errorf("could not link user to company: %w", err)
			}

			// Assign "Administrador de Empresa" role.
			var role domain.Role
			if err := tx.Where("rol_name = ? AND app_id = (SELECT app_id FROM \"tblApplications_APP\" WHERE app_code = 'ARTICNEXUS')",
				"Administrador de Empresa").First(&role).Error; err == nil {
				ur := domain.UserRole{
					UserID:    user.ID,
					CompanyID: company.ID,
					BranchID:  branch.ID,
					RoleID:    role.ID,
				}
				if err := companyUserRepo.AssignUserRole(ur); err != nil {
					return fmt.Errorf("could not assign admin role: %w", err)
				}
			}
		}

		result = company
		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	// License the company for its requested applications.
	// Always done outside the transaction (no rollback risk).
	appIDs := req.ApplicationIDs
	if len(appIDs) == 0 {
		// Default: license for ARTICNEXUS itself.
		var articApp struct{ AppID int64 }
		if err := s.db.Raw(`SELECT app_id FROM "tblApplications_APP" WHERE app_code = 'ARTICNEXUS' LIMIT 1`).Scan(&articApp).Error; err == nil && articApp.AppID > 0 {
			appIDs = []int64{articApp.AppID}
		}
	}
	if len(appIDs) > 0 {
		_ = s.companyAppRepo.BulkCreate(result.ID, appIDs)
	}

	resp := mapCompanyToResponse(result)
	return &resp, nil
}

// abbreviate returns a short uppercase prefix from a company name (max 3 chars).
func abbreviate(name string) string {
	if len(name) <= 3 {
		result := ""
		for _, c := range name {
			if c >= 'a' && c <= 'z' {
				c -= 32
			}
			result += string(c)
		}
		return result
	}
	result := ""
	for _, c := range name[:3] {
		if c >= 'a' && c <= 'z' {
			c -= 32
		}
		result += string(c)
	}
	return result
}

func (s *companyService) Update(id int64, req domain.UpdateCompanyRequest) (*domain.CompanyResponse, error) {
	company, err := s.companyRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		company.Name = *req.Name
	}
	if req.Status != nil {
		company.Status = *req.Status
	}
	if err := s.companyRepo.Update(company); err != nil {
		return nil, err
	}
	resp := mapCompanyToResponse(company)
	return &resp, nil
}

func (s *companyService) Delete(id int64) error {
	return s.companyRepo.Delete(id)
}

// ─── Branch operations ────────────────────────────────────────────────────────

func (s *companyService) GetBranch(branchID int64) (*domain.BranchResponse, error) {
	branch, err := s.branchRepo.FindByID(branchID)
	if err != nil {
		return nil, err
	}
	resp := mapBranchToResponse(branch)
	return &resp, nil
}

func (s *companyService) GetBranches(companyID int64, params domain.PaginationParams) ([]domain.BranchResponse, int64, error) {
	branches, total, err := s.branchRepo.FindByCompany(companyID, params)
	if err != nil {
		return nil, 0, err
	}
	responses := make([]domain.BranchResponse, len(branches))
	for i := range branches {
		responses[i] = mapBranchToResponse(&branches[i])
	}
	return responses, total, nil
}

func (s *companyService) CreateBranch(companyID int64, req domain.CreateBranchRequest) (*domain.BranchResponse, error) {
	// Verify the parent company exists.
	if _, err := s.companyRepo.FindByID(companyID); err != nil {
		return nil, err
	}
	branch := &domain.Branch{
		CompanyID:   companyID,
		Code:        req.Code,
		Name:        req.Name,
		Address:     req.Address,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		Status:      "active",
	}
	if err := s.branchRepo.Create(branch); err != nil {
		return nil, err
	}
	resp := mapBranchToResponse(branch)
	return &resp, nil
}

func (s *companyService) UpdateBranch(branchID int64, req domain.UpdateBranchRequest) (*domain.BranchResponse, error) {
	branch, err := s.branchRepo.FindByID(branchID)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		branch.Name = *req.Name
	}
	if req.Address != nil {
		branch.Address = req.Address
	}
	if req.PhoneNumber != nil {
		branch.PhoneNumber = req.PhoneNumber
	}
	if req.Email != nil {
		branch.Email = req.Email
	}
	if req.Status != nil {
		branch.Status = *req.Status
	}
	if err := s.branchRepo.Update(branch); err != nil {
		return nil, err
	}
	resp := mapBranchToResponse(branch)
	return &resp, nil
}

func (s *companyService) DeleteBranch(branchID int64) error {
	return s.branchRepo.Delete(branchID)
}

// ─── Company-user operations ──────────────────────────────────────────────────

func (s *companyService) GetCompanyUsers(companyID int64) ([]domain.CompanyUserResponse, error) {
	// Verify company exists.
	if _, err := s.companyRepo.FindByID(companyID); err != nil {
		return nil, err
	}

	users, err := s.companyUserRepo.FindUsersByCompany(companyID)
	if err != nil {
		return nil, err
	}

	// Fetch all role assignments for this company.
	allRoles, err := s.companyUserRepo.FindUserRolesByCompany(companyID)
	if err != nil {
		return nil, err
	}

	// Collect unique role IDs and branch IDs for batch lookup.
	roleIDSet := make(map[int64]struct{})
	branchIDSet := make(map[int64]struct{})
	for _, ur := range allRoles {
		roleIDSet[ur.RoleID] = struct{}{}
		branchIDSet[ur.BranchID] = struct{}{}
	}

	roleIDs := make([]int64, 0, len(roleIDSet))
	for id := range roleIDSet {
		roleIDs = append(roleIDs, id)
	}
	branchIDs := make([]int64, 0, len(branchIDSet))
	for id := range branchIDSet {
		branchIDs = append(branchIDs, id)
	}

	// Batch-fetch role names.
	roleNames := make(map[int64]string)
	if len(roleIDs) > 0 {
		var roles []domain.Role
		if err := s.db.Where("rol_id IN ?", roleIDs).Find(&roles).Error; err == nil {
			for _, r := range roles {
				roleNames[r.ID] = r.Name
			}
		}
	}

	// Batch-fetch branch names.
	branchNames := make(map[int64]string)
	if len(branchIDs) > 0 {
		var branches []domain.Branch
		if err := s.db.Where("bra_id IN ?", branchIDs).Find(&branches).Error; err == nil {
			for _, b := range branches {
				branchNames[b.ID] = b.Name
			}
		}
	}

	// Add fallback names for any IDs not found.
	for id := range roleIDSet {
		if roleNames[id] == "" {
			roleNames[id] = fmt.Sprintf("Rol #%d", id)
		}
	}
	for id := range branchIDSet {
		if branchNames[id] == "" {
			branchNames[id] = fmt.Sprintf("Sucursal #%d", id)
		}
	}

	// Group roles by user.
	userRolesMap := make(map[int64][]domain.CompanyUserRoleResponse)
	for _, ur := range allRoles {
		userRolesMap[ur.UserID] = append(userRolesMap[ur.UserID], domain.CompanyUserRoleResponse{
			RoleID:     ur.RoleID,
			RoleName:   roleNames[ur.RoleID],
			BranchID:   ur.BranchID,
			BranchName: branchNames[ur.BranchID],
		})
	}

	responses := make([]domain.CompanyUserResponse, len(users))
	for i, u := range users {
		responses[i] = domain.CompanyUserResponse{
			ID:       u.ID,
			Username: u.Username,
			Email:    u.Email,
			Status:   u.Status,
			Person:   mapPersonToResponse(&u.Person),
			Roles:    userRolesMap[u.ID],
		}
		if responses[i].Roles == nil {
			responses[i].Roles = []domain.CompanyUserRoleResponse{}
		}
	}
	return responses, nil
}

func (s *companyService) AddUserToCompany(companyID int64, req domain.AddUserToCompanyRequest) error {
	// Verify company exists.
	if _, err := s.companyRepo.FindByID(companyID); err != nil {
		return err
	}
	// Check if already linked.
	linked, err := s.companyUserRepo.IsUserInCompany(companyID, req.UserID)
	if err != nil {
		return err
	}
	if linked {
		return domain.ErrValidation(domain.ErrCodeUserAlreadyInCompany, "user is already assigned to this company")
	}
	return s.companyUserRepo.AddUserToCompany(companyID, req.UserID)
}

func (s *companyService) RemoveUserFromCompany(companyID, userID int64) error {
	return s.companyUserRepo.RemoveUserFromCompany(companyID, userID)
}

// ─── User-role operations ─────────────────────────────────────────────────────

func (s *companyService) AssignUserRole(companyID, userID int64, req domain.AssignUserRoleRequest) error {
	// Verify the user belongs to the company.
	linked, err := s.companyUserRepo.IsUserInCompany(companyID, userID)
	if err != nil {
		return err
	}
	if !linked {
		return domain.ErrValidation(domain.ErrCodeUserNotInCompany, "user does not belong to this company")
	}
	return s.companyUserRepo.AssignUserRole(domain.UserRole{
		UserID:    userID,
		CompanyID: companyID,
		BranchID:  req.BranchID,
		RoleID:    req.RoleID,
	})
}

func (s *companyService) RemoveUserRole(companyID, userID int64, req domain.AssignUserRoleRequest) error {
	// Protect the bootstrap super-admin: no one can touch their roles.
	if s.superAdminUser != "" {
		user, err := s.userRepo.FindByID(userID)
		if err == nil && user.Username == s.superAdminUser {
			return domain.ErrForbidden("no se pueden modificar los roles del super administrador")
		}
	}
	return s.companyUserRepo.RemoveUserRole(domain.UserRole{
		UserID:    userID,
		CompanyID: companyID,
		BranchID:  req.BranchID,
		RoleID:    req.RoleID,
	})
}
