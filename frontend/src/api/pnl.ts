import apiClient from './client';
import type { PLReport } from '../types';

export const pnlApi = {
  get: (params?: { granularity?: string; date_from?: string; date_to?: string }) =>
    apiClient.get<PLReport[]>('/pnl', { params }).then((r) => r.data),
};
