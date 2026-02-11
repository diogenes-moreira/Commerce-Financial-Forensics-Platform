import { useQuery } from '@tanstack/react-query';
import { exchangeRatesApi } from '../api/exchangeRates';

export function useExchangeRates(params?: { base_currency?: string; quote_currency?: string; page?: number; page_size?: number }) {
  return useQuery({
    queryKey: ['exchange-rates', params],
    queryFn: () => exchangeRatesApi.list(params),
  });
}

export function useExchangeRate(id: string) {
  return useQuery({
    queryKey: ['exchange-rate', id],
    queryFn: () => exchangeRatesApi.getById(id),
    enabled: !!id,
  });
}

export function useLatestExchangeRate(base: string, quote: string) {
  return useQuery({
    queryKey: ['exchange-rate-latest', base, quote],
    queryFn: () => exchangeRatesApi.getLatest(base, quote),
    enabled: !!base && !!quote,
  });
}
