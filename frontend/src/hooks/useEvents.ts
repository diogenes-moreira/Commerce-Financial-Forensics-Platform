import { useQuery } from '@tanstack/react-query';
import { eventsApi } from '../api/events';

export function useEvents(params?: { entity_type?: string; entity_id?: string; event_type?: string; date_from?: string; date_to?: string; page?: number; page_size?: number }) {
  return useQuery({
    queryKey: ['events', params],
    queryFn: () => eventsApi.list(params),
  });
}

export function useEvent(id: string) {
  return useQuery({
    queryKey: ['event', id],
    queryFn: () => eventsApi.getById(id),
    enabled: !!id,
  });
}
