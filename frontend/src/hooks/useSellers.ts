import { useQuery } from '@tanstack/react-query';
import { sellersApi } from '../api/sellers';

export function useSellers(params?: { status?: string; page?: number; page_size?: number }) {
  return useQuery({
    queryKey: ['sellers', params],
    queryFn: () => sellersApi.list(params),
  });
}

export function useSeller(id: string) {
  return useQuery({
    queryKey: ['seller', id],
    queryFn: () => sellersApi.getById(id),
    enabled: !!id,
  });
}
