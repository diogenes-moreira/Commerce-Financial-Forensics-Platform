import apiClient from './client';
import type { RetentionCohort } from '../types';

export const cohortsApi = {
  getRetention: (params?: { months?: number }) =>
    apiClient.get<RetentionCohort[]>('/cohorts/retention', { params }).then((r) => r.data),
};
