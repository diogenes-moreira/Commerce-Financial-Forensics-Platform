import apiClient from './client';
import type { Order, PaginatedResponse } from '../types';

export const ordersApi = {
  list: (params?: { seller_id?: string; customer_id?: string; status?: string; start_date?: string; end_date?: string; page?: number; page_size?: number }) =>
    apiClient.get<PaginatedResponse<Order>>('/orders', { params }).then((r) => r.data),
  getById: (id: string) => apiClient.get<Order>(`/orders/${id}`).then((r) => r.data),
  create: (data: { external_id: string; seller_id: string; customer_id: string; currency?: string }) =>
    apiClient.post<Order>('/orders', data).then((r) => r.data),
  addItem: (orderId: string, data: { product_id: string; sku: string; product_name: string; quantity: number; unit_price_cents: number; unit_cost_cents: number; currency?: string }) =>
    apiClient.post<Order>(`/orders/${orderId}/items`, data).then((r) => r.data),
  confirm: (id: string) => apiClient.post<Order>(`/orders/${id}/confirm`).then((r) => r.data),
  ship: (id: string) => apiClient.post<Order>(`/orders/${id}/ship`).then((r) => r.data),
  deliver: (id: string) => apiClient.post<Order>(`/orders/${id}/deliver`).then((r) => r.data),
  cancel: (id: string) => apiClient.post<Order>(`/orders/${id}/cancel`).then((r) => r.data),
  refund: (id: string) => apiClient.post<Order>(`/orders/${id}/refund`).then((r) => r.data),
};
