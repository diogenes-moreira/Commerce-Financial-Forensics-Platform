import { useQuery } from '@tanstack/react-query';
import { anomaliesApi } from '../api/anomalies';

export function useAnomalies() {
  return useQuery({
    queryKey: ['anomalies'],
    queryFn: anomaliesApi.list,
  });
}

export function useAnomaly(id: string) {
  return useQuery({
    queryKey: ['anomaly', id],
    queryFn: () => anomaliesApi.getById(id),
    enabled: !!id,
  });
}
