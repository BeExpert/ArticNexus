package service

import (
	"crypto/rand"
	"fmt"
	"log"

	"articnexus/backend/internal/config"
	"articnexus/backend/internal/domain"
	"articnexus/backend/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService handles business logic for user management.
type UserService interface {
	GetByID(id int64) (*domain.UserResponse, error)
	GetAll(params domain.PaginationParams) ([]domain.UserResponse, int64, error)
	Create(req domain.CreateUserRequest) (*domain.UserResponse, error)
	Update(id int64, req domain.UpdateUserRequest) (*domain.UserResponse, error)
	Delete(id int64) error
	ResetUserPassword(id int64, req domain.AdminResetPasswordRequest) (*domain.AdminResetPasswordResponse, error)
}

type userService struct {
	userRepo        repository.UserRepository
	personRepo      repository.PersonRepository
	companyUserRepo repository.CompanyUserRepository
	db              *gorm.DB
	emailService    EmailService
	cfg             *config.Config
}

// NewUserService returns a UserService implementation.
func NewUserService(
	db *gorm.DB,
	userRepo repository.UserRepository,
	personRepo repository.PersonRepository,
	companyUserRepo repository.CompanyUserRepository,
	emailService EmailService,
	cfg *config.Config,
) UserService {
	return &userService{
		db:              db,
		userRepo:        userRepo,
		personRepo:      personRepo,
		companyUserRepo: companyUserRepo,
		emailService:    emailService,
		cfg:             cfg,
	}
}

func (s *userService) GetByID(id int64) (*domain.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	companies, _ := s.userRepo.FindCompaniesByUserID(user.ID)
	resp := mapUserToResponse(user)
	for i := range companies {
		resp.Companies = append(resp.Companies, mapCompanyToResponse(&companies[i]))
	}
	return &resp, nil
}

func (s *userService) GetAll(params domain.PaginationParams) ([]domain.UserResponse, int64, error) {
	users, total, err := s.userRepo.FindAll(params)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]domain.UserResponse, len(users))
	for i := range users {
		responses[i] = mapUserToResponse(&users[i])
		companies, _ := s.userRepo.FindCompaniesByUserID(users[i].ID)
		for j := range companies {
			responses[i].Companies = append(responses[i].Companies, mapCompanyToResponse(&companies[j]))
		}
	}
	return responses, total, nil
}

// Create creates a Person and a User in a single database transaction.
func (s *userService) Create(req domain.CreateUserRequest) (*domain.UserResponse, error) {
	var plainPassword string
	var err error

	if req.SendCredentials {
		plainPassword, err = generateRandomPassword(12)
		if err != nil {
			return nil, fmt.Errorf("could not generate random password: %w", err)
		}
	} else {
		if req.Password == "" {
			return nil, fmt.Errorf("password is required when sendCredentials is false")
		}
		plainPassword = req.Password
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("could not hash password: %w", err)
	}

	var created *domain.User

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		personRepo := repository.NewPersonRepository(tx)
		userRepo := repository.NewUserRepository(tx)

		person := &domain.Person{
			FirstName:      req.FirstName,
			FirstSurname:   req.FirstSurname,
			SecondSurname:  req.SecondSurname,
			NationalID:     req.NationalID,
			BirthDate:      parseDate(req.BirthDate),
			PhoneAreaCode:  req.PhoneAreaCode,
			PrimaryPhone:   req.PrimaryPhone,
			SecondaryPhone: req.SecondaryPhone,
			Address:        req.Address,
			Status:         "active",
		}

		if err := personRepo.Create(person); err != nil {
			return fmt.Errorf("could not create person: %w", err)
		}

		user := &domain.User{
			PersonID: person.ID,
			Username: req.Username,
			Email:    req.Email,
			Password: string(hash),
			Status:   "active",
		}

		if err := userRepo.Create(user); err != nil {
			return fmt.Errorf("could not create user: %w", err)
		}

		// Optionally associate to a company.
		if req.CompanyID != nil {
			companyUserRepo := repository.NewCompanyUserRepository(tx)
			if err := companyUserRepo.AddUserToCompany(*req.CompanyID, user.ID); err != nil {
				return fmt.Errorf("could not add user to company: %w", err)
			}
		}

		user.Person = *person
		created = user
		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	// Send welcome credentials email when requested. A failure here is
	// non-blocking — the user account has already been created successfully.
	if req.SendCredentials && created.Email != "" {
		loginURL := s.cfg.FrontendURL + "/login"
		if sendErr := s.emailService.SendWelcomeCredentials(
			created.Email,
			created.Username,
			plainPassword,
			loginURL,
			s.cfg.SupportSMTPFrom,
		); sendErr != nil {
			log.Printf("[WARN] could not send welcome credentials to %s: %v", created.Email, sendErr)
		}
	}

	resp := mapUserToResponse(created)
	return &resp, nil
}

func (s *userService) Update(id int64, req domain.UpdateUserRequest) (*domain.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// ── User fields ─────────────────────────────────────────────────────
	if req.Username != nil {
		user.Username = *req.Username
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("could not hash password: %w", err)
		}
		user.Password = string(hash)
	}
	if req.Status != nil {
		user.Status = *req.Status
	}
	if req.ClearExpiry {
		user.PasswordExpiresAt = nil
	} else if req.PasswordExpiresAt != nil {
		user.PasswordExpiresAt = req.PasswordExpiresAt
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	// ── Person fields ───────────────────────────────────────────────────
	personDirty := false
	p := &user.Person
	if req.FirstName != nil {
		p.FirstName = *req.FirstName
		personDirty = true
	}
	if req.FirstSurname != nil {
		p.FirstSurname = *req.FirstSurname
		personDirty = true
	}
	if req.SecondSurname != nil {
		p.SecondSurname = req.SecondSurname
		personDirty = true
	}
	if req.NationalID != nil {
		p.NationalID = req.NationalID
		personDirty = true
	}
	if req.BirthDate != nil {
		p.BirthDate = parseDate(req.BirthDate)
		personDirty = true
	}
	if req.PhoneAreaCode != nil {
		p.PhoneAreaCode = req.PhoneAreaCode
		personDirty = true
	}
	if req.PrimaryPhone != nil {
		p.PrimaryPhone = req.PrimaryPhone
		personDirty = true
	}
	if req.SecondaryPhone != nil {
		p.SecondaryPhone = req.SecondaryPhone
		personDirty = true
	}
	if req.Address != nil {
		p.Address = req.Address
		personDirty = true
	}

	if personDirty {
		if err := s.personRepo.Update(p); err != nil {
			return nil, fmt.Errorf("could not update person: %w", err)
		}
	}

	resp := mapUserToResponse(user)
	return &resp, nil
}

func (s *userService) Delete(id int64) error {
	return s.userRepo.Delete(id)
}

// ResetUserPassword allows an admin to forcibly reset another user's password.
// If req.GenerateRandom is true, a cryptographically random 12-char password is
// generated, set, and returned in plain text (shown once, never stored).
// Otherwise req.Password is used.
func (s *userService) ResetUserPassword(id int64, req domain.AdminResetPasswordRequest) (*domain.AdminResetPasswordResponse, error) {
	if !req.GenerateRandom && req.Password == nil {
		return nil, fmt.Errorf("either generateRandom must be true or a password must be provided")
	}

	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	var plainPassword string
	if req.GenerateRandom {
		plainPassword, err = generateRandomPassword(12)
		if err != nil {
			return nil, fmt.Errorf("could not generate random password: %w", err)
		}
	} else {
		plainPassword = *req.Password
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("could not hash password: %w", err)
	}
	user.Password = string(hash)

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	resp := &domain.AdminResetPasswordResponse{}
	if req.GenerateRandom {
		resp.GeneratedPassword = plainPassword
	}
	return resp, nil
}

// generateRandomPassword returns a cryptographically random password of the
// given length using alphanumeric characters.
func generateRandomPassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b), nil
}
