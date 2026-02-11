# costForensics

Commerce Financial Forensics platform. Multi-tenant SaaS that unifies cloud cost analysis, commerce operations, financial ledger, margin forensics, and data integration — with database-per-tenant isolation, built with Go and React.

## What it does

### Cloud Cost Analysis
- **Cost Tracking** — Ingest and categorize cloud spend (AWS/GCP/Azure) by service, account, and time period
- **Budgets** — Set spending limits with automatic threshold alerts
- **Anomaly Detection** — Flag cost deviations with severity classification (low/medium/high/critical)
- **Reports** — Generate aggregated cost breakdowns by service and period

### Commerce Operations
- **Products** — Catalog with SKU, EAN, UPC codes, unit cost/price, margin calculation, status lifecycle
- **Sellers** — Seller registry with commission rates and external ID mapping (e.g. `a19999`, `seller-a199999`)
- **Customers** — Customer database with segmentation (enterprise/mid-market/smb/consumer)
- **Orders** — Full order lifecycle (pending → confirmed → shipped → delivered → cancelled → refunded) with line items, price/cost snapshots, discount and tax breakdown

### Financial
- **Double-Entry Ledger** — Immutable ledger entries (debit/credit) tied to orders, refunds, and adjustments across account codes (revenue, COGS, discount, commission, shipping)
- **Forensic Events** — Append-only event log with SHA-256 hash integrity for tamper detection, tracking every price change, discount application, and order state transition
- **Promotion Rules & Discounts** — Percentage, fixed amount, buy-X-get-Y, and tiered discounts with funding split (seller/platform/shared), usage limits, and rule versioning
- **Bidirectional Payments** — Track inbound payments from customers and outbound payments to sellers with status lifecycle (pending → processed → failed → refunded)
- **Exchange Rates** — Currency conversion table for multi-currency operations

### Analytics
- **Margin Decomposition** — Per-order breakdown: revenue, COGS, discounts (seller vs platform funded), commission, shipping, fees, gross/net margin, real margin %
- **Drift Detection** — Compare snapshotted order prices against current catalog to identify deviations, with severity classification (low < 5%, medium 5-15%, high 15-30%, critical > 30%)
- **P&L Report** — Profit & Loss with drill-down by day, month, quarter, or year
- **Retention Cohorts** — Customer cohort analysis by first-order period, tracking repeat purchase behavior

### Data & Integrations
- **CSV Import** — Upload CSV files with progress tracking, per-row error logging, and support for products, orders, sellers, customers, and cost records
- **Platform Connectors** — Sync data from WooCommerce, Shopify, MedusaJS, and custom OMS via batch sync and real-time webhooks
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
│   │   │   ├── costreport/             # Report with service breakdown
│   │   │   ├── product/                # Product catalog (SKU, EAN, UPC, cost/price)
│   │   │   ├── seller/                 # Seller with commission calculation
│   │   │   ├── customer/               # Customer with segmentation
│   │   │   ├── order/                  # Order + line items, state machine, margin calc
│   │   │   ├── ledger/                 # Double-entry ledger entries
│   │   │   ├── forensicevent/          # Immutable events with hash integrity
│   │   │   ├── discount/               # Promotion rules + discount applications
│   │   │   ├── payment/                # Bidirectional payments (inbound/outbound)
│   │   │   ├── exchangerate/           # Currency exchange rates
│   │   │   ├── margin/                 # Margin breakdown calculator
│   │   │   ├── drift/                  # Price/cost drift detector
│   │   │   ├── importjob/              # CSV/bucket import job tracking
│   │   │   └── integration/            # Platform connector management
│   │   ├── application/                # Use-case services (one per domain)
│   │   ├── adapter/
│   │   │   ├── http/                   # Gin handlers, middleware, DTOs
│   │   │   └── postgres/              # GORM models, mappers, repositories
│   │   └── platform/
│   │       ├── database/               # DB connections, tenant DB manager, migrator
│   │       ├── auth/                   # JWT validation
│   │       ├── logger/                 # Logrus setup
│   │       ├── importer/               # CSV parser
│   │       └── connector/              # WooCommerce, Shopify, MedusaJS, custom OMS
│   └── migrations/
├── frontend/
│   └── src/
│       ├── api/                        # Axios client + domain API modules
│       ├── features/                   # Feature pages (15 domains)
│       │   ├── dashboard/              # Summary cards, charts
│       │   ├── cloud-accounts/         # Cloud provider accounts
│       │   ├── costs/                  # Cost explorer
│       │   ├── budgets/                # Budget management
│       │   ├── anomalies/              # Anomaly detection
│       │   ├── reports/                # Cost reports
│       │   ├── products/               # Product catalog
│       │   ├── sellers/                # Seller management
│       │   ├── customers/              # Customer database
│       │   ├── orders/                 # Order management
│       │   ├── ledger/                 # Ledger entries
│       │   ├── events/                 # Forensic event timeline
│       │   ├── discounts/              # Promotion rules
│       │   ├── payments/               # Payment tracking
│       │   ├── exchange-rates/         # Exchange rate table
│       │   ├── margins/                # Margin decomposition
│       │   ├── drift/                  # Drift detection
│       │   ├── pnl/                    # P&L report with drill-down
│       │   ├── cohorts/                # Retention cohort analysis
│       │   ├── imports/                # Data import jobs
│       │   ├── integrations/           # Platform connectors
│       │   └── settings/               # App settings
│       ├── components/                 # AppLayout, DataTable, MoneyDisplay
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

### Cloud Costs

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

### Commerce

| Method | Endpoint | Description |
|---|---|---|
| POST | `/api/v1/products` | Create product |
| GET | `/api/v1/products` | List products |
| GET | `/api/v1/products/:id` | Get product |
| PUT | `/api/v1/products/:id` | Update product |
| POST | `/api/v1/sellers` | Create seller |
| GET | `/api/v1/sellers` | List sellers |
| GET | `/api/v1/sellers/:id` | Get seller |
| PUT | `/api/v1/sellers/:id` | Update seller |
| POST | `/api/v1/customers` | Create customer |
| GET | `/api/v1/customers` | List customers |
| GET | `/api/v1/customers/:id` | Get customer |
| PUT | `/api/v1/customers/:id` | Update customer |
| POST | `/api/v1/orders` | Create order |
| GET | `/api/v1/orders` | List orders |
| GET | `/api/v1/orders/:id` | Get order |
| POST | `/api/v1/orders/:id/items` | Add line item |
| POST | `/api/v1/orders/:id/confirm` | Confirm order |
| POST | `/api/v1/orders/:id/ship` | Mark shipped |
| POST | `/api/v1/orders/:id/deliver` | Mark delivered |
| POST | `/api/v1/orders/:id/cancel` | Cancel order |
| POST | `/api/v1/orders/:id/refund` | Refund order |

### Financial

| Method | Endpoint | Description |
|---|---|---|
| GET | `/api/v1/ledger/entries` | List ledger entries (filter by order, account, date) |
| GET | `/api/v1/ledger/balance` | Get account balance for period |
| GET | `/api/v1/events` | List forensic events (filter by entity, type, date) |
| POST | `/api/v1/promotion-rules` | Create promotion rule |
| GET | `/api/v1/promotion-rules` | List promotion rules |
| GET | `/api/v1/promotion-rules/:id` | Get promotion rule |
| POST | `/api/v1/promotion-rules/:id/disable` | Disable promotion rule |
| GET | `/api/v1/discount-applications` | List discount applications |
| POST | `/api/v1/payments` | Create payment |
| GET | `/api/v1/payments` | List payments (filter by order, direction, status) |
| GET | `/api/v1/payments/:id` | Get payment |
| POST | `/api/v1/payments/:id/process` | Mark payment processed |
| POST | `/api/v1/payments/:id/fail` | Mark payment failed |
| POST | `/api/v1/payments/:id/refund` | Mark payment refunded |
| POST | `/api/v1/exchange-rates` | Create exchange rate |
| GET | `/api/v1/exchange-rates` | List exchange rates |
| GET | `/api/v1/exchange-rates/latest` | Get latest rate for currency pair |
| GET | `/api/v1/exchange-rates/:id` | Get exchange rate |

### Analytics

| Method | Endpoint | Description |
|---|---|---|
| GET | `/api/v1/orders/:id/margin` | Margin breakdown for order |
| GET | `/api/v1/margins` | Batch margin report (filter by date, seller) |
| GET | `/api/v1/orders/:id/drift` | Drift analysis for order |
| GET | `/api/v1/drift-report` | Batch drift report |
| GET | `/api/v1/pnl` | P&L report (drill-down: day/month/quarter/year) |
| GET | `/api/v1/cohorts` | Retention cohort analysis |

### Data & Integrations

| Method | Endpoint | Description |
|---|---|---|
| POST | `/api/v1/imports/upload` | Upload CSV file (multipart) |
| GET | `/api/v1/imports` | List import jobs |
| GET | `/api/v1/imports/:id` | Get import job status |
| POST | `/api/v1/integrations` | Create integration |
| GET | `/api/v1/integrations` | List integrations |
| GET | `/api/v1/integrations/:id` | Get integration |
| PUT | `/api/v1/integrations/:id` | Update integration |
| DELETE | `/api/v1/integrations/:id` | Delete integration |
| POST | `/api/v1/integrations/:id/sync` | Trigger data sync |
| POST | `/api/v1/webhooks/:integration_id` | Receive webhook |

Full OpenAPI spec: [`docs/swagger.yaml`](docs/swagger.yaml)

## Database Tables (per tenant)

```
Cloud Costs:    cloud_accounts, cost_records, budgets, anomalies, cost_reports
Commerce:       products, sellers, customers, orders, order_items
Financial:      ledger_entries, forensic_events, promotion_rules, discount_applications,
                payments, exchange_rates
Data:           import_jobs, integrations
```

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
