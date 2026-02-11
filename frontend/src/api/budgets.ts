import apiClient from './client';
import type { Budget, BudgetAlert } from '../types';

export const budgetsApi = {
  list: () => apiClient.get<Budget[]>('/budgets').then((r) => r.data),
  getById: (id: string) => apiClient.get<Budget>(`/budgets/${id}`).then((r) => r.data),
  create: (data: { name: string; amount_cents: number; currency?: string; period_start: string; period_end: string; alert_threshold_pct?: number }) =>
    apiClient.post<Budget>('/budgets', data).then((r) => r.data),
  update: (id: string, data: { name: string }) =>
    apiClient.put<Budget>(`/budgets/${id}`, data).then((r) => r.data),
  recordSpend: (id: string, data: { amount_cents: number; currency?: string }) =>
    apiClient.post<{ budget: Budget; alert: BudgetAlert | null }>(`/budgets/${id}/spend`, data).then((r) => r.data),
};
