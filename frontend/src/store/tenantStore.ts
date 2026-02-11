import { create } from 'zustand';
import type { Tenant } from '../types';

interface TenantState {
  currentTenant: Tenant | null;
  setTenant: (tenant: Tenant) => void;
}

export const useTenantStore = create<TenantState>((set) => ({
  currentTenant: null,
  setTenant: (tenant: Tenant) => set({ currentTenant: tenant }),
}));
