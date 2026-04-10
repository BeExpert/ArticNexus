-- =============================================================
-- Seed: ArticNexus modules with display names
-- File: database/seeds/seed_articnexus_modules.sql
-- Usage: psql $DATABASE_URL -f seed_articnexus_modules.sql
--
-- Inserts all ArticNexus modules with human-readable display
-- names. Safe to run multiple times (upsert on app_id + mod_name).
-- Requires migration 020 (UNIQUE constraint on app_id, mod_name).
-- =============================================================

WITH app AS (
    SELECT app_id FROM "tblApplications_APP" WHERE app_code = 'ARTICNEXUS'
)
INSERT INTO "tblModules_MOD"
    (app_id, mod_name, mod_display_name, mod_description, mod_status, created_at, updated_at)
SELECT
    app.app_id,
    v.mod_name,
    v.display_name,
    v.description,
    'active',
    now(),
    now()
FROM app, (VALUES
    -- ── Aplicaciones ────────────────────────────────────────────────
    ('aplicaciones.ver',                'Ver aplicaciones',                    'Puede ver el listado de aplicaciones registradas'),
    ('aplicaciones.crear',              'Crear aplicaciones',                  'Puede registrar nuevas aplicaciones'),
    ('aplicaciones.editar',             'Editar aplicaciones',                 'Puede modificar aplicaciones existentes'),
    ('aplicaciones.eliminar',           'Eliminar aplicaciones',               'Puede eliminar aplicaciones del sistema'),
    ('aplicaciones.asignar_usuario',    'Asignar usuarios a aplicaciones',     'Puede asignar usuarios a aplicaciones'),
    ('aplicaciones.desasignar_usuario', 'Desasignar usuarios de aplicaciones', 'Puede desasignar usuarios de aplicaciones'),
    -- ── Roles ───────────────────────────────────────────────────────
    ('roles.ver',                'Ver roles',                    'Puede ver los roles disponibles'),
    ('roles.crear',              'Crear roles',                  'Puede crear nuevos roles'),
    ('roles.editar',             'Editar roles',                 'Puede modificar roles existentes'),
    ('roles.eliminar',           'Eliminar roles',               'Puede eliminar roles del sistema'),
    ('roles.asignar_usuario',    'Asignar usuarios a roles',     'Puede asignar usuarios a roles'),
    ('roles.desasignar_usuario', 'Desasignar usuarios de roles', 'Puede desasignar usuarios de roles'),
    ('roles.asignar_modulo',     'Asignar módulos a roles',      'Puede asignar módulos a roles'),
    ('roles.desasignar_modulo',  'Desasignar módulos de roles',  'Puede desasignar módulos de roles'),
    -- ── Usuarios ────────────────────────────────────────────────────
    ('usuarios.ver',               'Ver usuarios',                       'Puede ver el listado de usuarios del sistema'),
    ('usuarios.crear',             'Crear usuarios',                     'Puede crear nuevas cuentas de usuario'),
    ('usuarios.editar',            'Editar usuarios',                    'Puede modificar cuentas de usuario'),
    ('usuarios.eliminar',          'Eliminar usuarios',                  'Puede eliminar cuentas de usuario'),
    ('usuarios.reset_contrasena',  'Resetear contraseña de usuario',     'Puede restablecer la contraseña de otros usuarios'),
    ('usuarios.demo',              'Gestionar usuarios demo',            'Puede configurar usuarios con acceso temporal (fecha de expiración)'),
    -- ── Módulos ─────────────────────────────────────────────────────
    ('modulos.ver',      'Ver módulos',      'Puede ver los módulos de una aplicación'),
    ('modulos.crear',    'Crear módulos',    'Puede crear nuevos módulos'),
    ('modulos.editar',   'Editar módulos',   'Puede modificar módulos existentes'),
    ('modulos.eliminar', 'Eliminar módulos', 'Puede eliminar módulos del sistema'),
    -- ── Empresas ────────────────────────────────────────────────────
    ('empresas.ver',      'Ver empresas',      'Puede ver el listado de empresas afiliadas'),
    ('empresas.crear',    'Registrar empresas', 'Puede registrar nuevas empresas en el sistema'),
    ('empresas.editar',   'Editar empresas',   'Puede modificar los datos de empresas existentes'),
    ('empresas.eliminar', 'Eliminar empresas', 'Puede eliminar empresas del sistema'),
    -- ── Personas ────────────────────────────────────────────────────
    ('personas.ver',         'Ver personas',            'Puede ver las personas asignadas a una empresa'),
    ('personas.crear',       'Agregar personas',        'Puede agregar personas a una empresa'),
    ('personas.editar',      'Editar personas',         'Puede modificar los datos de personas'),
    ('personas.eliminar',    'Quitar personas',         'Puede quitar personas de una empresa'),
    ('personas.asignar_rol', 'Asignar rol a persona',   'Puede asignar y desasignar roles a personas de una empresa')
) AS v(mod_name, display_name, description)
ON CONFLICT (app_id, mod_name) DO UPDATE SET
    mod_display_name = EXCLUDED.mod_display_name,
    mod_description  = EXCLUDED.mod_description,
    updated_at       = now();
