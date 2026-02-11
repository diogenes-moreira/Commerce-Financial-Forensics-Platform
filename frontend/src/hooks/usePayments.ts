import { useQuery } from '@tanstack/react-query';
import { paymentsApi } from '../api/payments';

export function usePayments(params?: { direction?: string; status?: string; page?: number; page_size?: number }) {
  return useQuery({
    queryKey: ['payments', params],
    queryFn: () => paymentsApi.list(params),
  });
}

export function usePayment(id: string) {
  return useQuery({
    queryKey: ['payment', id],
    queryFn: () => paymentsApi.getById(id),
    enabled: !!id,
  });
}
