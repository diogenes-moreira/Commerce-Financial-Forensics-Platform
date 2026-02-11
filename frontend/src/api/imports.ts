import apiClient from './client';
import type { ImportJob, PaginatedResponse } from '../types';

export const importsApi = {
  list: (params?: { status?: string; entity_type?: string; page?: number; page_size?: number }) =>
    apiClient.get<PaginatedResponse<ImportJob>>('/imports', { params }).then((r) => r.data),
  getById: (id: string) => apiClient.get<ImportJob>(`/imports/${id}`).then((r) => r.data),
  upload: (file: File, entityType: string, name?: string) => {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('entity_type', entityType);
    if (name) formData.append('name', name);
    return apiClient.post<ImportJob>('/imports/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    }).then((r) => r.data);
  },
};
