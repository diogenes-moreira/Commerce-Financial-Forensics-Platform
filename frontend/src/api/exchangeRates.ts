import apiClient from './client';
import type { ExchangeRate, PaginatedResponse } from '../types';

export const exchangeRatesApi = {
  list: (params?: { base_currency?: string; quote_currency?: string; page?: number; page_size?: number }) =>
    apiClient.get<PaginatedResponse<ExchangeRate>>('/exchange-rates', { params }).then((r) => r.data),
  getById: (id: string) => apiClient.get<ExchangeRate>(`/exchange-rates/${id}`).then((r) => r.data),
  create: (data: { base_currency: string; quote_currency: string; rate: number; source?: string; effective_date: string }) =>
    apiClient.post<ExchangeRate>('/exchange-rates', data).then((r) => r.data),
  update: (id: string, data: { rate: number; source?: string }) =>
    apiClient.put<ExchangeRate>(`/exchange-rates/${id}`, data).then((r) => r.data),
  getLatest: (base: string, quote: string) =>
    apiClient.get<ExchangeRate>(`/exchange-rates/latest`, { params: { base_currency: base, quote_currency: quote } }).then((r) => r.data),
};
