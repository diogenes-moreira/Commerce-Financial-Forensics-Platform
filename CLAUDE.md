# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

costForensics — Cloud cost analysis platform (AWS/GCP/Azure). Multi-tenant SaaS with database-per-tenant isolation.

## Build & Run Commands

### Backend
```bash
cd backend && go build ./...        # Compile
cd backend && go test ./... -v -race # Run tests
cd backend && go run ./cmd/server   # Dev server (port 8080)
```

### Frontend
```bash
cd frontend && npm install          # Install deps
cd frontend && npm run dev          # Dev server (port 3000)
cd frontend && npm run build        # Production build
```

### Docker
```bash
docker-compose up -d                # Start PostgreSQL + backend
```

## Architecture

- **Backend**: Go, Gin, GORM, Logrus, PostgreSQL
- **Frontend**: React, TypeScript, Vite, Material UI, Zustand, TanStack Query, Recharts
- **Multi-tenancy**: Database-per-tenant (system DB for tenant registry, dedicated PG DB per tenant)
- **Pattern**: Hexagonal/clean architecture with rich domain models

### Backend Structure
- `domain/` — Pure domain entities, value objects, repository interfaces (NO external deps)
- `application/` — Use-case services (depends on domain only)
- `adapter/http/` — Gin handlers, middleware, DTOs
- `adapter/postgres/` — GORM models, mappers, repository implementations
- `platform/` — Cross-cutting: database connections, logger, JWT auth
- `config/` — YAML + env var configuration

### Key Rules
- Domain entities have unexported fields; mutation only via methods that enforce invariants
- GORM models are separate from domain entities; mappers convert between layers
- Tenant-scoped handlers get `*gorm.DB` from gin.Context via `c.MustGet("tenant_db")`
- Middleware chain: Recovery → RequestID → Logging → CORS → Auth → TenantResolver → TenantDB → Handler

### Frontend Structure
- `api/` — Axios client + per-domain API modules
- `features/` — Feature-based pages (dashboard, costs, budgets, anomalies, reports)
- `components/` — Shared layout and common components (DataTable, MoneyDisplay)
- `store/` — Zustand stores (auth, tenant)
- `hooks/` — TanStack Query hooks per domain
