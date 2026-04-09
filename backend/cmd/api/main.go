package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"articnexus/backend/internal/config"
	"articnexus/backend/internal/db"
	"articnexus/backend/internal/handler"
	"articnexus/backend/internal/repository"
	"articnexus/backend/internal/router"
	"articnexus/backend/internal/service"
	"articnexus/backend/pkg/logger"
)

func main() {
	// ── CLI flags ────────────────────────────────────────────────────────────
	migrate := flag.Bool("migrate", false, "run database migrations before starting the server")
	flag.Parse()

	// ── Configuration ────────────────────────────────────────────────────────
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	// ── Logging ──────────────────────────────────────────────────────────────
	if err := logger.Init(cfg.MaxLogSizeMB, cfg.MaxLogBackups); err != nil {
		log.Printf("warning: could not init file logger: %v — falling back to stderr only", err)
	}
	logger.Info(logger.App, fmt.Sprintf("ArticNexus starting (env=%s port=%s)", cfg.Environment, cfg.Port))

	// ── Session epoch ────────────────────────────────────────────────────────
	// A fresh 16-byte random value is generated on every boot.  Any JWT that
	// carries a different epoch is rejected, forcing re-login after restarts.
	epochBytes := make([]byte, 16)
	if _, err := rand.Read(epochBytes); err != nil {
		log.Fatalf("failed to generate session epoch: %v", err)
	}
	cfg.SessionEpoch = hex.EncodeToString(epochBytes)
	logger.Info(logger.App, "session epoch generated — all previous sessions invalidated")

	// ── Database ─────────────────────────────────────────────────────────────
	database, err := db.New(cfg.DatabaseURL, cfg.IsDevelopment(), db.PoolConfig{
		MaxOpenConns:    cfg.DBMaxOpenConns,
		MaxIdleConns:    cfg.DBMaxIdleConns,
		ConnMaxLifetime: cfg.DBConnMaxLifetime,
	})
	if err != nil {
		logger.Error(logger.DB, fmt.Sprintf("database connection error: %v", err))
		log.Fatalf("database connection error: %v", err)
	}
	logger.Info(logger.DB, fmt.Sprintf("database connected (url configured, env=%s)", cfg.Environment))

	sqlDB, _ := database.DB()
	defer sqlDB.Close()

	// ── Optional migrations ──────────────────────────────────────────────────
	if *migrate {
		logger.Info(logger.App, fmt.Sprintf("running migrations from %s", cfg.MigrationsDir))
		if err := db.Migrate(database, cfg.MigrationsDir); err != nil {
			logger.Error(logger.App, fmt.Sprintf("migration error: %v", err))
			log.Fatalf("migration error: %v", err)
		}
		logger.Info(logger.App, "migrations completed successfully")
	}

	// ── Seed super-admin ─────────────────────────────────────────────────────
	if err := db.SeedSuperAdmin(database, cfg); err != nil {
		logger.Error(logger.App, fmt.Sprintf("seed super-admin error: %v", err))
		log.Fatalf("seed error: %v", err)
	}

	// ── Seed ArticNexus modules (upsert on every boot) ───────────────────────
	if err := db.SeedModules(database); err != nil {
		logger.Error(logger.App, fmt.Sprintf("seed modules error: %v", err))
		log.Fatalf("seed modules error: %v", err)
	}

	// ── Seed all applications and their modules (idempotent) ─────────────────
	if err := db.SeedApplications(database); err != nil {
		logger.Error(logger.App, fmt.Sprintf("seed applications error: %v", err))
		log.Fatalf("seed applications error: %v", err)
	}
	// ── Seed ArticDev company, branch, demo users and DEMO roles (idempotent) ───
	if err := db.SeedArticDevAndDemoUsers(database, cfg); err != nil {
		logger.Error(logger.App, fmt.Sprintf("seed articdev error: %v", err))
		log.Fatalf("seed articdev error: %v", err)
	}
	logger.Info(logger.App, "all seeds completed successfully")

	// ── Repositories ─────────────────────────────────────────────────────────
	personRepo := repository.NewPersonRepository(database)
	userRepo := repository.NewUserRepository(database)
	companyRepo := repository.NewCompanyRepository(database)
	branchRepo := repository.NewBranchRepository(database)
	companyUserRepo := repository.NewCompanyUserRepository(database)
	appRepo := repository.NewApplicationRepository(database)
	moduleRepo := repository.NewModuleRepository(database)
	roleRepo := repository.NewRoleRepository(database)
	resetRepo := repository.NewPasswordResetRepository(database)
	demoLinkRepo := repository.NewDemoLinkRepository(database)

	// ── Services ─────────────────────────────────────────────────────────────
	emailSvc := service.NewEmailService(cfg)
	authSvc := service.NewAuthService(userRepo, personRepo, moduleRepo, resetRepo, emailSvc, cfg, cfg.JWTSecret, cfg.SuperAdminUser)
	userSvc := service.NewUserService(database, userRepo, personRepo, companyUserRepo)
	companySvc := service.NewCompanyService(database, companyRepo, branchRepo, companyUserRepo, roleRepo, personRepo, userRepo)
	appSvc := service.NewApplicationService(appRepo, moduleRepo)
	roleSvc := service.NewRoleService(roleRepo, appRepo)
	demoLinkSvc := service.NewDemoLinkService(demoLinkRepo, cfg.JWTSecret, emailSvc, cfg)

	// ── Handlers ─────────────────────────────────────────────────────────────
	authHandler := handler.NewAuthHandler(authSvc, companySvc, userSvc)
	userHandler := handler.NewUserHandler(userSvc)
	companyHandler := handler.NewCompanyHandler(companySvc)
	appHandler := handler.NewApplicationHandler(appSvc)
	roleHandler := handler.NewRoleHandler(roleSvc)
	statsHandler := handler.NewStatsHandler(database)
	contactHandler := handler.NewContactHandler(emailSvc)
	demoLinkHandler := handler.NewDemoLinkHandler(demoLinkSvc)

	// ── Router ───────────────────────────────────────────────────────────────
	r := router.New(
		cfg.JWTSecret,
		cfg.SessionEpoch,
		cfg.SuperAdminUser,
		cfg.AllowedOrigins,
		userRepo,
		moduleRepo,
		authHandler,
		userHandler,
		companyHandler,
		appHandler,
		roleHandler,
		statsHandler,
		contactHandler,
		demoLinkHandler,
	)

	// ── HTTP server ──────────────────────────────────────────────────────────
	addr := fmt.Sprintf(":%s", cfg.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start listening in a goroutine so the main goroutine can handle signals.
	serverErr := make(chan error, 1)
	go func() {
		logger.Info(logger.App, fmt.Sprintf("server listening on %s", addr))
		serverErr <- srv.ListenAndServe()
	}()

	// ── Graceful shutdown ────────────────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	// SIGHUP is sent when the controlling terminal closes (e.g. VS Code tab closed).
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	select {
	case err := <-serverErr:
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Error(logger.App, fmt.Sprintf("server error: %v", err))
			log.Fatalf("server error: %v", err)
		}
	case sig := <-quit:
		logger.Info(logger.App, fmt.Sprintf("received signal %s, shutting down gracefully", sig))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error(logger.App, fmt.Sprintf("graceful shutdown failed: %v", err))
		log.Fatalf("graceful shutdown failed: %v", err)
	}

	logger.Info(logger.App, "server stopped")
}
