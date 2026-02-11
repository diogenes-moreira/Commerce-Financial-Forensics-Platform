import apiClient from './client';
import type { ForensicEvent, PaginatedResponse } from '../types';

export const eventsApi = {
  list: (params?: { entity_type?: string; entity_id?: string; event_type?: string; date_from?: string; date_to?: string; page?: number; page_size?: number }) =>
    apiClient.get<PaginatedResponse<ForensicEvent>>('/events', { params }).then((r) => r.data),
  getById: (id: string) => apiClient.get<ForensicEvent>(`/events/${id}`).then((r) => r.data),
};
