import apiClient from './client';
import type { LedgerEntry, PaginatedResponse } from '../types';

export const ledgerApi = {
  list: (params?: { account_code?: string; date_from?: string; date_to?: string; page?: number; page_size?: number }) =>
    apiClient.get<PaginatedResponse<LedgerEntry>>('/ledger', { params }).then((r) => r.data),
  getById: (id: string) => apiClient.get<LedgerEntry>(`/ledger/${id}`).then((r) => r.data),
  getBalance: (params?: { account_code?: string; as_of?: string }) =>
    apiClient.get<{ account_code: string; balance_cents: number; currency: string }>('/ledger/balance', { params }).then((r) => r.data),
};
