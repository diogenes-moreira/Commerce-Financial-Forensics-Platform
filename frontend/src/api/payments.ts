import apiClient from './client';
import type { Payment, PaginatedResponse } from '../types';

export const paymentsApi = {
  list: (params?: { direction?: string; status?: string; page?: number; page_size?: number }) =>
    apiClient.get<PaginatedResponse<Payment>>('/payments', { params }).then((r) => r.data),
  getById: (id: string) => apiClient.get<Payment>(`/payments/${id}`).then((r) => r.data),
  create: (data: { order_id: string; direction: string; counterparty_id: string; amount_cents: number; currency?: string; method: string; external_ref?: string }) =>
    apiClient.post<Payment>('/payments', data).then((r) => r.data),
  updateStatus: (id: string, data: { status: string; failure_reason?: string }) =>
    apiClient.patch<Payment>(`/payments/${id}/status`, data).then((r) => r.data),
};
