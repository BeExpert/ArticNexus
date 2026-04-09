package router

import (
	"net/http"

	"articnexus/backend/internal/handler"
	"articnexus/backend/internal/middleware"
	"articnexus/backend/internal/repository"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

// New builds and returns the top-level Chi router with all routes registered.
func New(
	jwtSecret string,
	sessionEpoch string,
	superAdminUser string,
	allowedOrigins []string,
	userRepo repository.UserRepository,
	moduleRepo repository.ModuleRepository,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	companyHandler *handler.CompanyHandler,
	appHandler *handler.ApplicationHandler,
	roleHandler *handler.RoleHandler,
	statsHandler *handler.StatsHandler,
	contactHandler *handler.ContactHandler,
	demoLinkHandler *handler.DemoLinkHandler,
) http.Handler {
	r := chi.NewRouter()

	// ── Global middleware stack ──────────────────────────────────────────────
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.CORS(allowedOrigins))

	// ── Health check (unauthenticated) ───────────────────────────────────────
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// ── API v1 ───────────────────────────────────────────────────────────────
	r.Route("/api/v1", func(r chi.Router) {

		// Public auth routes.
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", authHandler.Login)
			r.Post("/forgot-password", authHandler.ForgotPassword)
			r.Post("/reset-password", authHandler.ResetPassword)
		})

		// Public contact form (no auth required).
		r.Post("/contact", contactHandler.Submit)

		// All routes below require a valid JWT.
		r.Group(func(r chi.Router) {
			r.Use(middleware.Authenticate(jwtSecret, sessionEpoch))
			r.Use(middleware.LoadPermissions(userRepo, moduleRepo, superAdminUser))

			// Auth — protected (self-service, no module required).
			r.Post("/auth/logout", authHandler.Logout)
			r.Get("/auth/me", authHandler.Me)
			r.Put("/auth/me", authHandler.UpdateMe)
			r.Get("/auth/me/companies", authHandler.MyCompanies)

			// Stats (dashboard).
			r.With(middleware.RequireModule("dashboard.ver")).Get("/stats", statsHandler.GetStats)

			// Users — granular per-action guards.
			r.Route("/users", func(r chi.Router) {
				r.With(middleware.RequireModule("usuarios.ver")).Get("/", userHandler.GetAll)
				r.With(middleware.RequireModule("usuarios.crear")).Post("/", userHandler.Create)
				r.Route("/{id}", func(r chi.Router) {
					r.With(middleware.RequireModule("usuarios.ver")).Get("/", userHandler.GetByID)
					r.With(middleware.RequireModule("usuarios.editar")).Put("/", userHandler.Update)
					r.With(middleware.RequireModule("usuarios.eliminar")).Delete("/", userHandler.Delete)
					r.With(middleware.RequireModule("usuarios.reset_contrasena")).Post("/reset-password", userHandler.ResetPassword)
				})
			})

			// Companies — granular per-action guards.
			r.Route("/companies", func(r chi.Router) {
				r.With(middleware.RequireModule("empresas.ver")).Get("/", companyHandler.GetAll)
				r.With(middleware.RequireModule("empresas.crear")).Post("/", companyHandler.Create)
				r.Route("/{id}", func(r chi.Router) {
					r.With(middleware.RequireModule("empresas.ver")).Get("/", companyHandler.GetByID)
					r.With(middleware.RequireModule("empresas.editar")).Put("/", companyHandler.Update)
					r.With(middleware.RequireModule("empresas.eliminar")).Delete("/", companyHandler.Delete)

					// Branches sub-resource.
					r.Route("/branches", func(r chi.Router) {
						r.With(middleware.RequireModule("sucursales.ver")).Get("/", companyHandler.GetBranches)
						r.With(middleware.RequireModule("sucursales.crear")).Post("/", companyHandler.CreateBranch)
						r.Route("/{branchId}", func(r chi.Router) {
							r.With(middleware.RequireModule("sucursales.ver")).Get("/", companyHandler.GetBranch)
							r.With(middleware.RequireModule("sucursales.editar")).Put("/", companyHandler.UpdateBranch)
							r.With(middleware.RequireModule("sucursales.eliminar")).Delete("/", companyHandler.DeleteBranch)
						})
					})

					// Company users sub-resource.
					r.Route("/users", func(r chi.Router) {
						r.With(middleware.RequireModule("personas.ver")).Get("/", companyHandler.GetCompanyUsers)
						r.With(middleware.RequireModule("personas.crear")).Post("/", companyHandler.AddUserToCompany)
						r.Route("/{userId}", func(r chi.Router) {
							r.With(middleware.RequireModule("personas.eliminar")).Delete("/", companyHandler.RemoveUserFromCompany)
							// User role assignment.
							r.With(middleware.RequireModule("roles.asignar_modulos")).Post("/roles", companyHandler.AssignUserRole)
							r.With(middleware.RequireModule("roles.asignar_modulos")).Delete("/roles", companyHandler.RemoveUserRole)
						})
					})
				})
			})

			// Applications + Modules — granular per-action guards.
			r.Route("/applications", func(r chi.Router) {
				r.With(middleware.RequireModule("aplicaciones.ver")).Get("/", appHandler.GetAll)
				r.With(middleware.RequireModule("aplicaciones.crear")).Post("/", appHandler.Create)
				r.Route("/{id}", func(r chi.Router) {
					r.With(middleware.RequireModule("aplicaciones.ver")).Get("/", appHandler.GetByID)
					r.With(middleware.RequireModule("aplicaciones.editar")).Put("/", appHandler.Update)
					r.With(middleware.RequireModule("aplicaciones.eliminar")).Delete("/", appHandler.Delete)

					// Modules sub-resource.
					r.Route("/modules", func(r chi.Router) {
						r.With(middleware.RequireModule("modulos.ver")).Get("/", appHandler.GetModules)
						r.With(middleware.RequireModule("modulos.crear")).Post("/", appHandler.CreateModule)
						r.Route("/{moduleId}", func(r chi.Router) {
							r.With(middleware.RequireModule("modulos.ver")).Get("/", appHandler.GetModule)
							r.With(middleware.RequireModule("modulos.editar")).Put("/", appHandler.UpdateModule)
							r.With(middleware.RequireModule("modulos.eliminar")).Delete("/", appHandler.DeleteModule)
						})
					})
				})
			})

			// Roles + Module assignments — granular per-action guards.
			r.Route("/roles", func(r chi.Router) {
				r.With(middleware.RequireModule("roles.ver")).Get("/", roleHandler.GetAll)
				r.With(middleware.RequireModule("roles.crear")).Post("/", roleHandler.Create)
				r.Route("/{id}", func(r chi.Router) {
					r.With(middleware.RequireModule("roles.ver")).Get("/", roleHandler.GetByID)
					r.With(middleware.RequireModule("roles.editar")).Put("/", roleHandler.Update)
					r.With(middleware.RequireModule("roles.eliminar")).Delete("/", roleHandler.Delete)

					// Module assignment sub-resource.
					r.Route("/modules", func(r chi.Router) {
						r.With(middleware.RequireModule("roles.asignar_modulos")).Get("/", roleHandler.GetModules)
						r.With(middleware.RequireModule("roles.asignar_modulos")).Post("/", roleHandler.AssignModules)
						r.With(middleware.RequireModule("roles.asignar_modulos")).Delete("/", roleHandler.RemoveModules)
					})
				})
			})

			// Demo Links.
			r.Route("/demo-links", func(r chi.Router) {
				r.With(middleware.RequireModule("demo_links.ver")).Get("/", demoLinkHandler.List)
				r.With(middleware.RequireModule("demo_links.crear")).Post("/", demoLinkHandler.Create)
				r.With(middleware.RequireModule("demo_links.eliminar")).Delete("/{id}", demoLinkHandler.Revoke)
			})
		})
	})

	return r
}
