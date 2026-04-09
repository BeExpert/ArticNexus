# Guía de base de datos — ArticNexus

## Estructura de migraciones

A diferencia de un sistema incremental, ArticNexus usa **un solo archivo de migración** como fuente de verdad:

```
database/
├── migrations/
│   └── 001_full_schema.sql    ← única migración: DDL + seeds completos
├── schema/
│   ├── articnexus_dbdiagram.dbml  ← schema V1 (referencia histórica)
│   └── schemaV2.dbml              ← schema V2 actualizado (abrir en dbdiagram.io)
└── seeds/
    ├── seed_articnexus_modules.sql ← referencia de módulos ARTICNEXUS (no ejecutar manualmente)
    ├── seed_oftadata_modules.sql   ← referencia de módulos OFTADATA
    └── seed_vetdata_modules.sql    ← referencia de módulos VETDATA
```

> Los archivos en `seeds/` son de referencia. El seeding real ocurre en Go (`db.SeedModules()` y `db.SeedApplications()`) ejecutado automáticamente en cada arranque del backend.

## Primera migración (base de datos vacía)

```bash
cd ArticNexus/backend
go run cmd/api/main.go -migrate
```

Esto ejecuta `001_full_schema.sql` que construye:

1. **14 tablas** con sus columnas finales, constraints e índices
2. **Función y 7 triggers** `fn_set_updated_at()` para auto-actualizar `updated_at`
3. **3 aplicaciones**: ARTICNEXUS, OFTADATA, VETDATA
4. **~85 módulos** con `mod_display_name` para las 3 apps
5. **9 roles** distribuidos entre las 3 apps
6. **Asignaciones rol-módulo** para todos los roles
7. **Empresa ArticDev** con sucursal Casa Matriz (AD-001)

La migración usa `IF NOT EXISTS`, `ON CONFLICT DO NOTHING` y `ON CONFLICT DO UPDATE` en todos lados — es segura ejecutarla sobre una BD ya inicializada.

## Arranque normal (sin -migrate)

Al arrancar sin la flag `-migrate`, el backend ejecuta automáticamente estos seeders idempotentes:

1. `SeedSuperAdmin()` — crea o actualiza el usuario super-admin
2. `SeedModules()` — upsert de ~34 módulos de ARTICNEXUS
3. `SeedApplications()` — upsert de las 3 apps + módulos de OftaData y VetData

Estos seeders usan `ON CONFLICT ... DO UPDATE` o `DO NOTHING`, por lo que nunca crean duplicados ni fallan si los datos ya existen.

## Agregar nuevas migraciones en el futuro

Cuando el esquema necesite cambiar (nueva columna, nueva tabla, etc.):

1. Actualizar `001_full_schema.sql` con el nuevo DDL
2. Crear un archivo adicional `002_<descripcion>.sql` con solo el `ALTER TABLE` o DDL incremental
3. Ejecutar `go run cmd/api/main.go -migrate` en el entorno objetivo

> En producción, nunca ejecutar `-migrate` automáticamente en el arranque. Hacerlo manualmente luego de revisar el archivo de migración.

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

Las tablas de entidad tienen columna `deleted_at` (GORM convention). GORM filtra automáticamente registros con `deleted_at IS NOT NULL` en todas las queries. Los registros "eliminados" no se borran físicamente de la BD.

## Visualizar el schema

Abrir [dbdiagram.io](https://dbdiagram.io) e importar el archivo `database/schema/schemaV2.dbml`.
