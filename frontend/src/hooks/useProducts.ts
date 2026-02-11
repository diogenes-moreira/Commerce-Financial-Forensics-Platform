import { useQuery } from '@tanstack/react-query';
import { productsApi } from '../api/products';

export function useProducts(params?: { category?: string; status?: string; page?: number; page_size?: number }) {
  return useQuery({
    queryKey: ['products', params],
    queryFn: () => productsApi.list(params),
  });
}

export function useProduct(id: string) {
  return useQuery({
    queryKey: ['product', id],
    queryFn: () => productsApi.getById(id),
    enabled: !!id,
  });
}
