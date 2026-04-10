package service

import (
	"articnexus/backend/internal/domain"
	"articnexus/backend/internal/repository"
)

// ApplicationService handles business logic for applications and their modules.
type ApplicationService interface {
	GetByID(id int64) (*domain.ApplicationResponse, error)
	// GetAll returns all applications. If companyID > 0 only licensed apps are returned.
	GetAll(companyID int64, params domain.PaginationParams) ([]domain.ApplicationResponse, int64, error)
	Create(req domain.CreateApplicationRequest) (*domain.ApplicationResponse, error)
	Update(id int64, req domain.UpdateApplicationRequest) (*domain.ApplicationResponse, error)
	Delete(id int64) error

	// Module operations.
	GetModule(moduleID int64) (*domain.ModuleResponse, error)
	GetModules(appID int64, params domain.PaginationParams) ([]domain.ModuleResponse, int64, error)
	CreateModule(req domain.CreateModuleRequest) (*domain.ModuleResponse, error)
	UpdateModule(moduleID int64, req domain.UpdateModuleRequest) (*domain.ModuleResponse, error)
	DeleteModule(moduleID int64) error
}

type applicationService struct {
	appRepo        repository.ApplicationRepository
	moduleRepo     repository.ModuleRepository
	companyAppRepo repository.CompanyApplicationRepository
}

// NewApplicationService returns an ApplicationService implementation.
func NewApplicationService(
	appRepo repository.ApplicationRepository,
	moduleRepo repository.ModuleRepository,
	companyAppRepo repository.CompanyApplicationRepository,
) ApplicationService {
	return &applicationService{appRepo: appRepo, moduleRepo: moduleRepo, companyAppRepo: companyAppRepo}
}

func (s *applicationService) GetByID(id int64) (*domain.ApplicationResponse, error) {
	app, err := s.appRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	resp := mapApplicationToResponse(app)
	return &resp, nil
}

func (s *applicationService) GetAll(companyID int64, params domain.PaginationParams) ([]domain.ApplicationResponse, int64, error) {
	var apps []domain.Application
	var total int64
	var err error

	if companyID > 0 {
		appIDs, aErr := s.companyAppRepo.GetLicensedAppIDs(companyID)
		if aErr != nil {
			return nil, 0, aErr
		}
		if len(appIDs) == 0 {
			return []domain.ApplicationResponse{}, 0, nil
		}
		apps, total, err = s.appRepo.FindAllByIDs(appIDs, params)
	} else {
		apps, total, err = s.appRepo.FindAll(params)
	}
	if err != nil {
		return nil, 0, err
	}
	responses := make([]domain.ApplicationResponse, len(apps))
	for i := range apps {
		responses[i] = mapApplicationToResponse(&apps[i])
	}
	return responses, total, nil
}

func (s *applicationService) Create(req domain.CreateApplicationRequest) (*domain.ApplicationResponse, error) {
	app := &domain.Application{Code: req.Code, Name: req.Name, Status: "active"}
	if err := s.appRepo.Create(app); err != nil {
		return nil, err
	}
	resp := mapApplicationToResponse(app)
	return &resp, nil
}

func (s *applicationService) Update(id int64, req domain.UpdateApplicationRequest) (*domain.ApplicationResponse, error) {
	app, err := s.appRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		app.Name = *req.Name
	}
	if req.Status != nil {
		app.Status = *req.Status
	}
	if err := s.appRepo.Update(app); err != nil {
		return nil, err
	}
	resp := mapApplicationToResponse(app)
	return &resp, nil
}

func (s *applicationService) Delete(id int64) error {
	return s.appRepo.Delete(id)
}

// ─── Module operations ────────────────────────────────────────────────────────

func (s *applicationService) GetModule(moduleID int64) (*domain.ModuleResponse, error) {
	module, err := s.moduleRepo.FindByID(moduleID)
	if err != nil {
		return nil, err
	}
	resp := mapModuleToResponse(module)
	return &resp, nil
}

func (s *applicationService) GetModules(appID int64, params domain.PaginationParams) ([]domain.ModuleResponse, int64, error) {
	modules, total, err := s.moduleRepo.FindByApplication(appID, params)
	if err != nil {
		return nil, 0, err
	}
	responses := make([]domain.ModuleResponse, len(modules))
	for i := range modules {
		responses[i] = mapModuleToResponse(&modules[i])
	}
	return responses, total, nil
}

func (s *applicationService) CreateModule(req domain.CreateModuleRequest) (*domain.ModuleResponse, error) {
	// Verify the parent application exists.
	if _, err := s.appRepo.FindByID(req.ApplicationID); err != nil {
		return nil, err
	}
	module := &domain.Module{
		ApplicationID: req.ApplicationID,
		Name:          req.Name,
		DisplayName:   req.DisplayName,
		MenuOption:    req.MenuOption,
		SubFunction:   req.SubFunction,
		Description:   req.Description,
		Status:        "active",
	}
	if err := s.moduleRepo.Create(module); err != nil {
		return nil, err
	}
	resp := mapModuleToResponse(module)
	return &resp, nil
}

func (s *applicationService) UpdateModule(moduleID int64, req domain.UpdateModuleRequest) (*domain.ModuleResponse, error) {
	module, err := s.moduleRepo.FindByID(moduleID)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		module.Name = *req.Name
	}
	if req.DisplayName != nil {
		module.DisplayName = req.DisplayName
	}
	if req.MenuOption != nil {
		module.MenuOption = req.MenuOption
	}
	if req.SubFunction != nil {
		module.SubFunction = req.SubFunction
	}
	if req.Description != nil {
		module.Description = req.Description
	}
	if req.Status != nil {
		module.Status = *req.Status
	}
	if err := s.moduleRepo.Update(module); err != nil {
		return nil, err
	}
	resp := mapModuleToResponse(module)
	return &resp, nil
}

func (s *applicationService) DeleteModule(moduleID int64) error {
	return s.moduleRepo.Delete(moduleID)
}
