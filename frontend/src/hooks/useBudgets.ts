import { useQuery } from '@tanstack/react-query';
import { budgetsApi } from '../api/budgets';

export function useBudgets() {
  return useQuery({
    queryKey: ['budgets'],
    queryFn: budgetsApi.list,
  });
}

export function useBudget(id: string) {
  return useQuery({
    queryKey: ['budget', id],
    queryFn: () => budgetsApi.getById(id),
    enabled: !!id,
  });
}
