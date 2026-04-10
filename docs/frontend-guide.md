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

El frontend usa archivos `.env` de Vite para configuración por entorno. Hay dos archivos versionados y uno local:

| Archivo | Propósito | Versionado |
|---|---|---|
| `.env.development` | Valores para `npm run dev` | Sí |
| `.env.production` | Valores para `npm run build` | Sí |
| `.env` | Overrides locales (nunca commitear) | No (gitignored) |

> Vite carga automáticamente el archivo según el modo (`development` o `production`).
> Un archivo `.env` local siempre tiene prioridad sobre los demás.

Variables disponibles (solo las que comienzan con `VITE_` llegan al navegador):

| Variable | Descripción | Default |
|---|---|---|
| `VITE_API_BASE_URL` | URL base de la API | `/api/v1` |
| `VITE_APP_TITLE` | Título de la aplicación | `ArticNexus` |
| `VITE_VETDATA_URL` | URL pública de VetData (links en la landing) | `https://vetdata.articdev.com` |
| `VITE_OFTADATA_URL` | URL pública de OftaData (links en la landing) | `https://oftadata.articdev.com` |

> **Nota:** En desarrollo, el proxy de Vite (`vite.config.js`) redirige `/api` al backend local.
> En producción, nginx hace lo mismo. Por eso `VITE_API_BASE_URL` es un path relativo (`/api/v1`) en ambos entornos.
> **Nunca tocar `vite.config.js` en el servidor** — toda configuración específica del entorno va en `.env`.

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

El archivo `vite.config.js` configura un proxy para que las peticiones a `/api` se redirijan al backend:

```js
proxy: {
  '/api': {
    target: 'http://localhost:8080',
    changeOrigin: true
  }
}
```

Este proxy **solo aplica en desarrollo** (`npm run dev`). En producción, nginx se encarga del mismo reenvío.

> **No** editar `vite.config.js` para cambiar URLs o puertos. Si se necesita apuntar a otro backend local, crear un archivo `frontend/.env` (gitignored) con:
> ```
> VITE_API_BASE_URL=http://otro-host:9090/api/v1
> ```

## 9. Lint

```bash
npm run lint
```
