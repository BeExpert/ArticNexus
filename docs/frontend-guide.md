# Guía de instalación — Frontend ArticNexus

## Prerrequisitos

| Herramienta | Versión mínima | Notas |
|---|---|---|
| Node.js | 18+ | [nodejs.org](https://nodejs.org) — verificar con `node -v` |
| npm | 9+ | Incluido con Node.js |

## 1. Instalar dependencias

```bash
cd NexusOftaData/ArticNexus/frontend
npm install
```

## 2. Configurar variables de entorno

```bash
cp .env.example .env
```

Variables disponibles en el frontend (archivo `.env` en `frontend/`):

| Variable | Descripción | Ejemplo |
|---|---|---|
| `VITE_API_BASE_URL` | URL base de la API (solo para peticiones directas desde landing) | `http://localhost:8080` |
| `VITE_VETDATA_URL` | URL pública de VetData (links en la landing) | `https://vetdata.articdev.com` |
| `VITE_OFTADATA_URL` | URL pública de OftaData (links en la landing) | `https://oftadata.articdev.com` |

> **Nota:** Las peticiones a `/api/v1/*` van a través del proxy de Vite en desarrollo y **no** requieren variable de entorno. El proxy está configurado en `vite.config.js`.

## 3. Modo desarrollo

```bash
npm run dev
```

Inicia el servidor de desarrollo en `http://localhost:5173`. El proxy de Vite redirige automáticamente `/api/v1/*` al backend en `http://localhost:8080`.

Asegurarse de que el backend esté corriendo antes de abrir el frontend.

## 4. Compilar para producción

```bash
npm run build
```

Genera los archivos estáticos en `frontend/dist/`. Los archivos pueden servirse desde cualquier servidor web estático (nginx, Caddy, etc.) o plataforma de despliegue (Netlify, Vercel, etc.).

## 5. Previsualizar el build de producción

```bash
npm run preview
```

## 6. Estructura del frontend

```
frontend/src/
├── main.js              ← Punto de entrada (Pinia + Router + mount)
├── App.vue              ← Shell principal (nav + layout)
├── router/index.js      ← Rutas + guard de autenticación
├── store/
│   ├── auth.js          ← Estado de autenticación, permisos, epoch
│   └── users.js         ← CRUD de usuarios
├── services/
│   ├── api.js           ← Instancia Axios con interceptor JWT
│   ├── userService.js   ← authService + userService
│   ├── companyService.js← CRUD de empresas + getMyCompanies()
│   └── branchService.js ← CRUD de sucursales
└── views/               ← Una vista por ruta
```

## 7. Dependencias exactas (package.json)

### Dependencias de producción

| Paquete | Versión |
|---|---|
| `vue` | `^3.4.15` |
| `vue-router` | `^4.2.5` |
| `pinia` | `^2.1.7` |
| `axios` | `^1.6.7` |
| `lucide-vue-next` | `^0.344.0` |

### Dependencias de desarrollo

| Paquete | Versión |
|---|---|
| `vite` | `^5.0.12` |
| `@vitejs/plugin-vue` | `^5.0.3` |
| `tailwindcss` | `^3.4.1` |
| `autoprefixer` | `^10.4.17` |
| `postcss` | `^8.4.35` |
| `eslint` | `^8.49.0` |
| `eslint-plugin-vue` | `^9.17.0` |

## 8. Proxy Vite (desarrollo)

El archivo `vite.config.js` configura un proxy para que las peticiones a `/api/v1` apunten al backend:

```js
proxy: {
  '/api': {
    target: 'http://localhost:8080',
    changeOrigin: true
  }
}
```

En producción, el proxy es reemplazado por la configuración de nginx (ver [deploy.md](deploy.md) si existe).

## 9. Lint

```bash
npm run lint
```
