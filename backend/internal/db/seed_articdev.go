package db

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"articnexus/backend/internal/config"
	"articnexus/backend/pkg/logger"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedArticDevAndDemoUsers ensures that:
//   - Company "ArticDev" and its headquarters branch exist.
//   - For each client application (OFTADATA, VETDATA) a dedicated demo user
//     (demo_oftadata / demo_vetdata) exists with a "DEMO" role that grants
//     read-only (*.ver) access to all modules of that application.
//
// The function is fully idempotent — safe to call on every boot.
func SeedArticDevAndDemoUsers(database *gorm.DB, cfg *config.Config) error {
	cost := cfg.BcryptCost
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}

	return database.Transaction(func(tx *gorm.DB) error {

		// ── 1. Upsert company "ArticDev" ─────────────────────────────────────
		var companyID int64
		if err := tx.Raw(
			`SELECT com_id FROM "tblCompanies_COM" WHERE com_name = 'ArticDev' AND deleted_at IS NULL LIMIT 1`,
		).Scan(&companyID).Error; err != nil {
			return fmt.Errorf("seed articdev: check company: %w", err)
		}
		if companyID == 0 {
			if err := tx.Raw(
				`INSERT INTO "tblCompanies_COM" (com_name, com_status, created_at, updated_at)
				 VALUES ('ArticDev', 'active', now(), now())
				 RETURNING com_id`,
			).Scan(&companyID).Error; err != nil {
				return fmt.Errorf("seed articdev: insert company: %w", err)
			}
		}

		// ── 2. Upsert branch "Matriz" (bra_code = 'HQ') ──────────────────────
		if err := tx.Exec(
			`INSERT INTO "tblBranches_BRA" (com_id, bra_code, bra_name, bra_status, created_at, updated_at)
			 VALUES (?, 'HQ', 'Matriz', 'active', now(), now())
			 ON CONFLICT (com_id, bra_code) DO NOTHING`,
			companyID,
		).Error; err != nil {
			return fmt.Errorf("seed articdev: insert branch: %w", err)
		}

		var branchID int64
		if err := tx.Raw(
			`SELECT bra_id FROM "tblBranches_BRA" WHERE com_id = ? AND bra_code = 'HQ' LIMIT 1`,
			companyID,
		).Scan(&branchID).Error; err != nil {
			return fmt.Errorf("seed articdev: select branch: %w", err)
		}

		// ── 2b. License all active apps for ArticDev ──────────────────────────
		if err := tx.Exec(
			`INSERT INTO "tblCompanyApplications_CAP" (com_id, app_id, cap_status)
			 SELECT ?, app_id, 'active'
			 FROM "tblApplications_APP"
			 WHERE app_status = 'active'
			 ON CONFLICT (com_id, app_id) DO NOTHING`,
			companyID,
		).Error; err != nil {
			return fmt.Errorf("seed articdev: license apps: %w", err)
		}

		// ── 3. Per-application demo users + DEMO roles ────────────────────────
		apps := []struct{ code, label string }{
			{"OFTADATA", "OftaData"},
			{"VETDATA", "VetData"},
		}

		for _, app := range apps {
			username := "demo_" + strings.ToLower(app.code)
			email := username + "@articdev.local"

			// ── 3a. Upsert person + user ──────────────────────────────────────
			var userID int64
			if err := tx.Raw(
				`SELECT usr_id FROM "tblUsers_USR" WHERE usr_username = ? LIMIT 1`,
				username,
			).Scan(&userID).Error; err != nil {
				return fmt.Errorf("seed articdev: check user %s: %w", username, err)
			}

			if userID == 0 {
				// Random initialization password — never shown, never re-used.
				passBytes := make([]byte, 32)
				if _, err := rand.Read(passBytes); err != nil {
					return fmt.Errorf("seed articdev: generate password for %s: %w", username, err)
				}
				hash, err := bcrypt.GenerateFromPassword([]byte(hex.EncodeToString(passBytes)), cost)
				if err != nil {
					return fmt.Errorf("seed articdev: hash password for %s: %w", username, err)
				}

				// Insert person.
				var personID int64
				if err := tx.Raw(
					`INSERT INTO "tblPersons_PER" (per_firstname, per_firstsurname, per_status, created_at, updated_at)
					 VALUES ('Demo', ?, 'active', now(), now())
					 RETURNING per_id`,
					app.label,
				).Scan(&personID).Error; err != nil {
					return fmt.Errorf("seed articdev: insert person for %s: %w", username, err)
				}

				// Insert user.
				if err := tx.Raw(
					`INSERT INTO "tblUsers_USR" (per_id, usr_username, usr_email, usr_password, usr_status, created_at, updated_at)
					 VALUES (?, ?, ?, ?, 'active', now(), now())
					 RETURNING usr_id`,
					personID, username, email, string(hash),
				).Scan(&userID).Error; err != nil {
					return fmt.Errorf("seed articdev: insert user %s: %w", username, err)
				}
			}

			// ── 3b. Resolve application ID ────────────────────────────────────
			var appID int64
			if err := tx.Raw(
				`SELECT app_id FROM "tblApplications_APP" WHERE app_code = ? LIMIT 1`,
				app.code,
			).Scan(&appID).Error; err != nil {
				return fmt.Errorf("seed articdev: get app_id for %s: %w", app.code, err)
			}
			if appID == 0 {
				return fmt.Errorf("seed articdev: application %s not found — run SeedApplications first", app.code)
			}

			// ── 3c. Upsert role "DEMO" ────────────────────────────────────────
			if err := tx.Exec(
				`INSERT INTO "tblRoles_ROL" (app_id, rol_name, rol_status, created_at, updated_at)
				 VALUES (?, 'DEMO', 'active', now(), now())
				 ON CONFLICT (app_id, rol_name) DO NOTHING`,
				appID,
			).Error; err != nil {
				return fmt.Errorf("seed articdev: insert role DEMO for %s: %w", app.code, err)
			}

			var roleID int64
			if err := tx.Raw(
				`SELECT rol_id FROM "tblRoles_ROL" WHERE app_id = ? AND rol_name = 'DEMO' LIMIT 1`,
				appID,
			).Scan(&roleID).Error; err != nil {
				return fmt.Errorf("seed articdev: get role_id for %s: %w", app.code, err)
			}

			// ── 3d. Assign all *.ver modules to the DEMO role ─────────────────
			if err := tx.Exec(
				`INSERT INTO "tblRoleModules_RMO" (rol_id, mod_id, created_at)
				 SELECT ?, mod_id, now()
				 FROM "tblModules_MOD"
				 WHERE app_id = ? AND mod_name LIKE '%.ver'
				 ON CONFLICT (rol_id, mod_id) DO NOTHING`,
				roleID, appID,
			).Error; err != nil {
				return fmt.Errorf("seed articdev: assign modules for %s: %w", app.code, err)
			}

			// ── 3e. Junction-table rows ───────────────────────────────────────
			if err := tx.Exec(
				`INSERT INTO "tblUserCompanies_UCO" (usr_id, com_id, created_at)
				 VALUES (?, ?, now())
				 ON CONFLICT (usr_id, com_id) DO NOTHING`,
				userID, companyID,
			).Error; err != nil {
				return fmt.Errorf("seed articdev: user company for %s: %w", username, err)
			}

			if err := tx.Exec(
				`INSERT INTO "tblUserBranches_UBR" (usr_id, bra_id, created_at)
				 VALUES (?, ?, now())
				 ON CONFLICT (usr_id, bra_id) DO NOTHING`,
				userID, branchID,
			).Error; err != nil {
				return fmt.Errorf("seed articdev: user branch for %s: %w", username, err)
			}

			if err := tx.Exec(
				`INSERT INTO "tblUserRoles_URO" (usr_id, com_id, bra_id, rol_id, created_at)
				 VALUES (?, ?, ?, ?, now())
				 ON CONFLICT (usr_id, com_id, bra_id, rol_id) DO NOTHING`,
				userID, companyID, branchID, roleID,
			).Error; err != nil {
				return fmt.Errorf("seed articdev: user role for %s: %w", username, err)
			}
		}

		logger.Info(logger.App, "[seed] ArticDev company, branch, demo users, and DEMO roles upserted")
		return nil
	})
}
