package service

import (
	"articnexus/backend/internal/domain"
	"articnexus/backend/internal/repository"
)

// RoleService handles business logic for roles and their module assignments.
type RoleService interface {
	GetByID(id int64) (*domain.RoleResponse, error)
	GetAll(params domain.PaginationParams) ([]domain.RoleResponse, int64, error)
	Create(req domain.CreateRoleRequest) (*domain.RoleResponse, error)
	Update(id int64, req domain.UpdateRoleRequest) (*domain.RoleResponse, error)
	Delete(id int64) error

	// Module assignment operations.
	GetModules(roleID int64) ([]domain.ModuleResponse, error)
	AssignModules(roleID int64, req domain.AssignModulesRequest) error
	RemoveModules(roleID int64, req domain.AssignModulesRequest) error
}

type roleService struct {
	roleRepo repository.RoleRepository
	appRepo  repository.ApplicationRepository
}

// NewRoleService returns a RoleService implementation.
func NewRoleService(roleRepo repository.RoleRepository, appRepo repository.ApplicationRepository) RoleService {
	return &roleService{roleRepo: roleRepo, appRepo: appRepo}
}

func (s *roleService) GetByID(id int64) (*domain.RoleResponse, error) {
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	resp := mapRoleToResponse(role)
	return &resp, nil
}

func (s *roleService) GetAll(params domain.PaginationParams) ([]domain.RoleResponse, int64, error) {
	roles, total, err := s.roleRepo.FindAll(params)
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
	// Verify the role exists.
	if _, err := s.roleRepo.FindByID(roleID); err != nil {
		return err
	}
	return s.roleRepo.AssignModules(roleID, req.ModuleIDs)
}

func (s *roleService) RemoveModules(roleID int64, req domain.AssignModulesRequest) error {
	return s.roleRepo.RemoveModules(roleID, req.ModuleIDs)
}
