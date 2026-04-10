# Guía de base de datos — ArticNexus

## Estructura de migraciones

ArticNexus usa **un solo archivo de migración** como fuente de verdad:

```
database/
├── migrations/
│   ├── 001_full_schema.up.sql     ← DDL + seeds completos (ejecutado por -migrate)
│   └── 001_full_schema.down.sql   ← rollback completo (solo uso manual)
├── schema/
│   ├── articnexus_dbdiagram.dbml  ← schema V1 (referencia histórica)
│   └── schemaV2.dbml              ← schema V2 actualizado (abrir en dbdiagram.io)
└── seeds/
    ├── seed_articnexus_modules.sql ← referencia de módulos ARTICNEXUS
    ├── seed_oftadata_modules.sql   ← referencia de módulos OFTADATA
    └── seed_vetdata_modules.sql    ← referencia de módulos VETDATA
```

> Los archivos en `seeds/` son solo de referencia. El seeding real ocurre en Go (`db.SeedModules()`, `db.SeedApplications()`, `db.SeedArticDevAndDemoUsers()`) ejecutado automáticamente en cada arranque del backend.

> Los archivos `.down.sql` nunca se ejecutan automáticamente. Son para rollback manual en caso de emergencia.

## Reiniciar la BD desde cero

Para borrar todo y empezar de nuevo con las migraciones y seeds:

```bash
cd ArticNexus/backend
./api -reset-db
```

Esto ejecuta en orden:
1. `DROP SCHEMA public CASCADE` + `CREATE SCHEMA public` (borra TODO)
2. Ejecuta `001_full_schema.up.sql` (crea tablas, índices, triggers, seeds)
3. Ejecuta los seeders de Go (super-admin, módulos, apps, ArticDev + demos)

> **Solo funciona en `APP_ENV=development`**. En producción el flag es rechazado.

Alternativamente, para solo ejecutar las migraciones sin borrar:

```bash
./api -migrate
```

O manualmente con psql:

```bash
# Reset completo manual
PGPASSWORD=postgres psql -h localhost -U postgres -d Nexus -f database/migrations/001_full_schema.down.sql
PGPASSWORD=postgres psql -h localhost -U postgres -d Nexus -f database/migrations/001_full_schema.up.sql
```

## Arranque normal (sin flags)

Al arrancar sin flags, el backend ejecuta automáticamente estos seeders idempotentes:

1. `SeedSuperAdmin()` — crea o actualiza el usuario super-admin
2. `SeedModules()` — upsert de ~41 módulos de ARTICNEXUS
3. `SeedApplications()` — upsert de las 3 apps + módulos de OftaData y VetData
4. `SeedArticDevAndDemoUsers()` — upsert de ArticDev company, branch, licencias, demo users

Todos usan `ON CONFLICT ... DO UPDATE` o `DO NOTHING`. Nunca crean duplicados.

## Agregar nuevas migraciones en el futuro

Cuando el esquema cambie (nueva columna, nueva tabla, etc.):

1. **Actualizar** `001_full_schema.up.sql` y `.down.sql` con el DDL actualizado
2. **Crear** un archivo incremental `002_<descripcion>.up.sql` con solo el `ALTER TABLE` o DDL nuevo
3. Ejecutar `./api -migrate`

> En producción, nunca ejecutar `-migrate` automáticamente. Revisar el archivo antes.

## Tablas del sistema

| Tabla | Propósito |
|---|---|
| `tblPersons_PER` | Personas físicas (datos del individuo) |
| `tblUsers_USR` | Cuentas de usuario (autenticación) |
| `tblCompanies_COM` | Empresas cliente del ecosistema |
| `tblBranches_BRA` | Sucursales de cada empresa |
| `tblUserCompanies_UCO` | Asignación usuario ↔ empresa |
| `tblUserBranches_UBR` | Asignación usuario ↔ sucursal |
| `tblApplications_APP` | Apps registradas (ARTICNEXUS, OFTADATA, VETDATA) |
| `tblCompanyApplications_CAP` | Licencias: qué apps puede usar cada empresa |
| `tblRoles_ROL` | Roles por aplicación |
| `tblModules_MOD` | Permisos granulares (recurso.acción) |
| `tblRoleModules_RMO` | Permisos que tiene cada rol |
| `tblUserRoles_URO` | Rol de un usuario en empresa+sucursal |
| `tblPasswordResetTokens_PRT` | Tokens de restablecimiento de contraseña |
| `tblDemoLinks_DML` | Links de acceso demo temporales |

## Convenciones de nomenclatura

- Tablas: `tbl<NombrePlural>_<SIGLA>` — ej. `tblUsers_USR`
- Columnas: `<sigla>_<camelCase>` — ej. `usr_username`, `per_firstName`
- Módulos: `recurso.acción` — ej. `empresas.ver`, `roles.asignar_modulo`
- Super-admin: identificado **solo** por variable de entorno `SUPER_ADMIN_USER`, nunca por columna en BD

## Soft-delete

Las tablas de entidad tienen columna `deleted_at` (GORM convention). GORM filtra automáticamente registros con `deleted_at IS NOT NULL` en todas las queries.

## Visualizar el schema

Abrir [dbdiagram.io](https://dbdiagram.io) e importar el archivo `database/schema/schemaV2.dbml`.
