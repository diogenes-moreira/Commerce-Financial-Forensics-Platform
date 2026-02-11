# costForensics

Cloud cost analysis platform for AWS, GCP, and Azure. Multi-tenant SaaS with database-per-tenant isolation, built with Go and React.

## What it does

- **Cost Tracking** — Ingest and categorize cloud spend by service, account, and time period
- **Budgets** — Set spending limits with automatic threshold alerts
- **Anomaly Detection** — Flag cost deviations with severity classification (low/medium/high/critical)
- **Reports** — Generate aggregated cost breakdowns by service and period
- **Multi-tenancy** — Full data isolation with a dedicated PostgreSQL database per tenant

## Architecture

```
┌──────────────┐       ┌──────────────────────────────────────────────┐
│   React SPA  │──────▶│              Go API (Gin)                    │
│  MUI + Vite  │  REST │                                              │
└──────────────┘       │  Middleware Chain                             │
                       │  Recovery → RequestID → Log → CORS → Auth    │
                       │  → TenantResolver → TenantDB → Handler       │
                       │                                              │
                       │  ┌─────────┐  ┌────────────┐  ┌──────────┐  │
                       │  │ Domain  │◀─│ Application │◀─│ Adapters │  │
                       │  │(entities│  │ (use cases) │  │(HTTP,DB) │  │
                       │  │  + VOs) │  └────────────┘  └──────────┘  │
                       │  └─────────┘                                 │
                       └──────────────┬───────────────────────────────┘
                                      │
                       ┌──────────────▼───────────────────────────────┐
                       │            PostgreSQL 16                      │
                       │                                              │
                       │  costforensics_system    (tenant registry)    │
                       │  cf_tenant_{id}          (per-tenant data)    │
                       └──────────────────────────────────────────────┘
```

### Backend Stack

| Component | Technology |
|---|---|
| Language | Go 1.23 |
| HTTP Framework | Gin |
| ORM | GORM |
| Database | PostgreSQL 16 |
| Auth | JWT (HMAC-SHA256) |
| Logging | Logrus (structured JSON) |
| Architecture | Hexagonal / Clean, rich domain models |

### Frontend Stack

| Component | Technology |
|---|---|
| Framework | React 19 |
| Language | TypeScript 5.7 |
| Build | Vite 6 |
| UI | Material UI 6 |
| Server State | TanStack Query 5 |
| Client State | Zustand 5 |
| Charts | Recharts 2 |
| HTTP Client | Axios |

## Quick Start

### Prerequisites

- Go 1.23+
- Node.js 20+
- Docker & Docker Compose

### 1. Clone and start PostgreSQL

```bash
git clone git@github.com:diogenes-moreira/Commerce-Financial-Forensics-Platform.git
cd Commerce-Financial-Forensics-Platform
docker-compose up -d postgres
```

### 2. Start the backend

```bash
cd backend
go mod download
go run ./cmd/server
```

The server starts on `http://localhost:8080`. It automatically runs database migrations on boot.

```bash
# Verify
curl http://localhost:8080/healthz
# {"status":"healthy"}
```

### 3. Start the frontend

```bash
cd frontend
npm install
npm run dev
```

Open `http://localhost:3000`. The Vite dev server proxies `/api` requests to the backend.

### Full Docker setup (alternative)

```bash
docker-compose up -d
# Backend: http://localhost:8080
```

## Project Structure

```
costForensics/
├── backend/
│   ├── cmd/server/main.go              # Entry point, DI wiring, graceful shutdown
│   ├── config/                         # YAML + env var configuration
│   ├── internal/
│   │   ├── domain/                     # Pure domain (NO external deps)
│   │   │   ├── shared/                 # Money, DateRange, Pagination value objects
│   │   │   ├── tenant/                 # Tenant entity + repository interface
│   │   │   ├── cloudaccount/           # Cloud account with sync state machine
│   │   │   ├── costrecord/             # Cost record with auto-categorization
│   │   │   ├── budget/                 # Budget with spend tracking + alerts
│   │   │   ├── anomaly/                # Anomaly with severity self-classification
│   │   │   └── costreport/             # Report with service breakdown
│   │   ├── application/                # Use-case services
│   │   ├── adapter/
│   │   │   ├── http/                   # Gin handlers, middleware, DTOs
│   │   │   └── postgres/              # GORM models, mappers, repositories
│   │   └── platform/                   # Database, logger, JWT auth
│   └── migrations/
├── frontend/
│   └── src/
│       ├── api/                        # Axios client + domain API modules
│       ├── features/                   # Dashboard, Costs, Budgets, Anomalies, Reports
│       ├── components/                 # AppLayout (Drawer+AppBar), DataTable, MoneyDisplay
│       ├── store/                      # Zustand (auth, tenant)
│       ├── hooks/                      # TanStack Query hooks
│       └── types/                      # TypeScript interfaces
├── docs/
│   ├── swagger.yaml                    # OpenAPI 3.0 specification
│   ├── MANUAL_IMPLEMENTACION.md        # Implementation manual (ES)
│   └── MANUAL_USO.md                   # User manual (ES)
├── docker-compose.yml
└── Makefile
```

## API Overview

All endpoints (except `/healthz`) require `Authorization: Bearer <jwt>`.

### Admin (Tenant Management)

| Method | Endpoint | Description |
|---|---|---|
| POST | `/api/v1/admin/tenants` | Create tenant (provisions DB) |
| GET | `/api/v1/admin/tenants` | List tenants |
| GET | `/api/v1/admin/tenants/:id` | Get tenant |
| PUT | `/api/v1/admin/tenants/:id` | Update tenant |
| DELETE | `/api/v1/admin/tenants/:id` | Delete tenant |

### Tenant-Scoped

| Method | Endpoint | Description |
|---|---|---|
| POST | `/api/v1/cloud-accounts` | Register cloud account |
| GET | `/api/v1/cloud-accounts` | List cloud accounts |
| GET | `/api/v1/cloud-accounts/:id` | Get cloud account |
| PUT | `/api/v1/cloud-accounts/:id` | Update cloud account |
| DELETE | `/api/v1/cloud-accounts/:id` | Delete cloud account |
| POST | `/api/v1/cloud-accounts/:id/sync` | Start cost sync |
| POST | `/api/v1/cost-records` | Create cost record |
| GET | `/api/v1/cost-records` | List costs (paginated, filterable) |
| GET | `/api/v1/cost-records/:id` | Get cost record |
| POST | `/api/v1/budgets` | Create budget |
| GET | `/api/v1/budgets` | List budgets |
| GET | `/api/v1/budgets/:id` | Get budget |
| PUT | `/api/v1/budgets/:id` | Update budget |
| POST | `/api/v1/budgets/:id/spend` | Record spend (may trigger alert) |
| GET | `/api/v1/anomalies` | List anomalies |
| GET | `/api/v1/anomalies/:id` | Get anomaly |
| POST | `/api/v1/anomalies/:id/resolve` | Resolve anomaly |
| POST | `/api/v1/cost-reports` | Generate cost report |
| GET | `/api/v1/cost-reports` | List reports |
| GET | `/api/v1/cost-reports/:id` | Get report |

Full OpenAPI spec: [`docs/swagger.yaml`](docs/swagger.yaml)

## Configuration

The backend loads `config/config.yaml` with environment variable overrides:

| Variable | Description | Default |
|---|---|---|
| `CF_SERVER_PORT` | HTTP port | `8080` |
| `CF_DB_HOST` | PostgreSQL host | `localhost` |
| `CF_DB_PORT` | PostgreSQL port | `5432` |
| `CF_DB_USER` | PostgreSQL user | `costforensics` |
| `CF_DB_PASSWORD` | PostgreSQL password | — |
| `CF_DB_SYSTEM_NAME` | System database name | `costforensics_system` |
| `CF_DB_SSLMODE` | SSL mode | `disable` |
| `CF_LOG_LEVEL` | Log level (debug/info/warn/error) | `info` |
| `CF_LOG_FORMAT` | Log format (text/json) | `text` |
| `CF_AUTH_JWT_SECRET` | JWT signing secret | — |

## Make Commands

```bash
make dev              # Run backend dev server
make build            # Compile backend binary
make test             # Run backend tests with race detector
make lint             # Run golangci-lint
make frontend-dev     # Run frontend dev server
make frontend-build   # Production frontend build
make docker-up        # Start Docker services
make docker-down      # Stop Docker services
```

## Documentation

- [OpenAPI 3.0 Spec](docs/swagger.yaml) — Full API specification with schemas and examples
- [Implementation Manual](docs/MANUAL_IMPLEMENTACION.md) — Architecture, setup, deployment, and how to extend
- [User Manual](docs/MANUAL_USO.md) — API usage with curl examples, frontend guide, error reference

## License

Private — All rights reserved.
