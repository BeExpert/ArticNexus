-- =============================================================
-- Seed: OftaData modules with display names
-- File: database/seeds/seed_oftadata_modules.sql
-- Usage: psql $DATABASE_URL -f seed_oftadata_modules.sql
--
-- Inserts all OftaData modules with human-readable display
-- names. Safe to run multiple times (upsert on app_id + mod_name).
-- Requires migration 020 (UNIQUE constraint on app_id, mod_name).
-- =============================================================

WITH app AS (
    SELECT app_id FROM "tblApplications_APP" WHERE app_code = 'OftaData'
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
    -- ── Pacientes ───────────────────────────────────────────────────
    ('pacientes.ver',        'Ver pacientes',        'Puede ver el listado de pacientes'),
    ('pacientes.crear',      'Registrar pacientes',  'Puede registrar nuevos pacientes'),
    ('pacientes.editar',     'Editar pacientes',     'Puede modificar los datos de pacientes'),
    ('pacientes.eliminar',   'Eliminar pacientes',   'Puede eliminar pacientes del sistema'),
    ('pacientes.reactivar',  'Reactivar pacientes',  'Puede reactivar pacientes inactivos'),
    -- ── Formularios ─────────────────────────────────────────────────
    ('formularios.ver',       'Ver formularios',       'Puede ver los formularios clínicos'),
    ('formularios.crear',     'Crear formulario',      'Puede crear formularios clínicos'),
    ('formularios.editar',    'Editar formulario',     'Puede modificar formularios existentes'),
    ('formularios.eliminar',  'Eliminar formulario',   'Puede eliminar formularios del sistema'),
    ('formularios.reactivar', 'Reactivar formulario',  'Puede reactivar formularios inactivos'),
    -- ── Citas ───────────────────────────────────────────────────────
    ('citas.ver',      'Ver agenda',    'Puede ver las citas agendadas'),
    ('citas.crear',    'Crear cita',    'Puede crear nuevas citas'),
    ('citas.editar',   'Editar cita',   'Puede modificar citas existentes'),
    ('citas.eliminar', 'Cancelar cita', 'Puede cancelar o eliminar citas'),
    -- ── Cirugía ─────────────────────────────────────────────────────
    ('cirugia.ver',      'Ver cirugías',       'Puede ver los registros de cirugías'),
    ('cirugia.crear',    'Registrar cirugía',  'Puede registrar nuevas cirugías'),
    ('cirugia.editar',   'Editar cirugía',     'Puede modificar registros de cirugías'),
    ('cirugia.eliminar', 'Eliminar cirugía',   'Puede eliminar registros de cirugías'),
    -- ── Catálogos ───────────────────────────────────────────────────
    ('catalogos.ver',       'Ver catálogos',       'Puede ver los catálogos del sistema'),
    ('catalogos.gestionar', 'Gestionar catálogos', 'Puede crear y modificar catálogos'),
    -- ── Notificaciones ──────────────────────────────────────────────
    ('notificaciones.ver',       'Ver notificaciones',       'Puede ver las notificaciones del sistema'),
    ('notificaciones.gestionar', 'Gestionar notificaciones', 'Puede configurar y gestionar notificaciones'),
    -- ── Administración ──────────────────────────────────────────────
    ('admin.usuarios',  'Administrar usuarios',       'Puede gestionar usuarios de la aplicación'),
    ('admin.reset_db',  'Resetear base de datos',     'Puede reinicializar los datos de la base de datos')
) AS v(mod_name, display_name, description)
ON CONFLICT (app_id, mod_name) DO UPDATE SET
    mod_display_name = EXCLUDED.mod_display_name,
    mod_description  = EXCLUDED.mod_description,
    updated_at       = now();
