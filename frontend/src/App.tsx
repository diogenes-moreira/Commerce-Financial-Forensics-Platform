import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import theme from './theme';
import AppLayout from './components/layout/AppLayout';
import DashboardPage from './features/dashboard/DashboardPage';
import CloudAccountsPage from './features/cloud-accounts/CloudAccountsPage';
import CostsPage from './features/costs/CostsPage';
import BudgetsPage from './features/budgets/BudgetsPage';
import AnomaliesPage from './features/anomalies/AnomaliesPage';
import ReportsPage from './features/reports/ReportsPage';
import SettingsPage from './features/settings/SettingsPage';
import ProductsPage from './features/products/ProductsPage';
import SellersPage from './features/sellers/SellersPage';
import CustomersPage from './features/customers/CustomersPage';
import OrdersPage from './features/orders/OrdersPage';
import LedgerPage from './features/ledger/LedgerPage';
import EventsPage from './features/events/EventsPage';
import DiscountsPage from './features/discounts/DiscountsPage';
import PaymentsPage from './features/payments/PaymentsPage';
import ExchangeRatesPage from './features/exchange-rates/ExchangeRatesPage';
import MarginsPage from './features/margins/MarginsPage';
import DriftPage from './features/drift/DriftPage';
import PnLPage from './features/pnl/PnLPage';
import CohortsPage from './features/cohorts/CohortsPage';
import ImportsPage from './features/imports/ImportsPage';
import IntegrationsPage from './features/integrations/IntegrationsPage';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 30_000,
      retry: 1,
    },
  },
});

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <ThemeProvider theme={theme}>
        <CssBaseline />
        <BrowserRouter>
          <Routes>
            <Route element={<AppLayout />}>
              <Route path="/" element={<DashboardPage />} />
              <Route path="/cloud-accounts" element={<CloudAccountsPage />} />
              <Route path="/costs" element={<CostsPage />} />
              <Route path="/budgets" element={<BudgetsPage />} />
              <Route path="/anomalies" element={<AnomaliesPage />} />
              <Route path="/reports" element={<ReportsPage />} />
              <Route path="/products" element={<ProductsPage />} />
              <Route path="/sellers" element={<SellersPage />} />
              <Route path="/customers" element={<CustomersPage />} />
              <Route path="/orders" element={<OrdersPage />} />
              <Route path="/ledger" element={<LedgerPage />} />
              <Route path="/events" element={<EventsPage />} />
              <Route path="/discounts" element={<DiscountsPage />} />
              <Route path="/payments" element={<PaymentsPage />} />
              <Route path="/exchange-rates" element={<ExchangeRatesPage />} />
              <Route path="/margins" element={<MarginsPage />} />
              <Route path="/drift" element={<DriftPage />} />
              <Route path="/pnl" element={<PnLPage />} />
              <Route path="/cohorts" element={<CohortsPage />} />
              <Route path="/imports" element={<ImportsPage />} />
              <Route path="/integrations" element={<IntegrationsPage />} />
              <Route path="/settings" element={<SettingsPage />} />
            </Route>
          </Routes>
        </BrowserRouter>
      </ThemeProvider>
    </QueryClientProvider>
  );
}
