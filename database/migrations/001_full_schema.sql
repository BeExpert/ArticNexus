-- =============================================================
-- ArticNexus — Full Schema + Seed (Production Initial State)
-- This single file represents the canonical database state.
-- Run once on a fresh database: psql $DATABASE_URL -f 001_full_schema.sql
-- Idempotent: safe to apply on already-initialized databases.
-- =============================================================

BEGIN;

-- =============================================================
-- SECTION 1: TABLES
-- =============================================================

CREATE TABLE IF NOT EXISTS "tblPersons_PER" (
    per_id             BIGSERIAL    PRIMARY KEY,
    per_firstName      VARCHAR(100),
    per_firstSurname   VARCHAR(100),
    per_secondSurname  VARCHAR(100),
    per_nationalId     VARCHAR(50),
    per_email          VARCHAR(255),
    per_birthDate      DATE,
    per_phoneAreaCode  VARCHAR(10),
    per_primaryPhone   VARCHAR(20),
    per_secondaryPhone VARCHAR(20),
    per_address        TEXT,
    per_status         VARCHAR(20)  NOT NULL DEFAULT 'active',
    created_at         TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at         TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS "tblUsers_USR" (
    usr_id                  BIGSERIAL    PRIMARY KEY,
    per_id                  BIGINT       NOT NULL REFERENCES "tblPersons_PER"(per_id),
    usr_username            VARCHAR(100) UNIQUE,
    usr_email               VARCHAR(255) UNIQUE,
    usr_password            VARCHAR(255),
    usr_status              VARCHAR(20),
    usr_password_expires_at TIMESTAMPTZ,
    created_at              TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at              TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS "tblCompanies_COM" (
    com_id     BIGSERIAL    PRIMARY KEY,
    com_name   VARCHAR(255) NOT NULL,
    com_status VARCHAR(20)  NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS "tblBranches_BRA" (
    bra_id          BIGSERIAL    PRIMARY KEY,
    com_id          BIGINT       NOT NULL REFERENCES "tblCompanies_COM"(com_id),
    bra_code        VARCHAR(50),
    bra_name        VARCHAR(150),
    bra_address     TEXT,
    bra_phoneNumber VARCHAR(20),
    bra_email       VARCHAR(255),
    bra_status      VARCHAR(20),
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ,
    UNIQUE (com_id, bra_code)
);

CREATE TABLE IF NOT EXISTS "tblUserCompanies_UCO" (
    usr_id     BIGINT      NOT NULL REFERENCES "tblUsers_USR"(usr_id)    ON DELETE CASCADE,
    com_id     BIGINT      NOT NULL REFERENCES "tblCompanies_COM"(com_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (usr_id, com_id)
);

CREATE TABLE IF NOT EXISTS "tblUserBranches_UBR" (
    usr_id     BIGINT      NOT NULL REFERENCES "tblUsers_USR"(usr_id)    ON DELETE CASCADE,
    bra_id     BIGINT      NOT NULL REFERENCES "tblBranches_BRA"(bra_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (usr_id, bra_id)
);

CREATE TABLE IF NOT EXISTS "tblApplications_APP" (
    app_id     BIGSERIAL    PRIMARY KEY,
    app_code   VARCHAR(50)  NOT NULL UNIQUE,
    app_name   VARCHAR(255) NOT NULL,
    app_status VARCHAR(20)  NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS "tblRoles_ROL" (
    rol_id     BIGSERIAL    PRIMARY KEY,
    app_id     BIGINT       NOT NULL REFERENCES "tblApplications_APP"(app_id),
    rol_name   VARCHAR(100) NOT NULL,
    rol_status VARCHAR(20)  NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    UNIQUE (app_id, rol_name)
);

CREATE TABLE IF NOT EXISTS "tblModules_MOD" (
    mod_id           BIGSERIAL    PRIMARY KEY,
    app_id           BIGINT       NOT NULL REFERENCES "tblApplications_APP"(app_id),
    mod_name         VARCHAR(100),
    mod_menuOption   VARCHAR(100),
    mod_subFunction  VARCHAR(100),
    mod_description  TEXT,
    mod_status       VARCHAR(20)  DEFAULT 'active',
    mod_display_name VARCHAR(255),
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at       TIMESTAMPTZ,
    CONSTRAINT uq_modules_app_mod_name UNIQUE (app_id, mod_name)
);

CREATE TABLE IF NOT EXISTS "tblRoleModules_RMO" (
    rol_id     BIGINT      NOT NULL REFERENCES "tblRoles_ROL"(rol_id)    ON DELETE CASCADE,
    mod_id     BIGINT      NOT NULL REFERENCES "tblModules_MOD"(mod_id)  ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (rol_id, mod_id)
);

CREATE TABLE IF NOT EXISTS "tblUserRoles_URO" (
    usr_id     BIGINT      NOT NULL REFERENCES "tblUsers_USR"(usr_id)       ON DELETE CASCADE,
    com_id     BIGINT      NOT NULL REFERENCES "tblCompanies_COM"(com_id)   ON DELETE CASCADE,
    bra_id     BIGINT      NOT NULL REFERENCES "tblBranches_BRA"(bra_id)   ON DELETE CASCADE,
    rol_id     BIGINT      NOT NULL REFERENCES "tblRoles_ROL"(rol_id)       ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (usr_id, com_id, bra_id, rol_id)
);

CREATE TABLE IF NOT EXISTS "tblPasswordResetTokens_PRT" (
    prt_id         SERIAL      PRIMARY KEY,
    usr_id         INTEGER     NOT NULL REFERENCES "tblUsers_USR"(usr_id) ON DELETE CASCADE,
    prt_token_hash VARCHAR(255) NOT NULL,
    prt_expires_at TIMESTAMPTZ NOT NULL,
    prt_used       BOOLEAN     DEFAULT FALSE,
    prt_created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS "tblDemoLinks_DML" (
    dml_id              BIGSERIAL    PRIMARY KEY,
    dml_token_hash      VARCHAR(255) NOT NULL UNIQUE,
    dml_app_code        VARCHAR(50)  NOT NULL,
    dml_demo_user_id    BIGINT       NOT NULL REFERENCES "tblUsers_USR"(usr_id),
    dml_expires_at      TIMESTAMPTZ  NOT NULL,
    dml_is_active       BOOLEAN      NOT NULL DEFAULT TRUE,
    dml_recipient_email VARCHAR(255),
    dml_created_by      BIGINT       NOT NULL REFERENCES "tblUsers_USR"(usr_id),
    dml_created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- =============================================================
-- SECTION 2: INDEXES
-- =============================================================

CREATE INDEX IF NOT EXISTS idx_prt_token_hash ON "tblPasswordResetTokens_PRT"(prt_token_hash);
CREATE INDEX IF NOT EXISTS idx_prt_usr_id     ON "tblPasswordResetTokens_PRT"(usr_id);
CREATE INDEX IF NOT EXISTS idx_demolinks_token_hash ON "tblDemoLinks_DML"(dml_token_hash);
CREATE INDEX IF NOT EXISTS idx_demolinks_app_code   ON "tblDemoLinks_DML"(dml_app_code);

-- =============================================================
-- SECTION 3: TRIGGERS (auto-update updated_at)
-- =============================================================

CREATE OR REPLACE FUNCTION fn_set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_persons_updated_at') THEN
        CREATE TRIGGER trg_persons_updated_at
            BEFORE UPDATE ON "tblPersons_PER"
            FOR EACH ROW EXECUTE FUNCTION fn_set_updated_at();
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_users_updated_at') THEN
        CREATE TRIGGER trg_users_updated_at
            BEFORE UPDATE ON "tblUsers_USR"
            FOR EACH ROW EXECUTE FUNCTION fn_set_updated_at();
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_companies_updated_at') THEN
        CREATE TRIGGER trg_companies_updated_at
            BEFORE UPDATE ON "tblCompanies_COM"
            FOR EACH ROW EXECUTE FUNCTION fn_set_updated_at();
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_branches_updated_at') THEN
        CREATE TRIGGER trg_branches_updated_at
            BEFORE UPDATE ON "tblBranches_BRA"
            FOR EACH ROW EXECUTE FUNCTION fn_set_updated_at();
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_applications_updated_at') THEN
        CREATE TRIGGER trg_applications_updated_at
            BEFORE UPDATE ON "tblApplications_APP"
            FOR EACH ROW EXECUTE FUNCTION fn_set_updated_at();
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_roles_updated_at') THEN
        CREATE TRIGGER trg_roles_updated_at
            BEFORE UPDATE ON "tblRoles_ROL"
            FOR EACH ROW EXECUTE FUNCTION fn_set_updated_at();
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_modules_updated_at') THEN
        CREATE TRIGGER trg_modules_updated_at
            BEFORE UPDATE ON "tblModules_MOD"
            FOR EACH ROW EXECUTE FUNCTION fn_set_updated_at();
    END IF;
END $$;

-- =============================================================
-- SECTION 4: SEED DATA — APPLICATIONS
-- =============================================================

INSERT INTO "tblApplications_APP" (app_code, app_name, app_status, created_at, updated_at)
VALUES
    ('ARTICNEXUS', 'ArticNexus', 'active', NOW(), NOW()),
    ('OFTADATA',   'OftaData',   'active', NOW(), NOW()),
    ('VETDATA',    'VetData',    'active', NOW(), NOW())
ON CONFLICT (app_code) DO NOTHING;

-- =============================================================
-- SECTION 5: SEED DATA — ARTICNEXUS MODULES
-- =============================================================

WITH app AS (SELECT app_id FROM "tblApplications_APP" WHERE app_code = 'ARTICNEXUS')
INSERT INTO "tblModules_MOD" (app_id, mod_name, mod_display_name, mod_description, mod_status, created_at, updated_at)
SELECT app.app_id, v.n, v.d, v.descr, 'active', NOW(), NOW()
FROM app, (VALUES
    ('dashboard.ver',               'Ver Dashboard',                   'Puede ver el panel de control'),
    ('empresas.ver',                'Ver empresas',                    'Puede ver el listado de empresas'),
    ('empresas.crear',              'Registrar empresas',              'Puede registrar nuevas empresas'),
    ('empresas.editar',             'Editar empresas',                 'Puede modificar datos de empresas'),
    ('empresas.eliminar',           'Eliminar empresas',               'Puede eliminar empresas del sistema'),
    ('sucursales.ver',              'Ver sucursales',                  'Puede ver las sucursales de una empresa'),
    ('sucursales.crear',            'Crear sucursal',                  'Puede registrar nuevas sucursales'),
    ('sucursales.editar',           'Editar sucursal',                 'Puede modificar datos de sucursales'),
    ('sucursales.eliminar',         'Eliminar sucursal',               'Puede eliminar sucursales'),
    ('personas.ver',                'Ver personas',                    'Puede ver personas asignadas a una empresa'),
    ('personas.editar',             'Editar personas',                 'Puede modificar datos de personas'),
    ('personas.eliminar',           'Quitar personas',                 'Puede quitar personas de una empresa'),
    ('usuarios.ver',                'Ver usuarios',                    'Puede ver el listado de usuarios'),
    ('usuarios.crear',              'Crear usuarios',                  'Puede crear nuevas cuentas de usuario'),
    ('usuarios.editar',             'Editar usuarios',                 'Puede modificar cuentas de usuario'),
    ('usuarios.eliminar',           'Eliminar usuarios',               'Puede eliminar cuentas de usuario'),
    ('usuarios.reset_contrasena',   'Resetear contraseña de usuario',  'Puede restablecer la contraseña de otros usuarios'),
    ('usuarios.demo',               'Gestionar usuarios demo',         'Puede configurar usuarios con acceso temporal'),
    ('roles.ver',                   'Ver roles',                       'Puede ver los roles disponibles'),
    ('roles.crear',                 'Crear roles',                     'Puede crear nuevos roles'),
    ('roles.editar',                'Editar roles',                    'Puede modificar roles existentes'),
    ('roles.eliminar',              'Eliminar roles',                  'Puede eliminar roles del sistema'),
    ('roles.asignar_usuario',       'Asignar usuarios a roles',        'Puede asignar usuarios a roles'),
    ('roles.desasignar_usuario',    'Desasignar usuarios de roles',    'Puede desasignar usuarios de roles'),
    ('roles.asignar_modulo',        'Asignar módulos a roles',         'Puede asignar módulos a roles'),
    ('roles.desasignar_modulo',     'Desasignar módulos de roles',     'Puede desasignar módulos de roles'),
    ('modulos.ver',                 'Ver módulos',                     'Puede ver los módulos de una aplicación'),
    ('modulos.crear',               'Crear módulos',                   'Puede crear nuevos módulos'),
    ('modulos.editar',              'Editar módulos',                  'Puede modificar módulos existentes'),
    ('modulos.eliminar',            'Eliminar módulos',                'Puede eliminar módulos del sistema'),
    ('aplicaciones.ver',            'Ver aplicaciones',                'Puede ver el listado de aplicaciones'),
    ('aplicaciones.crear',          'Crear aplicaciones',              'Puede registrar nuevas aplicaciones'),
    ('aplicaciones.editar',         'Editar aplicaciones',             'Puede modificar aplicaciones existentes'),
    ('aplicaciones.eliminar',       'Eliminar aplicaciones',           'Puede eliminar aplicaciones del sistema'),
    ('aplicaciones.asignar_usuario','Asignar usuarios a aplicaciones', 'Puede asignar usuarios a aplicaciones'),
    ('aplicaciones.desasignar_usuario','Desasignar usuarios de aplicaciones','Puede desasignar usuarios de aplicaciones'),
    ('demo_links.ver',              'Ver links de demo',               'Puede ver los links de demo generados'),
    ('demo_links.crear',            'Crear links de demo',             'Puede generar nuevos links de acceso demo'),
    ('demo_links.eliminar',         'Revocar links de demo',           'Puede revocar links de demo activos')
  ) AS v(n, d, descr)
ON CONFLICT (app_id, mod_name) DO UPDATE SET
    mod_display_name = EXCLUDED.mod_display_name,
    mod_description  = EXCLUDED.mod_description,
    updated_at       = NOW();

-- =============================================================
-- SECTION 6: SEED DATA — OFTADATA MODULES
-- =============================================================

WITH app AS (SELECT app_id FROM "tblApplications_APP" WHERE app_code = 'OFTADATA')
INSERT INTO "tblModules_MOD" (app_id, mod_name, mod_display_name, mod_description, mod_status, created_at, updated_at)
SELECT app.app_id, v.n, v.d, v.descr, 'active', NOW(), NOW()
FROM app, (VALUES
    ('pacientes.ver',           'Ver pacientes',            'Puede ver el listado de pacientes'),
    ('pacientes.crear',         'Registrar pacientes',      'Puede registrar nuevos pacientes'),
    ('pacientes.editar',        'Editar pacientes',         'Puede modificar los datos de pacientes'),
    ('pacientes.eliminar',      'Eliminar pacientes',       'Puede eliminar pacientes del sistema'),
    ('pacientes.reactivar',     'Reactivar pacientes',      'Puede reactivar pacientes inactivos'),
    ('formularios.ver',         'Ver formularios',          'Puede ver los formularios clínicos'),
    ('formularios.crear',       'Crear formulario',         'Puede crear formularios clínicos'),
    ('formularios.editar',      'Editar formulario',        'Puede modificar formularios existentes'),
    ('formularios.eliminar',    'Eliminar formulario',      'Puede eliminar formularios del sistema'),
    ('formularios.reactivar',   'Reactivar formulario',     'Puede reactivar formularios inactivos'),
    ('citas.ver',               'Ver agenda',               'Puede ver las citas agendadas'),
    ('citas.crear',             'Crear cita',               'Puede crear nuevas citas'),
    ('citas.editar',            'Editar cita',              'Puede modificar citas existentes'),
    ('citas.eliminar',          'Cancelar cita',            'Puede cancelar o eliminar citas'),
    ('cirugia.ver',             'Ver cirugías',             'Puede ver los registros de cirugías'),
    ('cirugia.crear',           'Registrar cirugía',        'Puede registrar nuevas cirugías'),
    ('cirugia.editar',          'Editar cirugía',           'Puede modificar registros de cirugías'),
    ('cirugia.eliminar',        'Eliminar cirugía',         'Puede eliminar registros de cirugías'),
    ('catalogos.ver',           'Ver catálogos',            'Puede ver los catálogos del sistema'),
    ('catalogos.gestionar',     'Gestionar catálogos',      'Puede crear y modificar catálogos'),
    ('notificaciones.ver',      'Ver notificaciones',       'Puede ver las notificaciones del sistema'),
    ('notificaciones.gestionar','Gestionar notificaciones', 'Puede configurar y gestionar notificaciones'),
    ('admin.usuarios',          'Administrar usuarios',     'Puede gestionar usuarios de la aplicación'),
    ('admin.reset_db',          'Resetear base de datos',   'Puede reinicializar los datos de la base de datos')
  ) AS v(n, d, descr)
ON CONFLICT (app_id, mod_name) DO UPDATE SET
    mod_display_name = EXCLUDED.mod_display_name,
    mod_description  = EXCLUDED.mod_description,
    updated_at       = NOW();

-- =============================================================
-- SECTION 7: SEED DATA — VETDATA MODULES
-- =============================================================

WITH app AS (SELECT app_id FROM "tblApplications_APP" WHERE app_code = 'VETDATA')
INSERT INTO "tblModules_MOD" (app_id, mod_name, mod_display_name, mod_description, mod_status, created_at, updated_at)
SELECT app.app_id, v.n, v.d, v.descr, 'active', NOW(), NOW()
FROM app, (VALUES
    ('clientes.ver',     'Ver clientes',          'Puede ver el listado de clientes'),
    ('clientes.crear',   'Registrar clientes',    'Puede registrar nuevos clientes'),
    ('mascotas.ver',     'Ver mascotas',          'Puede ver el listado de mascotas'),
    ('mascotas.crear',   'Registrar mascotas',    'Puede registrar nuevas mascotas'),
    ('vacunas.ver',      'Ver vacunas',           'Puede ver el historial de vacunación'),
    ('vacunas.crear',    'Registrar vacunas',     'Puede registrar nuevas vacunas'),
    ('vacunas.aplicar',  'Aplicar vacunas',       'Puede registrar la aplicación de vacunas'),
    ('medicamentos.ver', 'Ver medicamentos',      'Puede ver el catálogo de medicamentos'),
    ('medicamentos.crear','Registrar medicamentos','Puede registrar nuevos medicamentos'),
    ('tickets.ver',      'Ver tickets',           'Puede ver los tickets de soporte'),
    ('tickets.responder','Responder tickets',     'Puede responder y gestionar tickets'),
    ('inventario.ver',   'Ver inventario',        'Puede ver el inventario de productos'),
    ('inventario.crear', 'Gestionar inventario',  'Puede agregar y modificar el inventario'),
    ('consultas.ver',    'Ver consultas',         'Puede ver el historial de consultas clínicas'),
    ('consultas.crear',  'Registrar consultas',   'Puede crear y editar consultas clínicas'),
    ('cirugias.ver',     'Ver cirugías',          'Puede ver los registros de cirugías de mascotas'),
    ('cirugias.crear',   'Registrar cirugías',    'Puede registrar y editar cirugías'),
    ('libro.ver',        'Ver libro de mascota',  'Puede ver el libro de mascota completo'),
    ('libro.crear',      'Editar libro de mascota','Puede crear registros en el libro de mascota'),
    ('agenda.ver',       'Ver agenda',            'Puede ver la agenda, citas y calendario'),
    ('agenda.crear',     'Gestionar agenda',      'Puede crear y gestionar citas y disponibilidad'),
    ('admin.usuarios',   'Administrar usuarios',  'Puede gestionar usuarios de la aplicación')
  ) AS v(n, d, descr)
ON CONFLICT (app_id, mod_name) DO UPDATE SET
    mod_display_name = EXCLUDED.mod_display_name,
    mod_description  = EXCLUDED.mod_description,
    updated_at       = NOW();

-- =============================================================
-- SECTION 8: SEED DATA — ROLES
-- =============================================================

-- ArticNexus roles
INSERT INTO "tblRoles_ROL" (app_id, rol_name, rol_status, created_at, updated_at)
SELECT a.app_id, v.name, 'active', NOW(), NOW()
FROM "tblApplications_APP" a, (VALUES
    ('Super Administrador'),
    ('Administrador de Empresa')
) AS v(name)
WHERE a.app_code = 'ARTICNEXUS'
ON CONFLICT (app_id, rol_name) DO NOTHING;

-- OftaData roles
INSERT INTO "tblRoles_ROL" (app_id, rol_name, rol_status, created_at, updated_at)
SELECT a.app_id, v.name, 'active', NOW(), NOW()
FROM "tblApplications_APP" a, (VALUES
    ('SuperAdmin'),
    ('Secretaria'),
    ('Optometrista'),
    ('Recepcionista')
) AS v(name)
WHERE a.app_code = 'OFTADATA'
ON CONFLICT (app_id, rol_name) DO NOTHING;

-- VetData roles
INSERT INTO "tblRoles_ROL" (app_id, rol_name, rol_status, created_at, updated_at)
SELECT a.app_id, v.name, 'active', NOW(), NOW()
FROM "tblApplications_APP" a, (VALUES
    ('Admin'),
    ('Veterinario'),
    ('Asistente'),
    ('Recepcionista'),
    ('Cliente')
) AS v(name)
WHERE a.app_code = 'VETDATA'
ON CONFLICT (app_id, rol_name) DO NOTHING;

-- =============================================================
-- SECTION 9: SEED DATA — ROLE-MODULE ASSIGNMENTS
-- =============================================================

-- ArticNexus: Super Administrador → todos los módulos
INSERT INTO "tblRoleModules_RMO" (rol_id, mod_id, created_at)
SELECT r.rol_id, m.mod_id, NOW()
FROM "tblRoles_ROL" r
JOIN "tblApplications_APP" a ON a.app_id = r.app_id
JOIN "tblModules_MOD" m ON m.app_id = a.app_id
WHERE a.app_code = 'ARTICNEXUS'
  AND r.rol_name = 'Super Administrador'
  AND m.mod_status = 'active'
ON CONFLICT (rol_id, mod_id) DO NOTHING;

-- ArticNexus: Administrador de Empresa → subset
INSERT INTO "tblRoleModules_RMO" (rol_id, mod_id, created_at)
SELECT r.rol_id, m.mod_id, NOW()
FROM "tblRoles_ROL" r
JOIN "tblApplications_APP" a ON a.app_id = r.app_id
JOIN "tblModules_MOD" m ON m.app_id = a.app_id
WHERE a.app_code = 'ARTICNEXUS'
  AND r.rol_name = 'Administrador de Empresa'
  AND m.mod_name IN (
    'dashboard.ver',
    'empresas.ver', 'empresas.crear', 'empresas.editar',
    'sucursales.ver', 'sucursales.crear', 'sucursales.editar', 'sucursales.eliminar',
    'personas.ver', 'personas.editar', 'personas.eliminar',
    'roles.ver', 'roles.crear', 'roles.editar'
  )
ON CONFLICT (rol_id, mod_id) DO NOTHING;

-- OftaData: SuperAdmin → todos los módulos
INSERT INTO "tblRoleModules_RMO" (rol_id, mod_id, created_at)
SELECT r.rol_id, m.mod_id, NOW()
FROM "tblRoles_ROL" r
JOIN "tblApplications_APP" a ON a.app_id = r.app_id
JOIN "tblModules_MOD" m ON m.app_id = a.app_id
WHERE a.app_code = 'OFTADATA'
  AND r.rol_name = 'SuperAdmin'
  AND m.mod_status = 'active'
ON CONFLICT (rol_id, mod_id) DO NOTHING;

-- OftaData: Secretaria
INSERT INTO "tblRoleModules_RMO" (rol_id, mod_id, created_at)
SELECT r.rol_id, m.mod_id, NOW()
FROM "tblRoles_ROL" r
JOIN "tblApplications_APP" a ON a.app_id = r.app_id
JOIN "tblModules_MOD" m ON m.app_id = a.app_id
WHERE a.app_code = 'OFTADATA' AND r.rol_name = 'Secretaria'
  AND m.mod_name IN (
    'pacientes.ver', 'formularios.ver',
    'citas.ver', 'citas.crear', 'citas.editar', 'citas.eliminar',
    'cirugia.ver', 'notificaciones.ver'
  )
ON CONFLICT (rol_id, mod_id) DO NOTHING;

-- OftaData: Optometrista
INSERT INTO "tblRoleModules_RMO" (rol_id, mod_id, created_at)
SELECT r.rol_id, m.mod_id, NOW()
FROM "tblRoles_ROL" r
JOIN "tblApplications_APP" a ON a.app_id = r.app_id
JOIN "tblModules_MOD" m ON m.app_id = a.app_id
WHERE a.app_code = 'OFTADATA' AND r.rol_name = 'Optometrista'
  AND m.mod_name IN (
    'pacientes.ver', 'pacientes.crear', 'pacientes.editar',
    'formularios.ver', 'formularios.crear', 'formularios.editar',
    'citas.ver', 'citas.crear',
    'cirugia.ver', 'cirugia.crear', 'cirugia.editar',
    'catalogos.ver', 'notificaciones.ver'
  )
ON CONFLICT (rol_id, mod_id) DO NOTHING;

-- OftaData: Recepcionista
INSERT INTO "tblRoleModules_RMO" (rol_id, mod_id, created_at)
SELECT r.rol_id, m.mod_id, NOW()
FROM "tblRoles_ROL" r
JOIN "tblApplications_APP" a ON a.app_id = r.app_id
JOIN "tblModules_MOD" m ON m.app_id = a.app_id
WHERE a.app_code = 'OFTADATA' AND r.rol_name = 'Recepcionista'
  AND m.mod_name IN (
    'pacientes.ver', 'pacientes.crear',
    'citas.ver', 'citas.crear',
    'notificaciones.ver'
  )
ON CONFLICT (rol_id, mod_id) DO NOTHING;

-- VetData: Admin → todos los módulos
INSERT INTO "tblRoleModules_RMO" (rol_id, mod_id, created_at)
SELECT r.rol_id, m.mod_id, NOW()
FROM "tblRoles_ROL" r
JOIN "tblApplications_APP" a ON a.app_id = r.app_id
JOIN "tblModules_MOD" m ON m.app_id = a.app_id
WHERE a.app_code = 'VETDATA' AND r.rol_name = 'Admin'
  AND m.mod_status = 'active'
ON CONFLICT (rol_id, mod_id) DO NOTHING;

-- VetData: Veterinario
INSERT INTO "tblRoleModules_RMO" (rol_id, mod_id, created_at)
SELECT r.rol_id, m.mod_id, NOW()
FROM "tblRoles_ROL" r
JOIN "tblApplications_APP" a ON a.app_id = r.app_id
JOIN "tblModules_MOD" m ON m.app_id = a.app_id
WHERE a.app_code = 'VETDATA' AND r.rol_name = 'Veterinario'
  AND m.mod_name IN (
    'clientes.ver', 'clientes.crear',
    'mascotas.ver', 'mascotas.crear',
    'vacunas.ver', 'vacunas.crear', 'vacunas.aplicar',
    'medicamentos.ver', 'medicamentos.crear',
    'tickets.ver', 'tickets.responder',
    'consultas.ver', 'consultas.crear',
    'cirugias.ver', 'cirugias.crear',
    'libro.ver', 'libro.crear',
    'agenda.ver', 'agenda.crear'
  )
ON CONFLICT (rol_id, mod_id) DO NOTHING;

-- VetData: Asistente
INSERT INTO "tblRoleModules_RMO" (rol_id, mod_id, created_at)
SELECT r.rol_id, m.mod_id, NOW()
FROM "tblRoles_ROL" r
JOIN "tblApplications_APP" a ON a.app_id = r.app_id
JOIN "tblModules_MOD" m ON m.app_id = a.app_id
WHERE a.app_code = 'VETDATA' AND r.rol_name = 'Asistente'
  AND m.mod_name IN (
    'clientes.ver', 'mascotas.ver',
    'vacunas.ver', 'vacunas.aplicar',
    'medicamentos.ver', 'medicamentos.crear',
    'consultas.ver', 'cirugias.ver', 'libro.ver', 'agenda.ver'
  )
ON CONFLICT (rol_id, mod_id) DO NOTHING;

-- VetData: Recepcionista
INSERT INTO "tblRoleModules_RMO" (rol_id, mod_id, created_at)
SELECT r.rol_id, m.mod_id, NOW()
FROM "tblRoles_ROL" r
JOIN "tblApplications_APP" a ON a.app_id = r.app_id
JOIN "tblModules_MOD" m ON m.app_id = a.app_id
WHERE a.app_code = 'VETDATA' AND r.rol_name = 'Recepcionista'
  AND m.mod_name IN (
    'clientes.ver', 'clientes.crear',
    'mascotas.ver', 'vacunas.ver',
    'tickets.ver', 'tickets.responder',
    'inventario.ver', 'agenda.ver', 'agenda.crear'
  )
ON CONFLICT (rol_id, mod_id) DO NOTHING;

-- VetData: Cliente (read-only personal)
INSERT INTO "tblRoleModules_RMO" (rol_id, mod_id, created_at)
SELECT r.rol_id, m.mod_id, NOW()
FROM "tblRoles_ROL" r
JOIN "tblApplications_APP" a ON a.app_id = r.app_id
JOIN "tblModules_MOD" m ON m.app_id = a.app_id
WHERE a.app_code = 'VETDATA' AND r.rol_name = 'Cliente'
  AND m.mod_name IN ('mascotas.ver', 'vacunas.ver', 'medicamentos.ver', 'libro.ver')
ON CONFLICT (rol_id, mod_id) DO NOTHING;

-- =============================================================
-- SECTION 10: SEED DATA — ARTICDEV COMPANY + ADMIN USER SETUP
-- Note: The super-admin user is created by Go's SeedSuperAdmin()
-- at runtime using SUPER_ADMIN_USER env var. This section only
-- ensures the ArticDev company and branch exist. The role
-- assignment is done after the user exists.
-- =============================================================

INSERT INTO "tblCompanies_COM" (com_name, com_status, created_at, updated_at)
VALUES ('ArticDev', 'active', NOW(), NOW())
ON CONFLICT DO NOTHING;

INSERT INTO "tblBranches_BRA" (com_id, bra_code, bra_name, bra_address, bra_status, created_at, updated_at)
SELECT c.com_id, 'AD-001', 'Casa Matriz', 'San José, Costa Rica', 'active', NOW(), NOW()
FROM "tblCompanies_COM" c WHERE c.com_name = 'ArticDev'
ON CONFLICT (com_id, bra_code) DO NOTHING;

COMMIT;
