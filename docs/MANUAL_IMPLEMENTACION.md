# Manual de Implementacion — costForensics

## 1. Requisitos Previos

| Componente | Version minima | Proposito |
|---|---|---|
| Go | 1.23+ | Backend |
| Node.js | 20+ | Frontend |
| PostgreSQL | 16+ | Base de datos |
| Docker & Docker Compose | 24+ / 2.20+ | Contenedores |
| Git | 2.40+ | Control de versiones |

## 2. Clonar el Repositorio

```bash
git clone git@github.com:diogenes-moreira/Commerce-Financial-Forensics-Platform.git
cd Commerce-Financial-Forensics-Platform
```

## 3. Estructura del Proyecto

```
costForensics/
├── Makefile                 # Comandos raiz del monorepo
├── docker-compose.yml       # PostgreSQL + Backend
├── CLAUDE.md                # Guia para Claude Code
├── docs/
│   ├── swagger.yaml         # Especificacion OpenAPI 3.0
│   ├── MANUAL_IMPLEMENTACION.md
│   └── MANUAL_USO.md
├── backend/                 # Go API (Gin + GORM)
│   ├── cmd/server/main.go   # Entry point
│   ├── config/              # Configuracion YAML + env vars
│   ├── internal/
│   │   ├── domain/          # Entidades de dominio puras
│   │   ├── application/     # Servicios de aplicacion (use cases)
│   │   ├── adapter/         # HTTP handlers, repos PostgreSQL
│   │   └── platform/        # DB, logger, auth
│   └── migrations/          # SQL migrations
└── frontend/                # React SPA (Vite + MUI)
    └── src/
        ├── api/             # Cliente HTTP
        ├── features/        # Paginas por dominio
        ├── components/      # Componentes compartidos
        ├── store/           # Estado Zustand
        └── hooks/           # TanStack Query hooks
```

## 4. Configuracion

### 4.1 Variables de Entorno (Backend)

El backend carga configuracion desde `backend/config/config.yaml` con overrides via variables de entorno:

| Variable | Descripcion | Default |
|---|---|---|
| `CF_SERVER_PORT` | Puerto del servidor HTTP | `8080` |
| `CF_DB_HOST` | Host de PostgreSQL | `localhost` |
| `CF_DB_PORT` | Puerto de PostgreSQL | `5432` |
| `CF_DB_USER` | Usuario de PostgreSQL | `costforensics` |
| `CF_DB_PASSWORD` | Password de PostgreSQL | *(requerido)* |
| `CF_DB_SYSTEM_NAME` | Nombre de la DB del sistema | `costforensics_system` |
| `CF_DB_SSLMODE` | Modo SSL de PostgreSQL | `disable` |
| `CF_LOG_LEVEL` | Nivel de log (`debug`, `info`, `warn`, `error`) | `info` |
| `CF_LOG_FORMAT` | Formato de log (`text`, `json`) | `text` |
| `CF_AUTH_JWT_SECRET` | Secreto para firmar/validar JWT | *(requerido en produccion)* |

### 4.2 Archivo config.yaml

```yaml
server:
  port: 8080

database:
  host: localhost
  port: 5432
  user: costforensics
  password: localdev
  system_name: costforensics_system
  sslmode: disable
  max_open_conns: 25
  max_idle_conns: 5

log:
  level: debug
  format: text

auth:
  jwt_secret: dev-secret-change-in-production
```

> **IMPORTANTE**: En produccion, usar `CF_AUTH_JWT_SECRET` con un secreto fuerte y `CF_DB_SSLMODE=require`.

## 5. Instalacion y Ejecucion

### 5.1 Desarrollo Local con Docker (recomendado)

```bash
# 1. Levantar PostgreSQL
docker-compose up -d postgres

# 2. Instalar dependencias del backend
cd backend && go mod download && cd ..

# 3. Ejecutar backend (aplica migraciones automaticamente)
cd backend && go run ./cmd/server

# 4. En otra terminal, instalar y ejecutar frontend
cd frontend && npm install && npm run dev
```

- Backend: http://localhost:8080
- Frontend: http://localhost:3000
- El frontend tiene proxy configurado para `/api` -> `localhost:8080`

### 5.2 Docker Compose Completo

```bash
# Levanta PostgreSQL + Backend
docker-compose up -d

# Verificar salud
curl http://localhost:8080/healthz
# {"status":"healthy"}
```

### 5.3 Solo Backend (sin Docker)

Requisito: PostgreSQL corriendo localmente con la base de datos `costforensics_system` creada.

```bash
# Crear la base de datos del sistema
psql -U postgres -c "CREATE USER costforensics WITH PASSWORD 'localdev';"
psql -U postgres -c "CREATE DATABASE costforensics_system OWNER costforensics;"

# Ejecutar
cd backend && go run ./cmd/server
```

## 6. Arquitectura del Backend

### 6.1 Arquitectura Hexagonal

```
        ┌──────────────────────────────────┐
        │            Handlers              │ ← adapter/http/
        │      (Gin HTTP handlers)         │
        └──────────┬───────────────────────┘
                   │
        ┌──────────▼───────────────────────┐
        │       Application Services       │ ← application/
        │         (Use Cases)              │
        └──────────┬───────────────────────┘
                   │
        ┌──────────▼───────────────────────┐
        │        Domain Entities           │ ← domain/
        │  (Entidades, Value Objects,      │
        │   Repository Interfaces)         │
        └──────────┬───────────────────────┘
                   │ (implementado por)
        ┌──────────▼───────────────────────┐
        │     PostgreSQL Repositories      │ ← adapter/postgres/
        │    (GORM models + mappers)       │
        └─────────────────────────────────┘
```

**Regla fundamental**: El dominio NO importa paquetes externos. Los adapters implementan las interfaces del dominio.

### 6.2 Modelos de Dominio Ricos (anti-anemicos)

Los campos de las entidades son **no exportados**. Todo cambio de estado se realiza via metodos que validan invariantes:

```go
// NO hacer esto:
budget.Spent = newAmount  // ❌ campos no exportados

// Hacer esto:
alert, err := budget.RecordSpend(amount)  // ✅ devuelve alerta si supera umbral
```

Ejemplos de comportamiento en entidades:
- `CloudAccount.BeginSync()` — maquina de estados: solo permite sincronizar si esta `idle`
- `Budget.RecordSpend(amount)` — agrega gasto y retorna `*Alert` si se supera el umbral
- `CostAnomaly` — auto-clasifica severidad segun porcentaje de desviacion
- `Money.Add(other)` — retorna nuevo Money, valida misma moneda

### 6.3 Multi-tenancy: Database per Tenant

```
┌─────────────────────────────────┐
│     costforensics_system        │  ← Base de datos del sistema
│  ┌───────────────────────────┐  │
│  │ tenants                   │  │  ← Registro de tenants
│  │  id | name | slug | db   │  │
│  └───────────────────────────┘  │
└─────────────────────────────────┘

┌─────────────────────────────────┐
│    cf_tenant_550e8400           │  ← DB del Tenant A
│  cloud_accounts                 │
│  cost_records                   │
│  budgets / budget_alerts        │
│  anomalies                      │
│  cost_reports                   │
└─────────────────────────────────┘

┌─────────────────────────────────┐
│    cf_tenant_a1b2c3d4           │  ← DB del Tenant B
│  (mismas tablas)                │
└─────────────────────────────────┘
```

**Flujo de una request tenant-scoped**:
1. `Auth` middleware extrae `tenant_id` del JWT
2. `TenantResolver` busca el tenant en la DB del sistema y verifica que este activo
3. `TenantDB` obtiene (o crea lazy) la conexion a la DB del tenant via `TenantDBManager`
4. El handler accede a la DB del tenant con `c.MustGet("tenant_db").(*gorm.DB)`

### 6.4 Cadena de Middleware

```
Request → Recovery → RequestID → Logging → CORS → Auth → TenantResolver → TenantDB → Handler
```

| Middleware | Funcion |
|---|---|
| Recovery | Captura panics, devuelve 500 |
| RequestID | Genera/propaga `X-Request-ID` |
| Logging | Log estructurado con latencia, status, path |
| CORS | Headers CORS permisivos (dev) |
| Auth | Valida JWT, inyecta claims en contexto |
| TenantResolver | Resuelve tenant desde JWT, verifica estado activo |
| TenantDB | Inyecta `*gorm.DB` del tenant en contexto |

### 6.5 GORM Models vs Domain Entities

```
domain/tenant/entity.go          ←→  adapter/postgres/model/tenant_model.go
(campos no exportados,                (tags GORM, campos exportados)
 metodos de negocio)

                    adapter/postgres/mapper/tenant_mapper.go
                    (TenantToModel / TenantToDomain)
```

## 7. Arquitectura del Frontend

### 7.1 Stack

| Libreria | Proposito |
|---|---|
| React 19 | UI framework |
| TypeScript 5.7 | Tipado estatico |
| Vite 6 | Build tool + dev server |
| Material UI 6 | Componentes visuales |
| Zustand 5 | Estado del cliente (auth, tenant) |
| TanStack Query 5 | Cache de datos del servidor |
| Recharts 2 | Graficos de costos |
| Axios | Cliente HTTP |
| React Router 7 | Ruteo SPA |

### 7.2 Estructura de Carpetas

```
src/
├── api/              # Un modulo por dominio (costsApi, budgetsApi, etc.)
├── components/
│   ├── layout/       # AppLayout (Drawer + AppBar + Outlet)
│   └── common/       # DataTable generico, MoneyDisplay
├── features/         # Una carpeta por pagina/feature
│   ├── dashboard/    # DashboardPage con cards resumen + grafico Recharts
│   ├── cloud-accounts/
│   ├── costs/
│   ├── budgets/
│   ├── anomalies/
│   ├── reports/
│   └── settings/
├── hooks/            # Custom hooks con TanStack Query
├── store/            # Zustand: authStore, tenantStore
├── types/            # Interfaces TS que espejan los DTOs del backend
└── utils/            # formatCurrency, dateUtils
```

### 7.3 Patron de Datos

```
API Module (Axios) → TanStack Query Hook → Feature Page Component
     ↓
  Zustand Store (solo para estado del cliente: auth token, tenant seleccionado)
```

## 8. Base de Datos

### 8.1 Esquema del Sistema (`costforensics_system`)

```sql
CREATE TABLE tenants (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    db_name VARCHAR(255) NOT NULL UNIQUE,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### 8.2 Esquema del Tenant (aplicado a cada `cf_tenant_{id}`)

Tablas: `cloud_accounts`, `cost_records`, `budgets`, `budget_alerts`, `anomalies`, `cost_reports`.

Las migraciones se ejecutan automaticamente cuando se provisiona un nuevo tenant.

## 9. Autenticacion (JWT)

El sistema usa JWT con HMAC-SHA256. Claims requeridos:

```json
{
  "user_id": "uuid",
  "tenant_id": "uuid",
  "email": "user@example.com",
  "role": "admin",
  "exp": 1738368000,
  "iat": 1738281600
}
```

Para desarrollo, se puede generar un token usando la funcion `GenerateToken` del paquete `platform/auth`:

```go
validator := auth.NewJWTValidator("dev-secret-change-in-production")
token, _ := validator.GenerateToken(userID, tenantID, "admin@test.com", "admin")
```

## 10. Comandos Utiles

```bash
# Backend
make build          # Compilar binario
make dev            # Ejecutar en modo desarrollo
make test           # Ejecutar tests con race detector
make lint           # Ejecutar golangci-lint

# Frontend
make frontend-dev   # Vite dev server
make frontend-build # Build de produccion

# Docker
make docker-up      # Levantar servicios
make docker-down    # Detener servicios
```

## 11. Despliegue en Produccion

### 11.1 Checklist

- [ ] Configurar `CF_AUTH_JWT_SECRET` con secreto fuerte (32+ caracteres)
- [ ] Configurar `CF_DB_SSLMODE=require`
- [ ] Configurar `CF_DB_PASSWORD` con credencial segura
- [ ] Configurar `CF_LOG_FORMAT=json` para log aggregation
- [ ] Configurar `CF_LOG_LEVEL=info`
- [ ] Construir imagen Docker: `docker build -t costforensics-backend ./backend`
- [ ] Construir frontend: `cd frontend && npm run build` y servir `dist/` con Nginx o CDN
- [ ] Configurar CORS con dominios especificos (modificar `middleware/cors.go`)
- [ ] Configurar health checks apuntando a `/healthz`
- [ ] Configurar backups de PostgreSQL (sistema + todas las DBs de tenants)

### 11.2 Docker Build

```bash
# Backend
docker build -t costforensics-backend:latest ./backend

# Ejecutar
docker run -p 8080:8080 \
  -e CF_DB_HOST=your-db-host \
  -e CF_DB_PASSWORD=your-password \
  -e CF_AUTH_JWT_SECRET=your-secret \
  -e CF_LOG_FORMAT=json \
  costforensics-backend:latest
```

### 11.3 Escalabilidad

- El backend es stateless; se puede escalar horizontalmente detras de un load balancer
- `TenantDBManager` mantiene un pool de conexiones por tenant en memoria; en un cluster, cada instancia mantiene su propio pool
- PostgreSQL debe dimensionarse considerando: `max_connections >= (instancias_backend * max_open_conns * num_tenants)`

## 12. Agregar un Nuevo Dominio

Para agregar una nueva entidad al sistema:

1. **Domain** (`internal/domain/nuevo/`):
   - `entity.go` — Entidad con campos no exportados + constructor + metodos
   - `repository.go` — Interface del repositorio
   - `errors.go` — Errores de dominio

2. **Application** (`internal/application/nuevo/`):
   - `service.go` — Use cases que orquestan dominio + repositorio

3. **Adapter - Postgres** (`internal/adapter/postgres/`):
   - `model/nuevo_model.go` — Struct GORM con tags de DB
   - `mapper/nuevo_mapper.go` — Funciones ToModel/ToDomain
   - `tenant/nuevo_repo.go` — Implementacion del repositorio

4. **Adapter - HTTP** (`internal/adapter/http/`):
   - `dto/request/` — Structs de request con binding tags
   - `dto/response/` — Structs de response + funcion `FromDomain`
   - `handler/nuevo.go` — Handler HTTP

5. **Router** — Registrar rutas en `router.go`

6. **Migraciones** — Agregar SQL en `migrations/tenant/` o `migrations/system/`

7. **Frontend** — `api/nuevo.ts`, `hooks/useNuevo.ts`, `features/nuevo/NuevoPage.tsx`
