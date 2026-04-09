# ArticNexus

Servicio centralizado de autenticación, usuarios, roles, permisos y aplicaciones del ecosistema ArticDev.

## Stack

| Capa | Tecnología |
|---|---|
| Backend | Go 1.21, Chi v5, GORM, PostgreSQL |
| Frontend | Vue 3, Pinia, Vite, Tailwind CSS |
| Base de datos | PostgreSQL — migración única (`001_full_schema.sql`) |
| Logging | Lumberjack — 3 canales: app, security, db |

## Inicio rápido

```bash
# 1. Backend
cd backend
cp .env.example .env   # editar con DB_URL, JWT_SECRET, etc.
go mod download
go run cmd/api/main.go -migrate   # solo la primera vez
go run cmd/api/main.go

# 2. Frontend (en otra terminal)
cd frontend
npm install
npm run dev
```

Servidor: `http://localhost:8080` · Frontend: `http://localhost:5173`

## Documentación

- [docs/backend-guide.md](docs/backend-guide.md)
- [docs/frontend-guide.md](docs/frontend-guide.md)
- [docs/database-guide.md](docs/database-guide.md)
- [docs/env-vars.md](docs/env-vars.md)
- [docs/smtp-config.md](docs/smtp-config.md)