package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"articnexus/backend/internal/config"
	"articnexus/backend/internal/domain"
	"articnexus/backend/internal/repository"
	"articnexus/backend/pkg/logger"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService handles authentication and token management.
type AuthService interface {
	Login(req domain.LoginRequest) (*domain.LoginResponse, error)
	GetCurrentUser(userID int64) (*domain.UserResponse, error)
	UpdateProfile(userID int64, req domain.UpdateProfileRequest) (*domain.UserResponse, error)
	ForgotPassword(username string) error
	ResetPassword(token, newPassword string) error
}

type authService struct {
	userRepo       repository.UserRepository
	personRepo     repository.PersonRepository
	moduleRepo     repository.ModuleRepository
	resetRepo      repository.PasswordResetRepository
	emailSvc       EmailService
	cfg            *config.Config
	jwtSecret      string
	superAdminUser string
}

// NewAuthService returns an AuthService implementation.
func NewAuthService(
	userRepo repository.UserRepository,
	personRepo repository.PersonRepository,
	moduleRepo repository.ModuleRepository,
	resetRepo repository.PasswordResetRepository,
	emailSvc EmailService,
	cfg *config.Config,
	jwtSecret, superAdminUser string,
) AuthService {
	return &authService{
		userRepo:       userRepo,
		personRepo:     personRepo,
		moduleRepo:     moduleRepo,
		resetRepo:      resetRepo,
		emailSvc:       emailSvc,
		cfg:            cfg,
		jwtSecret:      jwtSecret,
		superAdminUser: superAdminUser,
	}
}

func (s *authService) Login(req domain.LoginRequest) (*domain.LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("credenciales incorrectas")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logger.Warn(logger.Security, fmt.Sprintf("login failed (wrong password): %s", req.Username))
		return nil, fmt.Errorf("credenciales incorrectas")
	}

	if user.Status != "active" {
		logger.Warn(logger.Security, fmt.Sprintf("login failed (user inactive): %s", req.Username))
		return nil, fmt.Errorf("tu cuenta está inactiva. Contacta al administrador para rehabilitarla")
	}

	if user.PasswordExpiresAt != nil && time.Now().After(*user.PasswordExpiresAt) {
		logger.Warn(logger.Security, fmt.Sprintf("login denied (demo session expired): %s", req.Username))
		return nil, fmt.Errorf("tu sesión de prueba ha expirado. Contacta al administrador para recuperar el acceso")
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("could not generate token: %w", err)
	}

	resp := mapUserToResponse(user)
	resp.IsSuperAdmin = s.superAdminUser != "" && user.Username == s.superAdminUser
	resp.Permissions = s.loadPermissions(user.ID, resp.IsSuperAdmin)
	logger.Info(logger.Security, fmt.Sprintf("login ok: %s (superAdmin=%v modules=%d)", req.Username, resp.IsSuperAdmin, len(resp.Permissions)))
	return &domain.LoginResponse{
		Token: token,
		User:  resp,
	}, nil
}

func (s *authService) GetCurrentUser(userID int64) (*domain.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	resp := mapUserToResponse(user)
	resp.IsSuperAdmin = s.superAdminUser != "" && user.Username == s.superAdminUser
	resp.Permissions = s.loadPermissions(user.ID, resp.IsSuperAdmin)
	return &resp, nil
}

// loadPermissions returns the ArticNexus platform module names the user
// can access.  Super-admins get ALL ArticNexus modules automatically.
func (s *authService) loadPermissions(userID int64, isSuperAdmin bool) []string {
	const appCode = "ARTICNEXUS"
	if isSuperAdmin {
		// Super-admin gets every ArticNexus module.
		perms, err := s.moduleRepo.FindNamesByAppCode(appCode)
		if err != nil || perms == nil {
			return []string{}
		}
		return perms
	}
	perms, err := s.userRepo.FindUserPermissions(userID, appCode)
	if err != nil || perms == nil {
		return []string{}
	}
	return perms
}

func (s *authService) generateToken(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"sub":   userID,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Duration(s.cfg.JWTExpHours) * time.Hour).Unix(),
		"epoch": s.cfg.SessionEpoch,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// UpdateProfile updates the authenticated user's own profile (no password/status).
func (s *authService) UpdateProfile(userID int64, req domain.UpdateProfileRequest) (*domain.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if req.Email != nil {
		user.Email = *req.Email
	}
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	// Person fields.
	p := &user.Person
	dirty := false
	if req.FirstName != nil {
		p.FirstName = *req.FirstName
		dirty = true
	}
	if req.FirstSurname != nil {
		p.FirstSurname = *req.FirstSurname
		dirty = true
	}
	if req.SecondSurname != nil {
		p.SecondSurname = req.SecondSurname
		dirty = true
	}
	if req.NationalID != nil {
		p.NationalID = req.NationalID
		dirty = true
	}
	if req.BirthDate != nil {
		p.BirthDate = parseDate(req.BirthDate)
		dirty = true
	}
	if req.PhoneAreaCode != nil {
		p.PhoneAreaCode = req.PhoneAreaCode
		dirty = true
	}
	if req.PrimaryPhone != nil {
		p.PrimaryPhone = req.PrimaryPhone
		dirty = true
	}
	if req.SecondaryPhone != nil {
		p.SecondaryPhone = req.SecondaryPhone
		dirty = true
	}
	if req.Address != nil {
		p.Address = req.Address
		dirty = true
	}

	if dirty {
		if err := s.personRepo.Update(p); err != nil {
			return nil, fmt.Errorf("could not update person: %w", err)
		}
	}

	resp := mapUserToResponse(user)
	resp.IsSuperAdmin = s.superAdminUser != "" && user.Username == s.superAdminUser
	resp.Permissions = s.loadPermissions(user.ID, resp.IsSuperAdmin)
	return &resp, nil
}

// ForgotPassword generates a reset token and emails it to the user.
// Always returns nil to avoid leaking whether the username exists.
func (s *authService) ForgotPassword(username string) error {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		// Silently succeed — don't reveal if the username exists.
		return nil
	}

	if user.Email == "" {
		// User has no email — can't send reset.
		return nil
	}

	// Invalidate any previous tokens.
	_ = s.resetRepo.InvalidateAllForUser(user.ID)

	// Generate a cryptographically random token.
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return fmt.Errorf("could not generate token: %w", err)
	}
	plainToken := hex.EncodeToString(raw)

	// Store the SHA-256 hash — the plain token only travels via email.
	hash := sha256.Sum256([]byte(plainToken))
	tokenHash := hex.EncodeToString(hash[:])

	resetToken := &domain.PasswordResetToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(time.Duration(s.cfg.PasswordResetExpMin) * time.Minute),
	}
	if err := s.resetRepo.Create(resetToken); err != nil {
		logger.Warn(logger.Security, fmt.Sprintf("[WARN] could not save reset token for user %d: %v", user.ID, err))
		return nil
	}

	resetURL := fmt.Sprintf("%s/reset-password?token=%s", s.cfg.FrontendURL, plainToken)
	if err := s.emailSvc.SendPasswordReset(user.Email, resetURL); err != nil {
		logger.Warn(logger.Security, fmt.Sprintf("[WARN] could not send reset email to %s: %v", user.Email, err))
	}
	return nil
}

// ResetPassword validates the token and sets a new password.
func (s *authService) ResetPassword(token, newPassword string) error {
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	rt, err := s.resetRepo.FindByTokenHash(tokenHash)
	if err != nil {
		return domain.ErrBadRequest(domain.ErrCodeAuthResetLinkInvalid, "enlace inválido o expirado")
	}
	if rt.Used {
		return domain.ErrBadRequest(domain.ErrCodeAuthResetLinkUsed, "este enlace ya fue utilizado")
	}
	if time.Now().After(rt.ExpiresAt) {
		return domain.ErrBadRequest(domain.ErrCodeAuthResetLinkExpired, "el enlace ha expirado")
	}

	// Hash new password.
	bcryptHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), s.cfg.BcryptCost)
	if err != nil {
		return fmt.Errorf("could not hash password: %w", err)
	}

	user, err := s.userRepo.FindByID(rt.UserID)
	if err != nil {
		return domain.ErrBadRequest(domain.ErrCodeAuthResetLinkInvalid, "usuario no encontrado")
	}
	user.Password = string(bcryptHash)
	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("could not update password: %w", err)
	}

	// Mark token as used.
	_ = s.resetRepo.MarkUsed(rt.ID)
	return nil
}
