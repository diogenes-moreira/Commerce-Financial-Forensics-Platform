import { useQuery } from '@tanstack/react-query';
import { integrationsApi } from '../api/integrations';

export function useIntegrations(params?: { platform?: string; status?: string; page?: number; page_size?: number }) {
  return useQuery({
    queryKey: ['integrations', params],
    queryFn: () => integrationsApi.list(params),
  });
}

export function useIntegration(id: string) {
  return useQuery({
    queryKey: ['integration', id],
    queryFn: () => integrationsApi.getById(id),
    enabled: !!id,
  });
}
