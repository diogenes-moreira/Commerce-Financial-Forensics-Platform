import apiClient from './client';
import type { Seller, PaginatedResponse } from '../types';

export const sellersApi = {
  list: (params?: { status?: string; page?: number; page_size?: number }) =>
    apiClient.get<PaginatedResponse<Seller>>('/sellers', { params }).then((r) => r.data),
  getById: (id: string) => apiClient.get<Seller>(`/sellers/${id}`).then((r) => r.data),
  create: (data: { external_id?: string; code?: string; name: string; email?: string; commission_pct?: number }) =>
    apiClient.post<Seller>('/sellers', data).then((r) => r.data),
  update: (id: string, data: { name: string; email?: string; commission_pct?: number }) =>
    apiClient.put<Seller>(`/sellers/${id}`, data).then((r) => r.data),
};
