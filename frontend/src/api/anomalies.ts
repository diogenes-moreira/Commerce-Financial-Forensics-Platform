import apiClient from './client';
import type { Anomaly } from '../types';

export const anomaliesApi = {
  list: () => apiClient.get<Anomaly[]>('/anomalies').then((r) => r.data),
  getById: (id: string) => apiClient.get<Anomaly>(`/anomalies/${id}`).then((r) => r.data),
  resolve: (id: string) => apiClient.post<Anomaly>(`/anomalies/${id}/resolve`).then((r) => r.data),
};
