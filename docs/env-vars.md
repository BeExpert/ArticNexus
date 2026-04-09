# Variables de entorno — ArticNexus Backend

Referencia completa de todas las variables de entorno soportadas por `backend/internal/config/config.go`.

> El archivo `.env.example` en `backend/` puede estar desactualizado. Este archivo es la fuente de verdad.

## Servidor HTTP

| Variable | Tipo | Requerido | Default | Descripción |
|---|---|---|---|---|
| `APP_PORT` / `PORT` | string | No | `8080` | Puerto en que escucha el servidor |
| `APP_ENV` / `ENVIRONMENT` | string | No | `development` | Entorno (`development` o `production`) |

En `development`, GORM registra queries con nivel `Warn`.

## Base de datos

| Variable | Tipo | Requerido | Default | Descripción |
|---|---|---|---|---|
| `DATABASE_URL` | string | Sí* | — | DSN completo de PostgreSQL |
| `DB_HOST` | string | Sí* | `localhost` | Host de PostgreSQL |
| `DB_PORT` | string | Sí* | `5432` | Puerto de PostgreSQL |
| `DB_USER` | string | Sí* | — | Usuario de PostgreSQL |
| `DB_PASSWORD` | string | Sí* | — | Contraseña de PostgreSQL |
| `DB_NAME` | string | Sí* | — | Nombre de la base de datos |
| `DB_SSL_MODE` | string | No | `disable` | `disable`, `require`, `verify-full` |
| `DB_MAX_OPEN_CONNS` | int | No | `25` | Máximo de conexiones abiertas al pool |
| `DB_MAX_IDLE_CONNS` | int | No | `10` | Conexiones idle en el pool |
| `DB_CONN_MAX_LIFETIME_MIN` | int | No | `60` | Tiempo de vida máximo de conexión (minutos) |

*`DATABASE_URL` tiene prioridad. Si no está, se construye el DSN desde `DB_HOST/PORT/USER/PASSWORD/NAME/SSL_MODE`.

**Ejemplo:**
```env
DATABASE_URL=postgres://postgres:secreto@localhost:5432/articnexus?sslmode=disable
```

## JWT y sesión

| Variable | Tipo | Requerido | Default | Descripción |
|---|---|---|---|---|
| `JWT_SECRET` | string | **Sí** | — | Clave secreta para firmar tokens. Mínimo 32 caracteres aleatorios. |
| `JWT_EXP_HOURS` | int | No | `24` | Tiempo de expiración del JWT en horas |

> **Epoch de sesión:** al arrancar, el backend genera un nonce aleatorio (`SessionEpoch`) que se embebe en cada JWT. Al reiniciar el backend, el epoch cambia y todos los tokens anteriores quedan inválidos, forzando re-login. Este mecanismo es completamente automático y no requiere configuración.

## CORS

| Variable | Tipo | Requerido | Default | Descripción |
|---|---|---|---|---|
| `ALLOWED_ORIGINS` | string (coma-separado) | No | `http://localhost:5173` | Orígenes permitidos para CORS |

**Ejemplo multi-origen:**
```env
ALLOWED_ORIGINS=https://nexus.articdev.com,https://admin.articdev.com
```

## Super-admin bootstrap

| Variable | Tipo | Requerido | Default | Descripción |
|---|---|---|---|---|
| `SUPER_ADMIN_USER` | string | Recomendado | `""` | Username del super-administrador |
| `SUPER_ADMIN_PASS` | string | Recomendado | `""` | Contraseña en texto plano (se hashea con bcrypt) |
| `SUPER_ADMIN_FORCE` | string | No | `0` | Si `1`, actualiza la contraseña del super-admin en cada boot |

Si no se configuran, el servidor arranca igualmente pero no habrá un usuario con acceso total.

## Migraciones

| Variable | Tipo | Requerido | Default | Descripción |
|---|---|---|---|---|
| `MIGRATIONS_DIR` | string | No | `../database/migrations` | Ruta al directorio de migraciones (relativa al binario) |

Solo relevante cuando se usa la flag `-migrate` al arrancar.

## Seguridad de contraseñas

| Variable | Tipo | Requerido | Default | Descripción |
|---|---|---|---|---|
| `BCRYPT_COST` | int | No | `12` | Factor de costo bcrypt (10-14 recomendado; mayor = más lento) |

## SMTP — Soporte técnico

| Variable | Tipo | Requerido | Default | Descripción |
|---|---|---|---|---|
| `SMTP_HOST` | string | No | `smtp.gmail.com` | Servidor SMTP |
| `SMTP_PORT` | string | No | `587` | Puerto SMTP (587 = TLS, 465 = SSL) |
| `SUPPORT_SMTP_USER` | string | No | `""` | Email de la cuenta de soporte |
| `SUPPORT_SMTP_PASSWORD` | string | No | `""` | App Password de Gmail (sin espacios) |
| `SUPPORT_SMTP_FROM` | string | No | `""` | Nombre y email de remitente. Ej: `ArticNexus <soporte@articdev.com>` |

Usada para: restablecimiento de contraseña, notificaciones técnicas.

## SMTP — Empresa / Negocio

| Variable | Tipo | Requerido | Default | Descripción |
|---|---|---|---|---|
| `BUSINESS_SMTP_USER` | string | No | `""` | Email de la cuenta comercial |
| `BUSINESS_SMTP_PASSWORD` | string | No | `""` | App Password de Gmail |
| `BUSINESS_SMTP_FROM` | string | No | `""` | Nombre y email de remitente. Ej: `ArticDev S.A. <info@articdev.com>` |
| `CONTACT_EMAIL` | string | No | `""` | Destino de los mensajes del formulario de contacto |

Usada para: formulario de contacto en la landing, propuestas comerciales.

Ver [smtp-config.md](smtp-config.md) para la guía completa.

## Frontend

| Variable | Tipo | Requerido | Default | Descripción |
|---|---|---|---|---|
| `FRONTEND_URL` | string | No | `http://localhost:5173` | URL base del frontend (para construir links de reset de contraseña en correos) |

## Tokens de reset de contraseña

| Variable | Tipo | Requerido | Default | Descripción |
|---|---|---|---|---|
| `PASSWORD_RESET_EXP_MIN` | int | No | `30` | Minutos de validez del token de restablecimiento de contraseña |

## Logging

| Variable | Tipo | Requerido | Default | Descripción |
|---|---|---|---|---|
| `LOG_MAX_SIZE_MB` | int | No | `10` | Tamaño máximo de cada archivo de log antes de rotar (MB) |
| `LOG_MAX_BACKUPS` | int | No | `3` | Número de archivos de log rotados a conservar |

Los logs se escriben en `backend/storage/logs/`:
- `app.log` — requests, seeds, ciclo de vida del servidor
- `security.log` — autenticación, autorización, tokens
- `db.log` — conectividad a la BD

---

## Ejemplo de .env completo (desarrollo)

```env
APP_PORT=8080
APP_ENV=development

DATABASE_URL=postgres://postgres:password@localhost:5432/articnexus?sslmode=disable

JWT_SECRET=dev-secret-change-in-production-minimum-32-chars
JWT_EXP_HOURS=24

ALLOWED_ORIGINS=http://localhost:5173

SUPER_ADMIN_USER=admin
SUPER_ADMIN_PASS=Admin123!
SUPER_ADMIN_FORCE=0

BCRYPT_COST=10

SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SUPPORT_SMTP_USER=soporte@articdev.com
SUPPORT_SMTP_PASSWORD=apppasswordhere
SUPPORT_SMTP_FROM=ArticNexus <soporte@articdev.com>
BUSINESS_SMTP_USER=info@articdev.com
BUSINESS_SMTP_PASSWORD=apppasswordhere
BUSINESS_SMTP_FROM=ArticDev S.A. <info@articdev.com>
CONTACT_EMAIL=info@articdev.com

FRONTEND_URL=http://localhost:5173
PASSWORD_RESET_EXP_MIN=30
LOG_MAX_SIZE_MB=10
LOG_MAX_BACKUPS=3
```
