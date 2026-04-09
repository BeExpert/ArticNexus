package db

import (
	"fmt"

	"articnexus/backend/internal/config"
	"articnexus/backend/internal/domain"
	"articnexus/backend/pkg/logger"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedSuperAdmin ensures a super-admin user exists in the database.
// When cfg.SuperAdminForce is true the password is re-hashed and updated on
// every boot; useful when you want to rotate the seed credential quickly.
func SeedSuperAdmin(database *gorm.DB, cfg *config.Config) error {
	if cfg.SuperAdminUser == "" || cfg.SuperAdminPass == "" {
		logger.Info(logger.App, "[seed] SUPER_ADMIN_USER or SUPER_ADMIN_PASS not set — skipping")
		return nil
	}

	cost := cfg.BcryptCost
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(cfg.SuperAdminPass), cost)
	if err != nil {
		return err
	}

	return database.Transaction(func(tx *gorm.DB) error {
		var user domain.User
		result := tx.Where("usr_username = ?", cfg.SuperAdminUser).Limit(1).Find(&user)
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected > 0 {
			// User already exists.
			if !cfg.SuperAdminForce {
				logger.Info(logger.App, fmt.Sprintf("[seed] super-admin '%s' already exists — skipping", cfg.SuperAdminUser))
				return nil
			}
			// Force mode: update password only.
			if err := tx.Model(&user).Update("usr_password", string(hash)).Error; err != nil {
				return err
			}
			logger.Info(logger.App, fmt.Sprintf("[seed] super-admin '%s' password refreshed (SUPER_ADMIN_FORCE=1)", cfg.SuperAdminUser))
			return nil
		}

		// Create a minimal Person record first.
		person := domain.Person{
			FirstName:    "Super",
			FirstSurname: "Admin",
			Status:       "active",
		}
		if err := tx.Create(&person).Error; err != nil {
			return err
		}

		// Create the User linked to that Person.
		// Super-admin status is NOT a DB column — it is derived at runtime
		// by matching this username against SUPER_ADMIN_USER in the env.
		user = domain.User{
			PersonID: person.ID,
			Username: cfg.SuperAdminUser,
			Email:    cfg.SuperAdminUser + "@articnexus.local",
			Password: string(hash),
			Status:   "active",
		}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		logger.Info(logger.App, fmt.Sprintf("[seed] super-admin '%s' created successfully", cfg.SuperAdminUser))
		return nil
	})
}

// SeedModules ensures all ArticNexus application modules are present in the
// database with up-to-date display names. Uses a raw SQL upsert keyed on
// (app_id, mod_name) so it is safe to run on every boot.
func SeedModules(database *gorm.DB) error {
	type mod struct{ name, display, desc string }
	modules := []mod{
		// ── Applications ────────────────────────────────────────────────
		{"aplicaciones.crear", "Crear aplicaciones", "Puede registrar nuevas aplicaciones"},
		{"aplicaciones.ver", "Ver aplicaciones", "Puede ver el listado de aplicaciones registradas"},
		{"aplicaciones.editar", "Editar aplicaciones", "Puede modificar aplicaciones existentes"},
		{"aplicaciones.eliminar", "Eliminar aplicaciones", "Puede eliminar aplicaciones del sistema"},
		{"aplicaciones.asignar_usuario", "Asignar usuarios a aplicaciones", "Puede asignar usuarios a aplicaciones"},
		{"aplicaciones.desasignar_usuario", "Desasignar usuarios de aplicaciones", "Puede desasignar usuarios de aplicaciones"},
		// ── Roles ───────────────────────────────────────────────────────
		{"roles.crear", "Crear roles", "Puede crear nuevos roles"},
		{"roles.ver", "Ver roles", "Puede ver los roles disponibles"},
		{"roles.editar", "Editar roles", "Puede modificar roles existentes"},
		{"roles.eliminar", "Eliminar roles", "Puede eliminar roles del sistema"},
		{"roles.asignar_usuario", "Asignar usuarios a roles", "Puede asignar usuarios a roles"},
		{"roles.desasignar_usuario", "Desasignar usuarios de roles", "Puede desasignar usuarios de roles"},
		{"roles.asignar_modulo", "Asignar módulos a roles", "Puede asignar módulos a roles"},
		{"roles.desasignar_modulo", "Desasignar módulos de roles", "Puede desasignar módulos de roles"},
		// ── Users ───────────────────────────────────────────────────────
		{"usuarios.crear", "Crear usuarios", "Puede crear nuevas cuentas de usuario"},
		{"usuarios.ver", "Ver usuarios", "Puede ver el listado de usuarios del sistema"},
		{"usuarios.editar", "Editar usuarios", "Puede modificar cuentas de usuario"},
		{"usuarios.eliminar", "Eliminar usuarios", "Puede eliminar cuentas de usuario"},
		{"usuarios.reset_contrasena", "Resetear contraseña de usuario", "Puede restablecer la contraseña de otros usuarios"},
		{"usuarios.demo", "Gestionar usuarios demo", "Puede configurar usuarios con acceso temporal (fecha de expiración)"},
		// ── Demo Links ──────────────────────────────────────────────────
		{"demo_links.ver", "Ver links de demo", "Puede ver los links de demo generados"},
		{"demo_links.crear", "Crear links de demo", "Puede generar nuevos links de acceso demo"},
		{"demo_links.eliminar", "Revocar links de demo", "Puede revocar links de demo activos"},
		// ── Modules ─────────────────────────────────────────────────────
		{"modulos.crear", "Crear módulos", "Puede crear nuevos módulos"},
		{"modulos.ver", "Ver módulos", "Puede ver los módulos de una aplicación"},
		{"modulos.editar", "Editar módulos", "Puede modificar módulos existentes"},
		{"modulos.eliminar", "Eliminar módulos", "Puede eliminar módulos del sistema"},
		// ── Companies ───────────────────────────────────────────────────
		{"empresas.crear", "Registrar empresas", "Puede registrar nuevas empresas en el sistema"},
		{"empresas.ver", "Ver empresas", "Puede ver el listado de empresas afiliadas"},
		{"empresas.editar", "Editar empresas", "Puede modificar los datos de empresas existentes"},
		{"empresas.eliminar", "Eliminar empresas", "Puede eliminar empresas del sistema"},
		// ── Persons ─────────────────────────────────────────────────────
		{"personas.ver", "Ver personas", "Puede ver las personas asignadas a una empresa"},
		{"personas.editar", "Editar personas", "Puede modificar los datos de personas"},
		{"personas.eliminar", "Quitar personas", "Puede quitar personas de una empresa"},
	}

	sql := `
INSERT INTO "tblModules_MOD"
    (app_id, mod_name, mod_display_name, mod_description, mod_status, created_at, updated_at)
SELECT a.app_id, @name, @display, @desc, 'active', now(), now()
FROM "tblApplications_APP" a
WHERE a.app_code = 'ARTICNEXUS'
ON CONFLICT (app_id, mod_name)
DO UPDATE SET
    mod_display_name = EXCLUDED.mod_display_name,
    mod_description  = EXCLUDED.mod_description,
    updated_at       = now()
`

	for _, m := range modules {
		if err := database.Exec(sql,
			map[string]interface{}{"name": m.name, "display": m.display, "desc": m.desc},
		).Error; err != nil {
			return fmt.Errorf("seed module %q: %w", m.name, err)
		}
	}

	logger.Info(logger.App, fmt.Sprintf("[seed] %d ArticNexus modules upserted", len(modules)))
	return nil
}

// SeedApplications ensures that the three core applications (ARTICNEXUS,
// OFTADATA, VETDATA) exist in the database along with all their modules.
// It is fully idempotent — safe to call on every boot.
func SeedApplications(database *gorm.DB) error {
	// ── Upsert applications ─────────────────────────────────────────────────
	apps := []struct{ code, name string }{
		{"ARTICNEXUS", "ArticNexus"},
		{"OFTADATA", "OftaData"},
		{"VETDATA", "VetData"},
	}
	for _, a := range apps {
		err := database.Exec(`
			INSERT INTO "tblApplications_APP" (app_code, app_name, app_status, created_at, updated_at)
			VALUES (@code, @name, 'active', now(), now())
			ON CONFLICT (app_code) DO NOTHING`,
			map[string]interface{}{"code": a.code, "name": a.name},
		).Error
		if err != nil {
			return fmt.Errorf("seed application %q: %w", a.code, err)
		}
	}

	// ── Module upsert template (same pattern as SeedModules) ────────────────
	modSQL := `
INSERT INTO "tblModules_MOD"
    (app_id, mod_name, mod_display_name, mod_description, mod_status, created_at, updated_at)
SELECT a.app_id, @name, @display, @desc, 'active', now(), now()
FROM "tblApplications_APP" a
WHERE a.app_code = @code
ON CONFLICT (app_id, mod_name)
DO UPDATE SET
    mod_display_name = EXCLUDED.mod_display_name,
    mod_description  = EXCLUDED.mod_description,
    updated_at       = now()
`
	type modDef struct{ code, name, display, desc string }

	// ── OftaData modules ────────────────────────────────────────────────────
	oftaModules := []modDef{
		// Pacientes
		{"OFTADATA", "pacientes.ver", "Ver pacientes", "Puede ver el listado de pacientes"},
		{"OFTADATA", "pacientes.crear", "Registrar pacientes", "Puede registrar nuevos pacientes"},
		{"OFTADATA", "pacientes.editar", "Editar pacientes", "Puede modificar los datos de pacientes"},
		{"OFTADATA", "pacientes.eliminar", "Eliminar pacientes", "Puede eliminar pacientes del sistema"},
		{"OFTADATA", "pacientes.reactivar", "Reactivar pacientes", "Puede reactivar pacientes inactivos"},
		// Formularios
		{"OFTADATA", "formularios.ver", "Ver formularios", "Puede ver los formularios clínicos"},
		{"OFTADATA", "formularios.crear", "Crear formulario", "Puede crear formularios clínicos"},
		{"OFTADATA", "formularios.editar", "Editar formulario", "Puede modificar formularios existentes"},
		{"OFTADATA", "formularios.eliminar", "Eliminar formulario", "Puede eliminar formularios del sistema"},
		{"OFTADATA", "formularios.reactivar", "Reactivar formulario", "Puede reactivar formularios inactivos"},
		// Citas
		{"OFTADATA", "citas.ver", "Ver agenda", "Puede ver las citas agendadas"},
		{"OFTADATA", "citas.crear", "Crear cita", "Puede crear nuevas citas"},
		{"OFTADATA", "citas.editar", "Editar cita", "Puede modificar citas existentes"},
		{"OFTADATA", "citas.eliminar", "Cancelar cita", "Puede cancelar o eliminar citas"},
		// Cirugía
		{"OFTADATA", "cirugia.ver", "Ver cirugías", "Puede ver los registros de cirugías"},
		{"OFTADATA", "cirugia.crear", "Registrar cirugía", "Puede registrar nuevas cirugías"},
		{"OFTADATA", "cirugia.editar", "Editar cirugía", "Puede modificar registros de cirugías"},
		{"OFTADATA", "cirugia.eliminar", "Eliminar cirugía", "Puede eliminar registros de cirugías"},
		// Catálogos
		{"OFTADATA", "catalogos.ver", "Ver catálogos", "Puede ver los catálogos del sistema"},
		{"OFTADATA", "catalogos.gestionar", "Gestionar catálogos", "Puede crear y modificar catálogos"},
		// Notificaciones
		{"OFTADATA", "notificaciones.ver", "Ver notificaciones", "Puede ver las notificaciones del sistema"},
		{"OFTADATA", "notificaciones.gestionar", "Gestionar notificaciones", "Puede configurar y gestionar notificaciones"},
		// Admin
		{"OFTADATA", "admin.usuarios", "Administrar usuarios", "Puede gestionar usuarios de la aplicación"},
		{"OFTADATA", "admin.reset_db", "Resetear base de datos", "Puede reinicializar los datos de la base de datos"},
	}

	// ── VetData modules ──────────────────────────────────────────────────────
	vetModules := []modDef{
		// Clientes
		{"VETDATA", "clientes.ver", "Ver clientes", "Puede ver el listado de clientes"},
		{"VETDATA", "clientes.crear", "Registrar clientes", "Puede registrar nuevos clientes"},
		// Mascotas
		{"VETDATA", "mascotas.ver", "Ver mascotas", "Puede ver el listado de mascotas"},
		{"VETDATA", "mascotas.crear", "Registrar mascotas", "Puede registrar nuevas mascotas"},
		// Vacunas
		{"VETDATA", "vacunas.ver", "Ver vacunas", "Puede ver el historial de vacunación"},
		{"VETDATA", "vacunas.crear", "Registrar vacunas", "Puede registrar nuevas vacunas"},
		{"VETDATA", "vacunas.aplicar", "Aplicar vacunas", "Puede registrar la aplicación de vacunas"},
		// Medicamentos
		{"VETDATA", "medicamentos.ver", "Ver medicamentos", "Puede ver el catálogo de medicamentos"},
		{"VETDATA", "medicamentos.crear", "Registrar medicamentos", "Puede registrar nuevos medicamentos"},
		// Tickets
		{"VETDATA", "tickets.ver", "Ver tickets", "Puede ver los tickets de soporte"},
		{"VETDATA", "tickets.responder", "Responder tickets", "Puede responder y gestionar tickets"},
		// Inventario
		{"VETDATA", "inventario.ver", "Ver inventario", "Puede ver el inventario de productos"},
		{"VETDATA", "inventario.crear", "Gestionar inventario", "Puede agregar y modificar el inventario"},
		// Consultas
		{"VETDATA", "consultas.ver", "Ver consultas", "Puede ver el historial de consultas clínicas"},
		{"VETDATA", "consultas.crear", "Registrar consultas", "Puede crear y editar consultas clínicas"},
		// Cirugías
		{"VETDATA", "cirugias.ver", "Ver cirugías", "Puede ver los registros de cirugías de mascotas"},
		{"VETDATA", "cirugias.crear", "Registrar cirugías", "Puede registrar y editar cirugías"},
		// Libro
		{"VETDATA", "libro.ver", "Ver libro de mascota", "Puede ver el libro de mascota completo"},
		{"VETDATA", "libro.crear", "Editar libro de mascota", "Puede crear registros en el libro de mascota"},
		// Agenda
		{"VETDATA", "agenda.ver", "Ver agenda", "Puede ver la agenda, citas y calendario"},
		{"VETDATA", "agenda.crear", "Gestionar agenda", "Puede crear y gestionar citas y disponibilidad"},
		// Admin
		{"VETDATA", "admin.usuarios", "Administrar usuarios", "Puede gestionar usuarios de la aplicación"},
	}

	allMods := append(oftaModules, vetModules...)
	for _, m := range allMods {
		if err := database.Exec(modSQL, map[string]interface{}{
			"code":    m.code,
			"name":    m.name,
			"display": m.display,
			"desc":    m.desc,
		}).Error; err != nil {
			return fmt.Errorf("seed module %q (%s): %w", m.name, m.code, err)
		}
	}

	logger.Info(logger.App, fmt.Sprintf("[seed] %d OftaData + VetData modules upserted", len(allMods)))
	return nil
}
