# ArticNexus — Documentación

> Documentación técnica del sistema ArticNexus

## Guías de instalación y configuración

- [backend-guide.md](backend-guide.md) — Instalación del backend Go, dependencias exactas, arranque
- [frontend-guide.md](frontend-guide.md) — Setup del frontend Vue 3, variables de entorno, build
- [database-guide.md](database-guide.md) — Migración única, estructura de tablas, convenciones
- [env-vars.md](env-vars.md) — Referencia completa de todas las variables de entorno
- [smtp-config.md](smtp-config.md) — Configuración de Gmail App Passwords

## Arquitectura y diseño

- [architecture.md](architecture.md) — Visión general de capas y patrones

## 🏗️ Arquitectura General

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│                 │    │                 │    │                 │
│   Frontend      │◄──►│    Backend      │◄──►│   Database      │
│   (Vue 3)       │    │    (Go/Chi)     │    │  (PostgreSQL)   │
│                 │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
        │                        │                        │
        │                        │                        │
        ▼                        ▼                        ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│                 │    │                 │    │                 │
│   Static Files  │    │   API Gateway   │    │   Migrations    │
│   (Vite Build)  │    │   Load Balancer │    │   Backups       │
│                 │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🎯 Funcionalidades Principales

### Gestión de Usuarios
- ✅ Registro y autenticación
- ✅ Perfiles de usuario
- ✅ Gestión de sesiones
- 🚧 Multi-factor authentication (planned)

### Gestión de Roles y Permisos
- ✅ Definición de roles
- ✅ Asignación de permisos
- ✅ Control de acceso basado en roles (RBAC)
- 🚧 Permisos granulares (planned)

### Multi-tenancy
- ✅ Gestión de empresas
- ✅ Aislamiento de datos por empresa
- ✅ Roles específicos por empresa
- 🚧 Configuración personalizada (planned)

### Gestión de Aplicaciones
- ✅ Registro de aplicaciones
- ✅ Módulos y funcionalidades
- ✅ Control de acceso por aplicación
- 🚧 SSO integration (planned)

## 🚀 Estado Actual del Proyecto

| Componente | Estado | Progreso |
|------------|--------|----------|
| Backend API | 🚧 Scaffolding | 25% |
| Frontend UI | 🚧 Scaffolding | 25% |
| Database Schema | ✅ Completo | 90% |
| Authentication | ❌ Pendiente | 0% |
| Authorization | ❌ Pendiente | 0% |
| Documentation | 🚧 En progreso | 50% |
| Testing | ❌ Pendiente | 0% |
| CI/CD | ❌ Pendiente | 0% |

## 🔄 Próximos Pasos

### Fase 1: Core Implementation
1. **Autenticación JWT** - Implementar login/logout completo
2. **CRUD de Usuarios** - Operaciones básicas de usuarios
3. **Gestión de Roles** - Asignación y validación de roles
4. **Frontend básico** - Interfaces para gestión de usuarios

### Fase 2: Advanced Features
1. **Autorización granular** - Permisos específicos por módulo
2. **Multi-tenancy** - Aislamiento completo por empresa
3. **API completa** - Todos los endpoints documentados
4. **Testing** - Cobertura de pruebas unitarias e integración

### Fase 3: Production Ready
1. **Security hardening** - Auditoría de seguridad
2. **Performance optimization** - Optimización de consultas y caching
3. **Monitoring** - Métricas y alertas
4. **Documentation** - Documentación completa

## 🤝 Contribuir

Para contribuir al proyecto:

1. **Fork** el repositorio
2. **Crear rama** para la funcionalidad (`git checkout -b feature/nueva-funcionalidad`)
3. **Commit** los cambios (`git commit -m 'Agregar nueva funcionalidad'`)
4. **Push** a la rama (`git push origin feature/nueva-funcionalidad`)
5. **Pull Request** para revisión

### Estándares de Desarrollo

- **Backend**: Seguir principios de arquitectura limpia
- **Frontend**: Utilizar Composition API de Vue 3
- **Database**: Migrations numeradas secuencialmente
- **Testing**: Cobertura mínima del 80%
- **Documentation**: Documentar todas las APIs y funcionalidades

## 📞 Soporte

- **Issues**: Reportar bugs en GitHub Issues
- **Discussions**: Preguntas y discusiones en GitHub Discussions
- **Documentation**: Consultar esta documentación

---

**Última actualización**: 3 de febrero de 2026  
**Versión**: 0.1.0 (Scaffolding inicial)