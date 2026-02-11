import apiClient from './client';
import type { Integration, PaginatedResponse } from '../types';

export const integrationsApi = {
  list: (params?: { platform?: string; status?: string; page?: number; page_size?: number }) =>
    apiClient.get<PaginatedResponse<Integration>>('/integrations', { params }).then((r) => r.data),
  getById: (id: string) => apiClient.get<Integration>(`/integrations/${id}`).then((r) => r.data),
  create: (data: { platform: string; name: string; config: Record<string, string> }) =>
    apiClient.post<Integration>('/integrations', data).then((r) => r.data),
  update: (id: string, data: { name?: string; config?: Record<string, string> }) =>
    apiClient.put<Integration>(`/integrations/${id}`, data).then((r) => r.data),
  delete: (id: string) => apiClient.delete(`/integrations/${id}`).then((r) => r.data),
  triggerSync: (id: string) =>
    apiClient.post(`/integrations/${id}/sync`).then((r) => r.data),
};
