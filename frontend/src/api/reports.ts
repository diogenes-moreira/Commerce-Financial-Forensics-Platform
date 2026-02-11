import apiClient from './client';
import type { CostReport } from '../types';

export const reportsApi = {
  list: () => apiClient.get<CostReport[]>('/cost-reports').then((r) => r.data),
  getById: (id: string) => apiClient.get<CostReport>(`/cost-reports/${id}`).then((r) => r.data),
  generate: (data: { name: string; report_type: string; period_start: string; period_end: string; currency?: string }) =>
    apiClient.post<CostReport>('/cost-reports', data).then((r) => r.data),
};
