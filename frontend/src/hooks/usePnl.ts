import { useQuery } from '@tanstack/react-query';
import { pnlApi } from '../api/pnl';

export function usePnl(params?: { granularity?: string; date_from?: string; date_to?: string }) {
  return useQuery({
    queryKey: ['pnl', params],
    queryFn: () => pnlApi.get(params),
  });
}
