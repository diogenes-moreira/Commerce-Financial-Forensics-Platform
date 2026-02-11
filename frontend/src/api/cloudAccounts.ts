import apiClient from './client';
import type { CloudAccount } from '../types';

export const cloudAccountsApi = {
  list: () => apiClient.get<CloudAccount[]>('/cloud-accounts').then((r) => r.data),
  getById: (id: string) => apiClient.get<CloudAccount>(`/cloud-accounts/${id}`).then((r) => r.data),
  create: (data: { provider: string; name: string; external_id: string }) =>
    apiClient.post<CloudAccount>('/cloud-accounts', data).then((r) => r.data),
  update: (id: string, data: { name: string }) =>
    apiClient.put<CloudAccount>(`/cloud-accounts/${id}`, data).then((r) => r.data),
  delete: (id: string) => apiClient.delete(`/cloud-accounts/${id}`),
  beginSync: (id: string) =>
    apiClient.post<CloudAccount>(`/cloud-accounts/${id}/sync`).then((r) => r.data),
};
