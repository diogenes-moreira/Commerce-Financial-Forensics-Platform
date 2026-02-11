import apiClient from './client';
import type { PromotionRule, DiscountApplication, PaginatedResponse } from '../types';

export const discountsApi = {
  listRules: (params?: { status?: string; page?: number; page_size?: number }) =>
    apiClient.get<PaginatedResponse<PromotionRule>>('/promotion-rules', { params }).then((r) => r.data),
  getRuleById: (id: string) => apiClient.get<PromotionRule>(`/promotion-rules/${id}`).then((r) => r.data),
  createRule: (data: { name: string; rule_type: string; value: number; conditions?: Record<string, any>; funding_source: string; funding_pct: number; max_usage_count?: number; valid_from: string; valid_to: string }) =>
    apiClient.post<PromotionRule>('/promotion-rules', data).then((r) => r.data),
  updateRule: (id: string, data: { name?: string; value?: number; conditions?: Record<string, any>; max_usage_count?: number; valid_to?: string; status?: string }) =>
    apiClient.put<PromotionRule>(`/promotion-rules/${id}`, data).then((r) => r.data),
  listApplications: (params?: { order_id?: string; promotion_rule_id?: string; page?: number; page_size?: number }) =>
    apiClient.get<PaginatedResponse<DiscountApplication>>('/discount-applications', { params }).then((r) => r.data),
};
