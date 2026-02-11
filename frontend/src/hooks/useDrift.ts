import { useQuery } from '@tanstack/react-query';
import { driftApi } from '../api/drift';

export function useDrift(params?: { severity?: string; date_from?: string; date_to?: string; page?: number; page_size?: number }) {
  return useQuery({
    queryKey: ['drift', params],
    queryFn: () => driftApi.list(params),
  });
}

export function useDriftForOrder(orderId: string) {
  return useQuery({
    queryKey: ['drift', orderId],
    queryFn: () => driftApi.getForOrder(orderId),
    enabled: !!orderId,
  });
}

export function useDriftReport(params?: { date_from?: string; date_to?: string }) {
  return useQuery({
    queryKey: ['drift-report', params],
    queryFn: () => driftApi.report(params),
  });
}
