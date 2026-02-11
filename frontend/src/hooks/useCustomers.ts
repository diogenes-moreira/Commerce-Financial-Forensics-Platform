import { useQuery } from '@tanstack/react-query';
import { customersApi } from '../api/customers';

export function useCustomers(params?: { segment?: string; page?: number; page_size?: number }) {
  return useQuery({
    queryKey: ['customers', params],
    queryFn: () => customersApi.list(params),
  });
}

export function useCustomer(id: string) {
  return useQuery({
    queryKey: ['customer', id],
    queryFn: () => customersApi.getById(id),
    enabled: !!id,
  });
}
