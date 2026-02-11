import apiClient from './client';
import type { MarginBreakdown, PaginatedResponse } from '../types';

export const marginsApi = {
  list: (params?: { date_from?: string; date_to?: string; page?: number; page_size?: number }) =>
    apiClient.get<PaginatedResponse<MarginBreakdown>>('/margins', { params }).then((r) => r.data),
  getForOrder: (orderId: string) =>
    apiClient.get<MarginBreakdown>(`/margins/order/${orderId}`).then((r) => r.data),
};
