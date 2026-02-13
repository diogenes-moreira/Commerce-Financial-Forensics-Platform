# Manual de Uso — costForensics

## 1. Introduccion

costForensics es una plataforma SaaS de analisis de costos cloud que permite a las organizaciones monitorear, presupuestar y detectar anomalias en sus gastos de AWS, GCP y Azure.

### Funcionalidades Principales

- **Dashboard** — Vista consolidada de gastos, presupuestos y anomalias
- **Cuentas Cloud** — Registro y sincronizacion de cuentas AWS/GCP/Azure
- **Explorador de Costos** — Consulta detallada con filtros por servicio, categoria y fecha
- **Presupuestos** — Definicion de limites de gasto con alertas automaticas
- **Anomalias** — Deteccion automatica de desviaciones de costos
- **Reportes** — Generacion de reportes agregados por periodo

## 2. Acceso al Sistema

### 2.1 Frontend (Interfaz Web)

Acceder a la aplicacion web en:
- Desarrollo: `http://localhost:3000`
- Produccion: URL proporcionada por el administrador

### 2.2 API REST

Base URL:
- Desarrollo: `http://localhost:8080/api/v1`
- Produccion: `https://api.costforensics.io/api/v1`

Todas las peticiones requieren autenticacion JWT:

```bash
# Header requerido en todas las peticiones
Authorization: Bearer <jwt_token>
```

### 2.3 Verificar Salud del Servicio

```bash
curl http://localhost:8080/healthz
```

Respuesta esperada:
```json
{"status": "healthy"}
```

## 3. Gestion de Tenants (Administrador)

Los tenants representan organizaciones independientes. Cada tenant tiene su propia base de datos aislada.

### 3.1 Crear un Tenant

```bash
curl -X POST http://localhost:8080/api/v1/admin/tenants \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Acme Corporation",
    "slug": "acme-corp"
  }'
```

Respuesta:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Acme Corporation",
  "slug": "acme-corp",
  "db_name": "cf_tenant_550e8400",
  "status": "active",
  "created_at": "2026-02-11T10:00:00Z",
  "updated_at": "2026-02-11T10:00:00Z"
}
```

> Al crear un tenant, el sistema automaticamente:
> 1. Crea una base de datos PostgreSQL dedicada
> 2. Ejecuta las migraciones del esquema
> 3. Registra el tenant en el sistema

### 3.2 Listar Tenants

```bash
curl http://localhost:8080/api/v1/admin/tenants \
  -H "Authorization: Bearer $TOKEN"
```

### 3.3 Actualizar un Tenant

```bash
curl -X PUT http://localhost:8080/api/v1/admin/tenants/{id} \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "Acme Corp Updated"}'
```

### 3.4 Eliminar un Tenant

```bash
curl -X DELETE http://localhost:8080/api/v1/admin/tenants/{id} \
  -H "Authorization: Bearer $TOKEN"
```

> **Nota**: Esto elimina el registro del tenant del sistema pero NO elimina la base de datos del tenant.

## 4. Cuentas Cloud

### 4.1 Registrar una Cuenta Cloud

```bash
curl -X POST http://localhost:8080/api/v1/cloud-accounts \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "aws",
    "name": "Production AWS",
    "external_id": "123456789012"
  }'
```

Proveedores soportados: `aws`, `gcp`, `azure`

### 4.2 Listar Cuentas

```bash
curl http://localhost:8080/api/v1/cloud-accounts \
  -H "Authorization: Bearer $TOKEN"
```

Respuesta:
```json
[
  {
    "id": "...",
    "provider": "aws",
    "name": "Production AWS",
    "external_id": "123456789012",
    "status": "active",
    "sync_status": "idle",
    "last_synced_at": null,
    "created_at": "2026-02-11T10:00:00Z",
    "updated_at": "2026-02-11T10:00:00Z"
  }
]
```

### 4.3 Iniciar Sincronizacion

```bash
curl -X POST http://localhost:8080/api/v1/cloud-accounts/{id}/sync \
  -H "Authorization: Bearer $TOKEN"
```

**Estados de sincronizacion**:
- `idle` — Sin sincronizacion en curso
- `running` — Sincronizacion en progreso
- `failed` — Ultima sincronizacion fallo

> Solo se puede iniciar una sincronizacion si el estado actual es `idle` o `failed`.

### 4.4 Actualizar Cuenta

```bash
curl -X PUT http://localhost:8080/api/v1/cloud-accounts/{id} \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "AWS Produccion Renombrada"}'
```

### 4.5 Eliminar Cuenta

```bash
curl -X DELETE http://localhost:8080/api/v1/cloud-accounts/{id} \
  -H "Authorization: Bearer $TOKEN"
```

## 5. Registros de Costos

### 5.1 Crear un Registro de Costo

```bash
curl -X POST http://localhost:8080/api/v1/cost-records \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "cloud_account_id": "550e8400-e29b-41d4-a716-446655440000",
    "service": "Amazon EC2",
    "amount_cents": 15023,
    "currency": "USD",
    "usage_date": "2026-02-10"
  }'
```

**Sobre los montos**: Todos los valores monetarios se expresan en **centavos**. Por ejemplo, `15023` = $150.23 USD.

**Categorizacion automatica**: El sistema asigna automaticamente una categoria basada en el nombre del servicio:

| Categoria | Servicios ejemplo |
|---|---|
| `compute` | EC2, Lambda, Fargate, Cloud Functions |
| `storage` | S3, EBS, GCS |
| `network` | VPC, CloudFront, ELB |
| `database` | RDS, DynamoDB, Aurora, Spanner |
| `analytics` | Athena, BigQuery, Redshift |
| `ml` | SageMaker, AI Platform |
| `other` | Cualquier otro servicio |

### 5.2 Listar Registros (con filtros)

```bash
# Sin filtros (paginado, 20 por pagina)
curl "http://localhost:8080/api/v1/cost-records" \
  -H "Authorization: Bearer $TOKEN"

# Con filtros
curl "http://localhost:8080/api/v1/cost-records?service=Amazon+EC2&start_date=2026-02-01&end_date=2026-02-28&page=1&page_size=50" \
  -H "Authorization: Bearer $TOKEN"

# Filtrar por cuenta y categoria
curl "http://localhost:8080/api/v1/cost-records?cloud_account_id={uuid}&category=compute" \
  -H "Authorization: Bearer $TOKEN"
```

**Parametros de filtro**:

| Parametro | Tipo | Descripcion |
|---|---|---|
| `cloud_account_id` | UUID | Filtrar por cuenta cloud |
| `service` | string | Filtrar por nombre de servicio |
| `category` | string | Filtrar por categoria |
| `start_date` | date | Fecha inicio (YYYY-MM-DD) |
| `end_date` | date | Fecha fin (YYYY-MM-DD) |
| `page` | int | Numero de pagina (default: 1) |
| `page_size` | int | Registros por pagina (default: 20, max: 100) |

Respuesta paginada:
```json
{
  "items": [
    {
      "id": "...",
      "cloud_account_id": "...",
      "service": "Amazon EC2",
      "category": "compute",
      "amount_cents": 15023,
      "currency": "USD",
      "usage_date": "2026-02-10",
      "tags": {},
      "created_at": "2026-02-11T10:00:00Z"
    }
  ],
  "total_count": 142,
  "page": 1,
  "page_size": 20
}
```

### 5.3 Obtener Registro por ID

```bash
curl http://localhost:8080/api/v1/cost-records/{id} \
  -H "Authorization: Bearer $TOKEN"
```

## 6. Presupuestos

### 6.1 Crear un Presupuesto

```bash
curl -X POST http://localhost:8080/api/v1/budgets \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Q1 2026 Infraestructura",
    "amount_cents": 5000000,
    "currency": "USD",
    "period_start": "2026-01-01",
    "period_end": "2026-03-31",
    "alert_threshold_pct": 80
  }'
```

- `amount_cents`: Limite del presupuesto en centavos ($50,000 = 5000000)
- `alert_threshold_pct`: Porcentaje de uso a partir del cual se generan alertas (default: 80)

### 6.2 Registrar Gasto

```bash
curl -X POST http://localhost:8080/api/v1/budgets/{id}/spend \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "amount_cents": 250000,
    "currency": "USD"
  }'
```

Respuesta (con alerta si se supera el umbral):
```json
{
  "budget": {
    "id": "...",
    "name": "Q1 2026 Infraestructura",
    "amount_cents": 5000000,
    "spent_cents": 4250000,
    "usage_pct": 85,
    "alert_threshold_pct": 80,
    "status": "active",
    ...
  },
  "alert": {
    "id": "...",
    "budget_id": "...",
    "threshold_pct": 80,
    "actual_pct": 85,
    "message": "Budget usage at 85% (threshold: 80%)",
    "triggered_at": "2026-02-11T15:30:00Z"
  }
}
```

> Si el gasto no supera el umbral, `alert` sera `null`.

### 6.3 Listar Presupuestos

```bash
curl http://localhost:8080/api/v1/budgets \
  -H "Authorization: Bearer $TOKEN"
```

### 6.4 Actualizar Presupuesto

```bash
curl -X PUT http://localhost:8080/api/v1/budgets/{id} \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "Q1 2026 Infra (Ajustado)"}'
```

## 7. Anomalias

Las anomalias son desviaciones significativas entre el costo esperado y el costo real de un servicio.

### 7.1 Clasificacion de Severidad

| Severidad | Desviacion |
|---|---|
| `low` | < 50% |
| `medium` | 50% - 99% |
| `high` | 100% - 199% |
| `critical` | >= 200% |

### 7.2 Listar Anomalias

```bash
curl http://localhost:8080/api/v1/anomalies \
  -H "Authorization: Bearer $TOKEN"
```

Respuesta:
```json
[
  {
    "id": "...",
    "cloud_account_id": "...",
    "service": "Amazon RDS",
    "expected_cents": 100000,
    "actual_cents": 350000,
    "currency": "USD",
    "deviation_pct": 250.0,
    "severity": "critical",
    "status": "open",
    "detected_at": "2026-02-11T08:00:00Z",
    "resolved_at": null
  }
]
```

### 7.3 Resolver una Anomalia

```bash
curl -X POST http://localhost:8080/api/v1/anomalies/{id}/resolve \
  -H "Authorization: Bearer $TOKEN"
```

> Una anomalia solo puede resolverse si su estado es `open`. Intentar resolver una anomalia ya resuelta retorna error 422.

## 8. Reportes de Costos

### 8.1 Generar un Reporte

```bash
curl -X POST http://localhost:8080/api/v1/cost-reports \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Reporte Febrero 2026",
    "report_type": "monthly",
    "period_start": "2026-02-01",
    "period_end": "2026-02-28",
    "currency": "USD"
  }'
```

Tipos de reporte: `daily`, `weekly`, `monthly`, `custom`

Respuesta:
```json
{
  "id": "...",
  "name": "Reporte Febrero 2026",
  "report_type": "monthly",
  "period_start": "2026-02-01",
  "period_end": "2026-02-28",
  "total_cents": 7500000,
  "currency": "USD",
  "breakdown": {
    "Amazon EC2": 3500000,
    "Amazon S3": 1200000,
    "Amazon RDS": 2800000
  },
  "generated_at": "2026-02-11T16:00:00Z"
}
```

El reporte agrega automaticamente todos los registros de costos dentro del periodo y genera un desglose por servicio.

### 8.2 Listar Reportes

```bash
curl http://localhost:8080/api/v1/cost-reports \
  -H "Authorization: Bearer $TOKEN"
```

### 8.3 Obtener Reporte por ID

```bash
curl http://localhost:8080/api/v1/cost-reports/{id} \
  -H "Authorization: Bearer $TOKEN"
```

## 9. Interfaz Web (Frontend)

### 9.1 Dashboard

La pagina principal muestra:
- **Cards resumen**: Gasto total MTD, presupuestos activos, anomalias abiertas, cuentas cloud
- **Grafico de barras**: Costos por categoria (Compute, Storage, Network, Database, Otros)

### 9.2 Cloud Accounts

Tabla con todas las cuentas cloud registradas mostrando:
- Nombre y proveedor
- Estado (active/inactive)
- Estado de sincronizacion (idle/running/failed)
- Ultima sincronizacion

### 9.3 Cost Explorer

Tabla paginada de todos los registros de costos con:
- Servicio y categoria
- Monto formateado
- Fecha de uso

### 9.4 Budgets

Tabla de presupuestos con:
- Nombre y limite
- Gasto acumulado
- Barra de progreso visual (verde < 80%, amarillo >= 80%, rojo >= 100%)
- Periodo

### 9.5 Anomalies

Tabla de anomalias con:
- Servicio afectado
- Costo esperado vs real
- Porcentaje de desviacion
- Chip de severidad con color (critical=rojo, high=rojo, medium=amarillo, low=azul)
- Estado (open/resolved)

### 9.6 Reports

Tabla de reportes generados con:
- Nombre y tipo
- Periodo cubierto
- Total consolidado
- Fecha de generacion

### 9.7 Settings

Pagina de configuracion para integraciones cloud, preferencias de notificacion y gestion de equipo.

## 10. Codigos de Error Comunes

| Codigo HTTP | Significado | Ejemplo |
|---|---|---|
| 400 | Request invalido | JSON malformado, campo requerido faltante |
| 401 | No autenticado | Token JWT invalido o expirado |
| 403 | Prohibido | Tenant suspendido |
| 404 | No encontrado | Recurso con el ID dado no existe |
| 422 | Error de negocio | Slug de tenant duplicado, cuenta ya sincronizando |
| 500 | Error interno | Error de base de datos inesperado |

Formato de error:
```json
{
  "error": "descripcion del error"
}
```

## 11. Importacion CSV

### 11.1 Subir un CSV

```bash
curl -X POST http://localhost:8080/api/v1/imports/upload \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@datos.csv" \
  -F "entity_type=products" \
  -F "name=Importacion Febrero 2026"
```

**Parametros del formulario**:

| Parametro | Tipo | Requerido | Descripcion |
|---|---|---|---|
| `file` | file | Si | Archivo CSV (la primera fila debe ser el encabezado) |
| `entity_type` | string | Si | Tipo de entidad: `products`, `orders`, `sellers`, `customers`, `cost_records` |
| `name` | string | No | Nombre descriptivo del job (default: `"CSV Import"`) |

### 11.2 Formato General

- **Delimitador**: coma (`,`)
- **Primera fila**: encabezados (nombres de columna)
- **Encoding**: UTF-8
- Los espacios al inicio de cada campo se recortan automaticamente

### 11.3 Formato por Tipo de Entidad

#### `products`

| Columna | Tipo | Requerido | Descripcion |
|---|---|---|---|
| `sku` | string | Si | Codigo unico del producto |
| `name` | string | Si | Nombre del producto |
| `description` | string | No | Descripcion del producto |
| `category` | string | No | Categoria (ej: `electronics`, `accessories`) |
| `ean` | string | No | Codigo EAN-13 |
| `upc` | string | No | Codigo UPC |
| `unit_cost_cents` | int | Si | Costo unitario en centavos |
| `unit_price_cents` | int | Si | Precio de venta en centavos |
| `currency` | string | No | Moneda (default: `USD`) |

Ejemplo:
```csv
sku,name,description,category,ean,upc,unit_cost_cents,unit_price_cents,currency
ELEC-LAPTOP-001,ProBook Laptop 15",15-inch business laptop,electronics,5901234123457,012345678905,65000,119900,USD
ELEC-CABLE-001,USB-C Cable 2m,Braided USB-C cable,accessories,5901234123495,012345678943,300,1499,USD
```

> **Nota**: `unit_price_cents` debe ser mayor o igual a `unit_cost_cents`.

#### `sellers`

| Columna | Tipo | Requerido | Descripcion |
|---|---|---|---|
| `external_id` | string | No | ID externo del vendedor |
| `code` | string | No | Codigo del vendedor (ej: `TD100`) |
| `name` | string | Si | Nombre del vendedor |
| `email` | string | No | Email de contacto |
| `commission_pct` | float | No | Porcentaje de comision (0-100) |

Ejemplo:
```csv
external_id,code,name,email,commission_pct
S-001,TD100,TechDirect,sales@techdirect.com,8.0
S-002,GH200,GadgetHub,contact@gadgethub.com,10.0
```

#### `customers`

| Columna | Tipo | Requerido | Descripcion |
|---|---|---|---|
| `external_id` | string | No | ID externo del cliente |
| `name` | string | Si | Nombre del cliente |
| `email` | string | No | Email de contacto |
| `segment` | string | No | Segmento: `enterprise`, `mid-market`, `smb`, `consumer` |

Ejemplo:
```csv
external_id,name,email,segment
C-001,Acme Corporation,procurement@acme.com,enterprise
C-002,Jane Smith,jane.smith@email.com,consumer
```

#### `orders`

| Columna | Tipo | Requerido | Descripcion |
|---|---|---|---|
| `external_id` | string | Si | ID externo de la orden |
| `seller_id` | UUID | Si | ID del vendedor en el sistema |
| `customer_id` | UUID | Si | ID del cliente en el sistema |
| `currency` | string | No | Moneda (default: `USD`) |

Ejemplo:
```csv
external_id,seller_id,customer_id,currency
ORD-001,550e8400-e29b-41d4-a716-446655440000,660e8400-e29b-41d4-a716-446655440001,USD
ORD-002,550e8400-e29b-41d4-a716-446655440000,660e8400-e29b-41d4-a716-446655440002,USD
```

#### `cost_records`

| Columna | Tipo | Requerido | Descripcion |
|---|---|---|---|
| `cloud_account_id` | UUID | Si | ID de la cuenta cloud |
| `service` | string | Si | Nombre del servicio (ej: `Amazon EC2`, `Cloud SQL`) |
| `amount_cents` | int | Si | Monto en centavos |
| `currency` | string | No | Moneda (default: `USD`) |
| `usage_date` | date | Si | Fecha de uso (`YYYY-MM-DD`) |

Ejemplo:
```csv
cloud_account_id,service,amount_cents,currency,usage_date
550e8400-e29b-41d4-a716-446655440000,Amazon EC2,1500000,USD,2026-02-15
550e8400-e29b-41d4-a716-446655440000,Amazon S3,110000,USD,2026-02-15
```

> La categoria (`compute`, `storage`, `network`, etc.) se asigna automaticamente segun el nombre del servicio.

### 11.4 Consultar Estado de una Importacion

```bash
curl http://localhost:8080/api/v1/imports/{id} \
  -H "Authorization: Bearer $TOKEN"
```

Respuesta:
```json
{
  "id": "...",
  "name": "Importacion Febrero 2026",
  "source": "csv_upload",
  "entity_type": "products",
  "status": "completed",
  "total_rows": 150,
  "processed_rows": 150,
  "failed_rows": 0,
  "error_log": [],
  "started_at": "2026-02-11T10:00:01Z",
  "completed_at": "2026-02-11T10:00:05Z",
  "created_at": "2026-02-11T10:00:00Z"
}
```

**Estados posibles**:
- `pending` — Job creado, aun no comienza a procesar
- `processing` — Procesando filas del CSV
- `completed` — Todas las filas procesadas
- `failed` — Error critico (ej: CSV malformado)

### 11.5 Listar Importaciones

```bash
curl "http://localhost:8080/api/v1/imports?status=completed&entity_type=products&page=1&page_size=20" \
  -H "Authorization: Bearer $TOKEN"
```

## 12. Especificacion OpenAPI

La especificacion completa de la API esta disponible en `docs/swagger.yaml`.

Para visualizarla interactivamente, se puede usar:
- [Swagger Editor](https://editor.swagger.io/) — pegar el contenido del YAML
- [Swagger UI](https://petstore.swagger.io/) — apuntar a la URL del archivo

## 13. Limites y Consideraciones

- **Paginacion**: Maximo 100 registros por pagina (`page_size` max: 100)
- **Montos**: Siempre en centavos (int64). No se aceptan montos negativos.
- **Fechas**: Formato `YYYY-MM-DD` para fechas, ISO 8601 para timestamps
- **Moneda**: Default `USD`. Se valida consistencia de moneda en operaciones aritmeticas
- **Presupuestos**: La alerta se genera al momento del `RecordSpend`, no retroactivamente
- **Anomalias**: La severidad se calcula una sola vez al momento de la deteccion
