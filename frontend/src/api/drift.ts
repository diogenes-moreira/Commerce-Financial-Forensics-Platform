import apiClient from './client';
import type { DriftResult, PaginatedResponse } from '../types';

export const driftApi = {
  list: (params?: { severity?: string; date_from?: string; date_to?: string; page?: number; page_size?: number }) =>
    apiClient.get<PaginatedResponse<DriftResult>>('/drift', { params }).then((r) => r.data),
  getForOrder: (orderId: string) =>
    apiClient.get<DriftResult[]>(`/drift/order/${orderId}`).then((r) => r.data),
  report: (params?: { date_from?: string; date_to?: string }) =>
    apiClient.get<{ total_drifts: number; by_severity: Record<string, number>; by_field: Record<string, number> }>('/drift/report', { params }).then((r) => r.data),
};
