# Guía de instalación — Backend ArticNexus

## Prerrequisitos

| Herramienta | Versión mínima | Notas |
|---|---|---|
| Go | 1.21 | [golang.org/dl](https://golang.org/dl/) — verificar con `go version` |
| PostgreSQL | 14+ | Puede ser local, Docker, o Neon/Railway/Supabase en prod |
| Git | cualquier | Para clonar el repositorio |

## 1. Clonar y preparar el entorno

```bash
git clone <repo-url>
cd NexusOftaData/ArticNexus/backend
cp .env.example .env
```

Editar `.env` con los valores reales (ver [env-vars.md](env-vars.md) para la referencia completa).

## 2. Variables de entorno mínimas

Las siguientes variables son **obligatorias** — el backend no arrancará sin ellas:

```env
DATABASE_URL=postgres://user:pass@localhost:5432/articnexus?sslmode=disable
JWT_SECRET=cadena-aleatoria-de-al-menos-32-caracteres
SUPER_ADMIN_USER=admin
SUPER_ADMIN_PASS=contraseña-segura
```

Para la guía completa de todas las variables ver [env-vars.md](env-vars.md).

## 3. Descargar dependencias

```bash
go mod download
go mod verify
```

`go mod verify` comprueba que los hashes de los módulos descargados coinciden exactamente con los registrados en `go.sum`. Si falla, hay una corrupción o manipulación del caché.

## 4. Primera migración (base de datos vacía)

```bash
go run cmd/api/main.go -migrate
```

Este comando ejecuta `database/migrations/001_full_schema.sql` que crea las 14 tablas, triggers, índices y carga todos los datos semilla (3 aplicaciones, ~85 módulos, roles base, empresa ArticDev).

> **Nota:** La flag `-migrate` solo debe usarse en la primera instalación o al realizar migraciones intencionadas. En arranques normales se omite.

## 5. Arranque normal

```bash
go run cmd/api/main.go
```

Al arrancar, el servidor ejecuta automáticamente:
1. Inicialización del logger (archivos en `storage/logs/`)
2. Generación de `SessionEpoch` aleatorio (invalida sesiones anteriores)
3. Conexión a PostgreSQL
4. Upsert del super-admin (si `SUPER_ADMIN_USER` está configurado)
5. Upsert de módulos de ArticNexus
6. Upsert de módulos de OftaData y VetData
7. Inicio del servidor HTTP en el puerto configurado (default: `8080`)

## 6. Compilar binario de producción

```bash
go build -o api cmd/api/main.go
./api
```

## 7. Dependencias exactas

### Dependencias directas (go.mod)

| Módulo | Versión | Propósito |
|---|---|---|
| `github.com/go-chi/chi/v5` | `v5.1.0` | Router HTTP |
| `github.com/golang-jwt/jwt/v5` | `v5.2.1` | JWT — autenticación stateless con epoch de sesión |
| `github.com/joho/godotenv` | `v1.5.1` | Carga de archivos `.env` |
| `golang.org/x/crypto` | `v0.31.0` | bcrypt para hash de contraseñas |
| `gorm.io/driver/postgres` | `v1.5.9` | Driver PostgreSQL para GORM |
| `gorm.io/gorm` | `v1.25.12` | ORM — mapeo de entidades y queries |
| `gopkg.in/natefinch/lumberjack.v2` | `v2.2.1` | Rotación de archivos de log |

### Dependencias indirectas (go.mod)

| Módulo | Versión |
|---|---|
| `github.com/jackc/pgx/v5` | `v5.5.5` |
| `github.com/jackc/puddle/v2` | `v2.2.1` |
| `github.com/jackc/pgpassfile` | `v1.0.0` |
| `github.com/jackc/pgservicefile` | `v0.0.0-20221227161230-091c0ba34f0a` |
| `github.com/jinzhu/inflection` | `v1.0.0` |
| `github.com/jinzhu/now` | `v1.1.5` |
| `golang.org/x/sync` | `v0.10.0` |
| `golang.org/x/text` | `v0.21.0` |

### Verificación de integridad

El archivo `go.sum` contiene los hashes `h1:` de cada módulo. Para verificar que las copias locales son íntegras:

```bash
go mod verify
# output esperado: all modules verified
```

Hashes de módulos directos:

| Módulo | Hash h1: |
|---|---|
| `go-chi/chi/v5 v5.1.0` | `h1:acVI1TYaD+hhedDJ3r54HyA6sExp3HfXq7QWEEY/xMw=` |
| `golang-jwt/jwt/v5 v5.2.1` | `h1:OuVbFODueb089Lh128TAcimifWaLhJwVflnrgM17wHk=` |
| `joho/godotenv v1.5.1` | `h1:7eLL/+HRGLY0ldzfGMeQkb7vMd0as4CfYvUVzLqw0N0=` |
| `x/crypto v0.31.0` | `h1:ihbySMvVjLAeSH1IbfcRTkD/iNscyz8rGzjF/E5hV6U=` |
| `gorm.io/driver/postgres v1.5.9` | `h1:DkegyItji119OlcaLjqN11kHoUgZ/j13E0jkJZgD6A8=` |
| `gorm.io/gorm v1.25.12` | `h1:I0u8i2hWQItBq1WfE0o2+WuL9+8L21K9e2HHSTE/0f8=` |

## 8. Sistema de logs

Al arrancar el backend crea tres archivos de log en `storage/logs/`:

| Archivo | Contenido |
|---|---|
| `app.log` | Requests HTTP, seeds, arranque, shutdown |
| `security.log` | Login OK/fallido, tokens inválidos, accesos denegados |
| `db.log` | Errores de conexión y queries a la BD |

Cada archivo rota automáticamente al superar `LOG_MAX_SIZE_MB` (default 10 MB) y conserva `LOG_MAX_BACKUPS` copias (default 3). Todos los logs también se emiten a stderr.

## 9. Verificación post-instalación

```bash
# Verificar que el servidor responde
curl http://localhost:8080/health
# Respuesta esperada: {"status":"ok"}

# Verificar las tres aplicaciones en la BD
psql $DATABASE_URL -c 'SELECT app_code, app_name FROM "tblApplications_APP";'
# Esperado: ARTICNEXUS, OFTADATA, VETDATA
```

## 10. SMTP — correo electrónico

Ver [smtp-config.md](smtp-config.md) para la guía completa de configuración del SMTP.
