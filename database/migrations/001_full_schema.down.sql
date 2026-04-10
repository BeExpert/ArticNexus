-- =============================================================
-- ArticNexus — Full Schema Rollback
-- File: database/migrations/001_full_schema.down.sql
--
-- Drops ALL tables, indexes, triggers, and functions.
-- USE WITH EXTREME CAUTION — this destroys all data.
--
-- Usage (manual only — never run automatically):
--   psql $DATABASE_URL -f 001_full_schema.down.sql
-- =============================================================

BEGIN;

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS "tblDemoLinks_DML" CASCADE;
DROP TABLE IF EXISTS "tblPasswordResetTokens_PRT" CASCADE;
DROP TABLE IF EXISTS "tblUserRoles_URO" CASCADE;
DROP TABLE IF EXISTS "tblRoleModules_RMO" CASCADE;
DROP TABLE IF EXISTS "tblModules_MOD" CASCADE;
DROP TABLE IF EXISTS "tblRoles_ROL" CASCADE;
DROP TABLE IF EXISTS "tblCompanyApplications_CAP" CASCADE;
DROP TABLE IF EXISTS "tblUserBranches_UBR" CASCADE;
DROP TABLE IF EXISTS "tblUserCompanies_UCO" CASCADE;
DROP TABLE IF EXISTS "tblBranches_BRA" CASCADE;
DROP TABLE IF EXISTS "tblApplications_APP" CASCADE;
DROP TABLE IF EXISTS "tblCompanies_COM" CASCADE;
DROP TABLE IF EXISTS "tblUsers_USR" CASCADE;
DROP TABLE IF EXISTS "tblPersons_PER" CASCADE;

-- Drop function (triggers are dropped with their tables via CASCADE)
DROP FUNCTION IF EXISTS fn_set_updated_at() CASCADE;

COMMIT;
