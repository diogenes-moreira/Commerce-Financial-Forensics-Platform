import apiClient from './client';
import type { CostRecord, PaginatedResponse } from '../types';

export interface CostRecordFilters {
  cloud_account_id?: string;
  service?: string;
  category?: string;
  start_date?: string;
  end_date?: string;
  page?: number;
  page_size?: number;
}

export const costsApi = {
  list: (filters: CostRecordFilters = {}) =>
    apiClient.get<PaginatedResponse<CostRecord>>('/cost-records', { params: filters }).then((r) => r.data),
  getById: (id: string) => apiClient.get<CostRecord>(`/cost-records/${id}`).then((r) => r.data),
  create: (data: { cloud_account_id: string; service: string; amount_cents: number; currency?: string; usage_date: string }) =>
    apiClient.post<CostRecord>('/cost-records', data).then((r) => r.data),
};
