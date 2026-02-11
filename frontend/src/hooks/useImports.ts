import { useQuery } from '@tanstack/react-query';
import { importsApi } from '../api/imports';

export function useImports(params?: { status?: string; entity_type?: string; page?: number; page_size?: number }) {
  return useQuery({
    queryKey: ['imports', params],
    queryFn: () => importsApi.list(params),
  });
}

export function useImportJob(id: string) {
  return useQuery({
    queryKey: ['import', id],
    queryFn: () => importsApi.getById(id),
    enabled: !!id,
  });
}
