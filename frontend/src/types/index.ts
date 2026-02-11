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
