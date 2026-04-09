# Arquitectura del Sistema ArticNexus

## 🏗️ Visión General

ArticNexus está diseñado como un **servicio centralizado de autenticación y autorización** que sigue los principios de **arquitectura limpia** y **Domain-Driven Design (DDD)**.

## 🎯 Principios Arquitectónicos

### 1. Separación de Responsabilidades
- **Frontend**: Interfaz de usuario reactiva
- **Backend**: Lógica de negocio y API REST
- **Database**: Persistencia y integridad de datos

### 2. Arquitectura Limpia (Clean Architecture)
```
┌─────────────────────────────────────────────────┐
│                   External                      │
│  ┌─────────────────────────────────────────┐    │
│  │              Interface                  │    │
│  │  ┌─────────────────────────────────┐    │    │
│  │  │           Application           │    │    │
│  │  │  ┌─────────────────────────┐    │    │    │
│  │  │  │       Domain            │    │    │    │
│  │  │  │   ┌─────────────────┐   │    │    │    │
│  │  │  │   │   Entities      │   │    │    │    │
│  │  │  │   │   Rules         │   │    │    │    │
│  │  │  │   └─────────────────┘   │    │    │    │
│  │  │  └─────────────────────────┘    │    │    │
│  │  └─────────────────────────────────┘    │    │
│  └─────────────────────────────────────────┘    │
└─────────────────────────────────────────────────┘

External: Database, Web Frameworks, UI
Interface: Controllers, Gateways, Presenters  
Application: Use Cases, Services
Domain: Entities, Business Rules
```

### 3. Dependency Inversion
- Las dependencias fluyen hacia adentro
- El dominio no depende de frameworks externos
- Interfaces definen contratos entre capas

## 🏢 Arquitectura de Componentes

### Backend (Go)
```
cmd/
└── server/           # Entry point
    └── main.go

internal/
├── config/          # Configuration management
├── db/              # Database connection
├── domain/          # Business entities
├── repository/      # Data access layer
├── service/         # Business logic layer
├── handler/         # HTTP request handlers
└── middleware/      # HTTP middleware
```

#### Flujo de Datos Backend
```
HTTP Request → Middleware → Handler → Service → Repository → Database
                  ↑           ↓         ↓          ↓
              CORS/Auth    Validation  Business   Data
              Logging      Transform   Logic      Access
```

### Frontend (Vue 3)
```
src/
├── assets/          # Static assets
├── components/      # Reusable components  
├── views/          # Page components
├── router/         # Navigation routing
├── services/       # API communication
├── store/          # State management (Pinia)
└── main.js         # Application entry
```

#### Flujo de Datos Frontend
```
User Interaction → Vue Component → Pinia Store → API Service → Backend
       ↑              ↓              ↓            ↓
   UI Updates     State Change   HTTP Request   Business Logic
```

## 🔐 Modelo de Seguridad

### Autenticación
- **JWT (JSON Web Tokens)** para sesiones stateless
- **Refresh tokens** para renovación automática
- **Password hashing** con bcrypt

### Autorización
- **RBAC (Role-Based Access Control)**
- **Multi-tenancy** a nivel de empresa
- **Permisos granulares** por módulo/función

```
User → Role → Permissions → Resources
  ↓      ↓         ↓           ↓
Juan → Admin → [CRUD] → [Users, Roles]
```

### Flujo de Autenticación
```
1. Usuario envía credenciales
2. Backend valida con database
3. Genera JWT con claims
4. Frontend almacena token
5. Requests incluyen Authorization header
6. Backend valida JWT en cada request
```

## 🗄️ Modelo de Datos

### Entidades Principales
```
Company 1---* User *---* Role
   ↑             ↓       ↓
   └── Application 1---* Module
```

### Relaciones Clave
- **User belongs to Company** (opcional para usuarios globales)
- **Role belongs to Company** (opcional para roles globales)
- **User has many Roles** (many-to-many)
- **Application belongs to Company** (opcional para apps globales)
- **Module belongs to Application** (one-to-many)

## 🌐 Patrones de API

### REST Endpoints
```
GET    /api/v1/users           # List users
GET    /api/v1/users/:id       # Get user
POST   /api/v1/users           # Create user  
PUT    /api/v1/users/:id       # Update user
DELETE /api/v1/users/:id       # Delete user
```

### Response Format
```json
{
  "success": true,
  "data": {...},
  "message": "Operation successful",
  "errors": null,
  "meta": {
    "pagination": {...},
    "timestamp": "2026-02-03T10:00:00Z"
  }
}
```

### Error Handling
```json
{
  "success": false,
  "data": null,
  "message": "Validation failed",
  "errors": [
    {
      "field": "email",
      "code": "INVALID_FORMAT",
      "message": "Invalid email format"
    }
  ]
}
```

## 🔄 Flujo de Desarrollo

### 1. Feature Development
```
1. Define domain entity
2. Create repository interface
3. Implement repository
4. Create service layer
5. Add HTTP handlers
6. Create frontend components
7. Wire up routing
8. Add tests
```

### 2. Database Changes
```
1. Create migration script
2. Update domain entities
3. Modify repositories
4. Update API contracts
5. Adjust frontend types
```

## 📊 Consideraciones de Escalabilidad

### Backend Scaling
- **Stateless design** para horizontal scaling
- **Database connection pooling**
- **Caching layer** (Redis) para sesiones
- **Load balancing** para múltiples instancias

### Frontend Scaling
- **Component-based architecture** para reutilización
- **Lazy loading** de rutas
- **CDN** para assets estáticos
- **Progressive Web App (PWA)** features

### Database Scaling
- **Indexing** en columnas frecuentemente consultadas
- **Read replicas** para consultas de solo lectura
- **Partitioning** por company_id para multi-tenancy
- **Connection pooling** para optimizar recursos

## 🚀 Deployment Architecture

### Production Environment
```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│             │    │             │    │             │
│ Load        │    │ Application │    │ Database    │
│ Balancer    │◄──►│ Servers     │◄──►│ Cluster     │
│ (Nginx)     │    │ (Go API)    │    │ (PostgreSQL)│
│             │    │             │    │             │
└─────────────┘    └─────────────┘    └─────────────┘
       │                   │                   │
       ▼                   ▼                   ▼
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│             │    │             │    │             │
│ Static      │    │ Cache       │    │ Monitoring  │
│ Assets      │    │ Layer       │    │ & Logs      │
│ (CDN)       │    │ (Redis)     │    │ (ELK Stack) │
│             │    │             │    │             │
└─────────────┘    └─────────────┘    └─────────────┘
```

---

**Estado**: 🚧 Documentación en desarrollo  
**Próxima actualización**: Implementación de autenticación JWT