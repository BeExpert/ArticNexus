package domain

import "time"

// RoleModule maps the junction table tblRoleModules_RMO.
type RoleModule struct {
	RoleID    int64     `gorm:"column:rol_id;primaryKey"`
	ModuleID  int64     `gorm:"column:mod_id;primaryKey"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (RoleModule) TableName() string { return "tblRoleModules_RMO" }

// UserCompany maps the junction table tblUserCompanies_UCO.
type UserCompany struct {
	UserID    int64     `gorm:"column:usr_id;primaryKey"`
	CompanyID int64     `gorm:"column:com_id;primaryKey"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (UserCompany) TableName() string { return "tblUserCompanies_UCO" }

// UserBranch maps the junction table tblUserBranches_UBR.
type UserBranch struct {
	UserID    int64     `gorm:"column:usr_id;primaryKey"`
	BranchID  int64     `gorm:"column:bra_id;primaryKey"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (UserBranch) TableName() string { return "tblUserBranches_UBR" }

// UserRole maps the junction table tblUserRoles_URO.
type UserRole struct {
	UserID    int64     `gorm:"column:usr_id;primaryKey"`
	CompanyID int64     `gorm:"column:com_id;primaryKey"`
	BranchID  int64     `gorm:"column:bra_id;primaryKey"`
	RoleID    int64     `gorm:"column:rol_id;primaryKey"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (UserRole) TableName() string { return "tblUserRoles_URO" }

// ─── DTOs for company-user management ────────────────────────────────────────

// AddUserToCompanyRequest links an existing user to a company.
type AddUserToCompanyRequest struct {
	UserID int64 `json:"userId" validate:"required"`
}

// AssignUserRoleRequest assigns a role to a user within a company + branch.
type AssignUserRoleRequest struct {
	BranchID int64 `json:"branchId" validate:"required"`
	RoleID   int64 `json:"roleId"   validate:"required"`
}

// CompanyUserResponse represents a user within a company context.
type CompanyUserResponse struct {
	ID       int64          `json:"id"`
	Username string         `json:"username"`
	Email    string         `json:"email"`
	Status   string         `json:"status"`
	Person   PersonResponse `json:"person"`
	Roles    []CompanyUserRoleResponse `json:"roles"`
}

// CompanyUserRoleResponse represents a role assignment for a user in a company.
type CompanyUserRoleResponse struct {
	RoleID     int64  `json:"roleId"`
	RoleName   string `json:"roleName"`
	BranchID   int64  `json:"branchId"`
	BranchName string `json:"branchName"`
}
