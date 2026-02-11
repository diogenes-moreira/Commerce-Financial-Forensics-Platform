import { useQuery } from '@tanstack/react-query';
import { ledgerApi } from '../api/ledger';

export function useLedger(params?: { account_code?: string; date_from?: string; date_to?: string; page?: number; page_size?: number }) {
  return useQuery({
    queryKey: ['ledger', params],
    queryFn: () => ledgerApi.list(params),
  });
}

export function useLedgerEntry(id: string) {
  return useQuery({
    queryKey: ['ledger', id],
    queryFn: () => ledgerApi.getById(id),
    enabled: !!id,
  });
}

export function useLedgerBalance(params?: { account_code?: string; as_of?: string }) {
  return useQuery({
    queryKey: ['ledger-balance', params],
    queryFn: () => ledgerApi.getBalance(params),
  });
}
