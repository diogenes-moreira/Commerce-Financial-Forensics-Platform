import apiClient from './client';
import type { Customer, PaginatedResponse } from '../types';

export const customersApi = {
  list: (params?: { segment?: string; page?: number; page_size?: number }) =>
    apiClient.get<PaginatedResponse<Customer>>('/customers', { params }).then((r) => r.data),
  getById: (id: string) => apiClient.get<Customer>(`/customers/${id}`).then((r) => r.data),
  create: (data: { external_id?: string; name: string; email?: string; segment?: string }) =>
    apiClient.post<Customer>('/customers', data).then((r) => r.data),
  update: (id: string, data: { name: string; email?: string; segment?: string }) =>
    apiClient.put<Customer>(`/customers/${id}`, data).then((r) => r.data),
};
