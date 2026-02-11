import { useQuery } from '@tanstack/react-query';
import { ordersApi } from '../api/orders';

export function useOrders(params?: { seller_id?: string; customer_id?: string; status?: string; start_date?: string; end_date?: string; page?: number; page_size?: number }) {
  return useQuery({
    queryKey: ['orders', params],
    queryFn: () => ordersApi.list(params),
  });
}

export function useOrder(id: string) {
  return useQuery({
    queryKey: ['order', id],
    queryFn: () => ordersApi.getById(id),
    enabled: !!id,
  });
}
