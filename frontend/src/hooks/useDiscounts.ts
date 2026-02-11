import { useQuery } from '@tanstack/react-query';
import { discountsApi } from '../api/discounts';

export function usePromotionRules(params?: { status?: string; page?: number; page_size?: number }) {
  return useQuery({
    queryKey: ['promotion-rules', params],
    queryFn: () => discountsApi.listRules(params),
  });
}

export function usePromotionRule(id: string) {
  return useQuery({
    queryKey: ['promotion-rule', id],
    queryFn: () => discountsApi.getRuleById(id),
    enabled: !!id,
  });
}

export function useDiscountApplications(params?: { order_id?: string; promotion_rule_id?: string; page?: number; page_size?: number }) {
  return useQuery({
    queryKey: ['discount-applications', params],
    queryFn: () => discountsApi.listApplications(params),
  });
}
