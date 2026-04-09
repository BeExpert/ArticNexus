-- =============================================================
-- Seed: VetData modules with display names
-- File: database/seeds/seed_vetdata_modules.sql
-- Usage: psql $DATABASE_URL -f seed_vetdata_modules.sql
--
-- Inserts all VetData modules with human-readable display
-- names. Safe to run multiple times (upsert on app_id + mod_name).
-- Requires migration 020 (UNIQUE constraint on app_id, mod_name).
-- Covers modules from migrations 014 and 016.
-- =============================================================

WITH app AS (
    SELECT app_id FROM "tblApplications_APP" WHERE app_code = 'VETDATA'
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
    -- ── Clientes ────────────────────────────────────────────────────
    ('clientes.ver',   'Ver clientes',       'Puede ver el listado de clientes'),
    ('clientes.crear', 'Registrar clientes', 'Puede registrar nuevos clientes'),
    -- ── Mascotas ────────────────────────────────────────────────────
    ('mascotas.ver',   'Ver mascotas',       'Puede ver el listado de mascotas'),
    ('mascotas.crear', 'Registrar mascotas', 'Puede registrar nuevas mascotas'),
    -- ── Vacunas ─────────────────────────────────────────────────────
    ('vacunas.ver',     'Ver vacunas',       'Puede ver el historial de vacunación'),
    ('vacunas.crear',   'Registrar vacunas', 'Puede registrar nuevas vacunas'),
    ('vacunas.aplicar', 'Aplicar vacunas',   'Puede registrar la aplicación de vacunas'),
    -- ── Medicamentos ────────────────────────────────────────────────
    ('medicamentos.ver',   'Ver medicamentos',       'Puede ver el catálogo de medicamentos'),
    ('medicamentos.crear', 'Registrar medicamentos', 'Puede registrar nuevos medicamentos'),
    -- ── Tickets ─────────────────────────────────────────────────────
    ('tickets.ver',       'Ver tickets',      'Puede ver los tickets de soporte'),
    ('tickets.responder', 'Responder tickets', 'Puede responder y gestionar tickets'),
    -- ── Inventario ──────────────────────────────────────────────────
    ('inventario.ver',   'Ver inventario',       'Puede ver el inventario de productos'),
    ('inventario.crear', 'Gestionar inventario', 'Puede agregar y modificar el inventario'),
    -- ── Consultas clínicas ──────────────────────────────────────────
    ('consultas.ver',   'Ver consultas',       'Puede ver el historial de consultas clínicas'),
    ('consultas.crear', 'Registrar consultas', 'Puede crear y editar consultas clínicas'),
    -- ── Cirugías ────────────────────────────────────────────────────
    ('cirugias.ver',   'Ver cirugías',       'Puede ver los registros de cirugías de mascotas'),
    ('cirugias.crear', 'Registrar cirugías', 'Puede registrar y editar cirugías'),
    -- ── Libro de Mascota ────────────────────────────────────────────
    ('libro.ver',   'Ver libro de mascota',    'Puede ver el libro de mascota completo'),
    ('libro.crear', 'Editar libro de mascota', 'Puede crear registros en el libro de mascota'),
    -- ── Agenda ──────────────────────────────────────────────────────
    ('agenda.ver',   'Ver agenda',    'Puede ver la agenda, citas y calendario'),
    ('agenda.crear', 'Gestionar agenda', 'Puede crear y gestionar citas y disponibilidad'),
    -- ── Administración ──────────────────────────────────────────────
    ('admin.usuarios', 'Administrar usuarios', 'Puede gestionar usuarios de la aplicación')
) AS v(mod_name, display_name, description)
ON CONFLICT (app_id, mod_name) DO UPDATE SET
    mod_display_name = EXCLUDED.mod_display_name,
    mod_description  = EXCLUDED.mod_description,
    updated_at       = now();
