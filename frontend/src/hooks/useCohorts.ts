import { useQuery } from '@tanstack/react-query';
import { cohortsApi } from '../api/cohorts';

export function useRetentionCohorts(params?: { months?: number }) {
  return useQuery({
    queryKey: ['cohorts-retention', params],
    queryFn: () => cohortsApi.getRetention(params),
  });
}
