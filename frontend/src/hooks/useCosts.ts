import { useQuery } from '@tanstack/react-query';
import { costsApi, type CostRecordFilters } from '../api/costs';

export function useCosts(filters: CostRecordFilters = {}) {
  return useQuery({
    queryKey: ['costs', filters],
    queryFn: () => costsApi.list(filters),
  });
}

export function useCost(id: string) {
  return useQuery({
    queryKey: ['cost', id],
    queryFn: () => costsApi.getById(id),
    enabled: !!id,
  });
}
