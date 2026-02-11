export interface Tenant {
  id: string;
  name: string;
  slug: string;
  db_name: string;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface CloudAccount {
  id: string;
  provider: string;
  name: string;
  external_id: string;
  status: string;
  sync_status: string;
  last_synced_at: string | null;
  created_at: string;
  updated_at: string;
}

export interface CostRecord {
  id: string;
  cloud_account_id: string;
  service: string;
  category: string;
  amount_cents: number;
  currency: string;
  usage_date: string;
  tags: Record<string, string>;
  created_at: string;
}

export interface Budget {
  id: string;
  name: string;
  amount_cents: number;
  currency: string;
  spent_cents: number;
  usage_pct: number;
  period_start: string;
  period_end: string;
  alert_threshold_pct: number;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface BudgetAlert {
  id: string;
  budget_id: string;
  threshold_pct: number;
  actual_pct: number;
  message: string;
  triggered_at: string;
}

export interface Anomaly {
  id: string;
  cloud_account_id: string;
  service: string;
  expected_cents: number;
  actual_cents: number;
  currency: string;
  deviation_pct: number;
  severity: string;
  status: string;
  detected_at: string;
  resolved_at: string | null;
}

export interface CostReport {
  id: string;
  name: string;
  report_type: string;
  period_start: string;
  period_end: string;
  total_cents: number;
  currency: string;
  breakdown: Record<string, number>;
  generated_at: string;
}

export interface PaginatedResponse<T> {
  items: T[];
  total_count: number;
  page: number;
  page_size: number;
}

export interface Product {
  id: string;
  sku: string;
  ean: string;
  upc: string;
  name: string;
  description: string;
  category: string;
  unit_cost_cents: number;
  unit_price_cents: number;
  currency: string;
  gross_margin_pct: number;
  status: string;
  metadata: Record<string, string>;
  created_at: string;
  updated_at: string;
}

export interface Seller {
  id: string;
  external_id: string;
  code: string;
  name: string;
  email: string;
  commission_pct: number;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface Customer {
  id: string;
  external_id: string;
  name: string;
  email: string;
  segment: string;
  metadata: Record<string, string>;
  created_at: string;
  updated_at: string;
}

export interface OrderItem {
  id: string;
  order_id: string;
  product_id: string;
  sku: string;
  product_name: string;
  quantity: number;
  unit_price_cents: number;
  unit_cost_cents: number;
  discount_cents: number;
  total_cents: number;
  created_at: string;
}

export interface Order {
  id: string;
  external_id: string;
  seller_id: string;
  customer_id: string;
  status: string;
  subtotal_cents: number;
  discount_cents: number;
  shipping_cents: number;
  tax_cents: number;
  total_cents: number;
  currency: string;
  order_date: string;
  items: OrderItem[];
  metadata: Record<string, string>;
  created_at: string;
  updated_at: string;
}

// Phase 2 types

export interface LedgerEntry {
  id: string;
  order_id: string | null;
  account_code: string;
  side: string;
  amount_cents: number;
  currency: string;
  description: string;
  reference_type: string;
  reference_id: string;
  effective_date: string;
  created_at: string;
}

export interface ForensicEvent {
  id: string;
  entity_type: string;
  entity_id: string;
  event_type: string;
  event_time: string;
  effective_time: string;
  payload: Record<string, any>;
  actor_id: string;
  hash_integrity: string;
  created_at: string;
}

export interface PromotionRule {
  id: string;
  name: string;
  rule_type: string;
  value: number;
  conditions: Record<string, any>;
  funding_source: string;
  funding_pct: number;
  max_usage_count: number;
  current_usage: number;
  valid_from: string;
  valid_to: string;
  status: string;
  version: number;
  created_at: string;
  updated_at: string;
}

export interface DiscountApplication {
  id: string;
  order_id: string;
  promotion_rule_id: string;
  rule_version: number;
  discount_cents: number;
  currency: string;
  funding_source: string;
  seller_share_cents: number;
  platform_share_cents: number;
  applied_at: string;
}

export interface Payment {
  id: string;
  order_id: string;
  direction: string;
  counterparty_id: string;
  amount_cents: number;
  currency: string;
  method: string;
  external_ref: string;
  status: string;
  processed_at: string | null;
  failure_reason: string;
  metadata: Record<string, string>;
  created_at: string;
  updated_at: string;
}

export interface ExchangeRate {
  id: string;
  base_currency: string;
  quote_currency: string;
  rate: number;
  inverse_rate: number;
  source: string;
  effective_date: string;
  created_at: string;
}

// Phase 3 types

export interface MarginBreakdown {
  order_id: string;
  revenue_cents: number;
  cogs_cents: number;
  discount_cents: number;
  seller_discount_cents: number;
  platform_discount_cents: number;
  commission_cents: number;
  shipping_cost_cents: number;
  platform_fee_cents: number;
  tax_cents: number;
  gross_margin_cents: number;
  net_margin_cents: number;
  real_margin_pct: number;
  currency: string;
  calculated_at: string;
}

export interface DriftResult {
  order_id: string;
  field_name: string;
  expected_value: number;
  actual_value: number;
  drift_pct: number;
  severity: string;
  detected_at: string;
}

export interface PLReport {
  period: string;
  granularity: string;
  revenue_cents: number;
  cogs_cents: number;
  gross_profit_cents: number;
  discount_cents: number;
  commission_cents: number;
  shipping_cents: number;
  platform_fees_cents: number;
  net_profit_cents: number;
  gross_margin_pct: number;
  net_margin_pct: number;
  order_count: number;
  currency: string;
  children?: PLReport[];
}

export interface RetentionCohort {
  cohort_month: string;
  total_customers: number;
  retention: Record<string, number>;
}

// Phase 4 types

export interface ImportJob {
  id: string;
  name: string;
  source: string;
  entity_type: string;
  status: string;
  total_rows: number;
  processed_rows: number;
  failed_rows: number;
  error_log: string[];
  source_uri: string;
  started_at: string | null;
  completed_at: string | null;
  created_at: string;
}

// Phase 5 types

export interface Integration {
  id: string;
  platform: string;
  name: string;
  config: Record<string, string>;
  status: string;
  sync_status: string;
  last_synced_at: string | null;
  created_at: string;
  updated_at: string;
}
