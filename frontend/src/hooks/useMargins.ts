import { useQuery } from '@tanstack/react-query';
import { marginsApi } from '../api/margins';

export function useMargins(params?: { date_from?: string; date_to?: string; page?: number; page_size?: number }) {
  return useQuery({
    queryKey: ['margins', params],
    queryFn: () => marginsApi.list(params),
  });
}

export function useMarginForOrder(orderId: string) {
  return useQuery({
    queryKey: ['margin', orderId],
    queryFn: () => marginsApi.getForOrder(orderId),
    enabled: !!orderId,
  });
}
