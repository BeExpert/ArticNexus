package service

import (
	"fmt"

	"articnexus/backend/internal/domain"
	"articnexus/backend/internal/repository"
)

// RoleService handles business logic for roles and their module assignments.
type RoleService interface {
	GetByID(id int64) (*domain.RoleResponse, error)
	// GetAll returns all roles. If companyID > 0 only roles for licensed apps are returned.
	GetAll(companyID int64, params domain.PaginationParams) ([]domain.RoleResponse, int64, error)
	Create(req domain.CreateRoleRequest) (*domain.RoleResponse, error)
	Update(id int64, req domain.UpdateRoleRequest) (*domain.RoleResponse, error)
	Delete(id int64) error

	// Module assignment operations.
	GetModules(roleID int64) ([]domain.ModuleResponse, error)
	AssignModules(roleID int64, req domain.AssignModulesRequest) error
	RemoveModules(roleID int64, req domain.AssignModulesRequest) error
}

type roleService struct {
	roleRepo       repository.RoleRepository
	appRepo        repository.ApplicationRepository
	moduleRepo     repository.ModuleRepository
	companyAppRepo repository.CompanyApplicationRepository
}

// NewRoleService returns a RoleService implementation.
func NewRoleService(
	roleRepo repository.RoleRepository,
	appRepo repository.ApplicationRepository,
	moduleRepo repository.ModuleRepository,
	companyAppRepo repository.CompanyApplicationRepository,
) RoleService {
	return &roleService{
		roleRepo:       roleRepo,
		appRepo:        appRepo,
		moduleRepo:     moduleRepo,
		companyAppRepo: companyAppRepo,
	}
}

func (s *roleService) GetByID(id int64) (*domain.RoleResponse, error) {
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	resp := mapRoleToResponse(role)
	return &resp, nil
}

func (s *roleService) GetAll(companyID int64, params domain.PaginationParams) ([]domain.RoleResponse, int64, error) {
	var roles []domain.Role
	var total int64
	var err error

	if companyID > 0 {
		appIDs, aErr := s.companyAppRepo.GetLicensedAppIDs(companyID)
		if aErr != nil {
			return nil, 0, aErr
		}
		if len(appIDs) == 0 {
			return []domain.RoleResponse{}, 0, nil
		}
		roles, total, err = s.roleRepo.FindAllByAppIDs(appIDs, params)
	} else {
		roles, total, err = s.roleRepo.FindAll(params)
	}
	if err != nil {
		return nil, 0, err
	}
	responses := make([]domain.RoleResponse, len(roles))
	for i := range roles {
		responses[i] = mapRoleToResponse(&roles[i])
	}
	return responses, total, nil
}

func (s *roleService) Create(req domain.CreateRoleRequest) (*domain.RoleResponse, error) {
	// Verify the parent application exists.
	if _, err := s.appRepo.FindByID(req.ApplicationID); err != nil {
		return nil, err
	}
	// If a company scope is provided, verify the company is licensed for this app.
	if req.CompanyID > 0 {
		ok, err := s.companyAppRepo.IsLicensed(req.CompanyID, req.ApplicationID)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, domain.ErrForbidden("la empresa no tiene licencia para esta aplicación")
		}
	}
	role := &domain.Role{
		ApplicationID: req.ApplicationID,
		Name:          req.Name,
		Status:        "active",
	}
	if err := s.roleRepo.Create(role); err != nil {
		return nil, err
	}
	resp := mapRoleToResponse(role)
	return &resp, nil
}

func (s *roleService) Update(id int64, req domain.UpdateRoleRequest) (*domain.RoleResponse, error) {
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		role.Name = *req.Name
	}
	if req.Status != nil {
		role.Status = *req.Status
	}
	if err := s.roleRepo.Update(role); err != nil {
		return nil, err
	}
	resp := mapRoleToResponse(role)
	return &resp, nil
}

func (s *roleService) Delete(id int64) error {
	return s.roleRepo.Delete(id)
}

func (s *roleService) GetModules(roleID int64) ([]domain.ModuleResponse, error) {
	modules, err := s.roleRepo.FindModules(roleID)
	if err != nil {
		return nil, err
	}
	responses := make([]domain.ModuleResponse, len(modules))
	for i := range modules {
		responses[i] = mapModuleToResponse(&modules[i])
	}
	return responses, nil
}

func (s *roleService) AssignModules(roleID int64, req domain.AssignModulesRequest) error {
	// Verify the role exists and get its application scope.
	role, err := s.roleRepo.FindByID(roleID)
	if err != nil {
		return err
	}

	// Security: ensure all requested modules belong to the same app as the role.
	if len(req.ModuleIDs) > 0 {
		mods, err := s.moduleRepo.FindByIDs(req.ModuleIDs)
		if err != nil {
			return err
		}
		for _, m := range mods {
			if m.ApplicationID != role.ApplicationID {
				return domain.ErrValidation(
					domain.ErrCodeValidation,
					fmt.Sprintf("el módulo '%s' no pertenece a la misma aplicación que el rol", m.Name),
				)
			}
		}
	}
	return s.roleRepo.AssignModules(roleID, req.ModuleIDs)
}

func (s *roleService) RemoveModules(roleID int64, req domain.AssignModulesRequest) error {
	return s.roleRepo.RemoveModules(roleID, req.ModuleIDs)
}
