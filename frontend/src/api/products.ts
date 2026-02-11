import apiClient from './client';
import type { Product, PaginatedResponse } from '../types';

export const productsApi = {
  list: (params?: { category?: string; status?: string; page?: number; page_size?: number }) =>
    apiClient.get<PaginatedResponse<Product>>('/products', { params }).then((r) => r.data),
  getById: (id: string) => apiClient.get<Product>(`/products/${id}`).then((r) => r.data),
  create: (data: { sku: string; ean?: string; upc?: string; name: string; description?: string; category?: string; unit_cost_cents: number; unit_price_cents: number; currency?: string }) =>
    apiClient.post<Product>('/products', data).then((r) => r.data),
  update: (id: string, data: { name: string; description?: string; category?: string; unit_price_cents?: number; currency?: string }) =>
    apiClient.put<Product>(`/products/${id}`, data).then((r) => r.data),
};
