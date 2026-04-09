package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"articnexus/backend/internal/config"
	"articnexus/backend/internal/domain"
	"articnexus/backend/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// DemoLinkService manages public demo links with JWT tokens.
type DemoLinkService interface {
	Create(actorID int64, req domain.CreateDemoLinkRequest) (*domain.DemoLinkResponse, error)
	List() ([]domain.DemoLinkResponse, error)
	Revoke(id int64) error
}

type demoLinkService struct {
	repo      repository.DemoLinkRepository
	jwtSecret string
	emailSvc  EmailService
	cfg       *config.Config
}

func NewDemoLinkService(
	repo repository.DemoLinkRepository,
	jwtSecret string,
	emailSvc EmailService,
	cfg *config.Config,
) DemoLinkService {
	return &demoLinkService{
		repo:      repo,
		jwtSecret: jwtSecret,
		emailSvc:  emailSvc,
		cfg:       cfg,
	}
}

func (s *demoLinkService) Create(actorID int64, req domain.CreateDemoLinkRequest) (*domain.DemoLinkResponse, error) {
	if req.ExpiresInHours <= 0 {
		req.ExpiresInHours = 24
	}

	// ── Auto-resolve demo user from app code if not provided ──────────────────
	if req.DemoUserID == 0 {
		id, err := s.repo.FindDefaultDemoUserByApp(req.AppCode)
		if err != nil {
			return nil, fmt.Errorf("resolver usuario demo para %s: %w", req.AppCode, err)
		}
		req.DemoUserID = id
	}

	expiresAt := time.Now().Add(time.Duration(req.ExpiresInHours) * time.Hour)

	// ── Rotate demo user password ─────────────────────────────────────────────
	// A fresh random password is generated and stored on every demo link creation
	// so that each invitation invalidates the previous one's direct-login credentials.
	passBytes := make([]byte, 16)
	if _, err := rand.Read(passBytes); err != nil {
		return nil, fmt.Errorf("generate demo password: %w", err)
	}
	tempPass := hex.EncodeToString(passBytes)

	cost := s.cfg.BcryptCost
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(tempPass), cost)
	if err != nil {
		return nil, fmt.Errorf("hash demo password: %w", err)
	}
	if err := s.repo.UpdateDemoUserPassword(req.DemoUserID, string(hash)); err != nil {
		return nil, fmt.Errorf("rotate demo user password: %w", err)
	}

	// ── Build JWT with custom demo_link claims ────────────────────────────────
	claims := jwt.MapClaims{
		"sub":  req.DemoUserID,
		"app":  req.AppCode,
		"type": "demo_link",
		"iat":  time.Now().Unix(),
		"exp":  expiresAt.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, fmt.Errorf("sign demo link token: %w", err)
	}

	// Store SHA-256 hash — never store the raw JWT.
	h := sha256.Sum256([]byte(tokenStr))
	tokenHash := fmt.Sprintf("%x", h)

	link := &domain.DemoLink{
		TokenHash:      tokenHash,
		AppCode:        req.AppCode,
		DemoUserID:     req.DemoUserID,
		ExpiresAt:      expiresAt,
		IsActive:       true,
		RecipientEmail: req.RecipientEmail,
		CreatedBy:      actorID,
	}
	if err := s.repo.Create(link); err != nil {
		return nil, err
	}

	// ── Send invitation email (best-effort) ───────────────────────────────────
	if req.RecipientEmail != nil && *req.RecipientEmail != "" {
		demoUsername := "demo_" + strings.ToLower(req.AppCode)
		appName := req.AppCode
		demoURL := fmt.Sprintf("%s/demo/enter?token=%s", s.cfg.AppBaseURL(req.AppCode), tokenStr)
		guestName := "Estimado usuario"
		if req.GuestName != nil && *req.GuestName != "" {
			guestName = *req.GuestName
		}
		_ = s.emailSvc.SendDemoInvitation(guestName, *req.RecipientEmail, appName, demoURL, tempPass, demoUsername)
	}

	resp := mapDemoLinkToResponse(*link)
	resp.Token = &tokenStr
	return resp, nil
}

func (s *demoLinkService) List() ([]domain.DemoLinkResponse, error) {
	links, err := s.repo.ListAll()
	if err != nil {
		return nil, err
	}
	out := make([]domain.DemoLinkResponse, len(links))
	for i, l := range links {
		out[i] = *mapDemoLinkToResponse(l)
	}
	return out, nil
}

func (s *demoLinkService) Revoke(id int64) error {
	if err := s.repo.Revoke(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return fmt.Errorf("revoke demo link: %w", err)
	}
	return nil
}

func mapDemoLinkToResponse(l domain.DemoLink) *domain.DemoLinkResponse {
	return &domain.DemoLinkResponse{
		ID:             l.ID,
		AppCode:        l.AppCode,
		DemoUserID:     l.DemoUserID,
		ExpiresAt:      l.ExpiresAt,
		IsActive:       l.IsActive,
		RecipientEmail: l.RecipientEmail,
		CreatedBy:      l.CreatedBy,
		CreatedAt:      l.CreatedAt,
	}
}
