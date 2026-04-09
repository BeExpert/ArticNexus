package service

import (
	"articnexus/backend/internal/domain"
	"time"
)

// mapUserToResponse converts a User model + embedded Person to a UserResponse DTO.
// Note: IsSuperAdmin is NOT set here — it is computed by authService which has
// access to the SUPER_ADMIN_USER env value. Other services leave it false.
func mapUserToResponse(user *domain.User) domain.UserResponse {
	return domain.UserResponse{
		ID:                user.ID,
		Username:          user.Username,
		Email:             user.Email,
		Status:            user.Status,
		PasswordExpiresAt: user.PasswordExpiresAt,
		CreatedAt:         user.CreatedAt,
		UpdatedAt:         user.UpdatedAt,
		Person:            mapPersonToResponse(&user.Person),
		Companies:         []domain.CompanyResponse{},
	}
}

// mapPersonToResponse converts a Person model to a PersonResponse DTO.
func mapPersonToResponse(p *domain.Person) domain.PersonResponse {
	var birthDate *string
	if p.BirthDate != nil {
		s := p.BirthDate.Format("2006-01-02")
		birthDate = &s
	}
	return domain.PersonResponse{
		ID:             p.ID,
		FirstName:      p.FirstName,
		FirstSurname:   p.FirstSurname,
		SecondSurname:  p.SecondSurname,
		NationalID:     p.NationalID,
		Email:          p.Email,
		BirthDate:      birthDate,
		PhoneAreaCode:  p.PhoneAreaCode,
		PrimaryPhone:   p.PrimaryPhone,
		SecondaryPhone: p.SecondaryPhone,
		Address:        p.Address,
		Status:         p.Status,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
}

// parseDate parses a "YYYY-MM-DD" string into a *time.Time, returning nil on empty/nil input.
func parseDate(s *string) *time.Time {
	if s == nil || *s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", *s)
	if err != nil {
		return nil
	}
	return &t
}

// mapCompanyToResponse converts a Company model to a CompanyResponse DTO.
func mapCompanyToResponse(c *domain.Company) domain.CompanyResponse {
	return domain.CompanyResponse{
		ID:        c.ID,
		Name:      c.Name,
		Status:    c.Status,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

// mapBranchToResponse converts a Branch model to a BranchResponse DTO.
func mapBranchToResponse(b *domain.Branch) domain.BranchResponse {
	return domain.BranchResponse{
		ID:          b.ID,
		CompanyID:   b.CompanyID,
		Code:        b.Code,
		Name:        b.Name,
		Address:     b.Address,
		PhoneNumber: b.PhoneNumber,
		Email:       b.Email,
		Status:      b.Status,
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
	}
}

// mapApplicationToResponse converts an Application model to an ApplicationResponse DTO.
func mapApplicationToResponse(a *domain.Application) domain.ApplicationResponse {
	return domain.ApplicationResponse{
		ID:        a.ID,
		Code:      a.Code,
		Name:      a.Name,
		Status:    a.Status,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

// mapModuleToResponse converts a Module model to a ModuleResponse DTO.
func mapModuleToResponse(m *domain.Module) domain.ModuleResponse {
	return domain.ModuleResponse{
		ID:            m.ID,
		ApplicationID: m.ApplicationID,
		Name:          m.Name,
		DisplayName:   m.DisplayName,
		MenuOption:    m.MenuOption,
		SubFunction:   m.SubFunction,
		Description:   m.Description,
		Status:        m.Status,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

// mapRoleToResponse converts a Role model to a RoleResponse DTO.
func mapRoleToResponse(r *domain.Role) domain.RoleResponse {
	return domain.RoleResponse{
		ID:            r.ID,
		ApplicationID: r.ApplicationID,
		AppName:       r.Application.Name,
		Name:          r.Name,
		Status:        r.Status,
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
	}
}
