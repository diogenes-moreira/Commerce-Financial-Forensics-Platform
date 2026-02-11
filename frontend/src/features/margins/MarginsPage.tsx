import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import CircularProgress from '@mui/material/CircularProgress';
import Chip from '@mui/material/Chip';
import DataTable, { type Column } from '../../components/common/DataTable';
import { useMargins } from '../../hooks/useMargins';
import { formatCents } from '../../utils/formatCurrency';
import type { MarginBreakdown } from '../../types';

function marginColor(pct: number): 'success' | 'warning' | 'error' {
  if (pct >= 20) return 'success';
  if (pct >= 5) return 'warning';
  return 'error';
}

const columns: Column<MarginBreakdown>[] = [
  { key: 'order_id', header: 'Order ID', render: (r) => r.order_id.slice(0, 8) + '...' },
  { key: 'revenue', header: 'Revenue', render: (r) => formatCents(r.revenue_cents, r.currency) },
  { key: 'cogs', header: 'COGS', render: (r) => formatCents(r.cogs_cents, r.currency) },
  { key: 'discount', header: 'Discount', render: (r) => formatCents(r.discount_cents, r.currency) },
  { key: 'commission', header: 'Commission', render: (r) => formatCents(r.commission_cents, r.currency) },
  { key: 'gross_margin', header: 'Gross Margin', render: (r) => formatCents(r.gross_margin_cents, r.currency) },
  { key: 'net_margin', header: 'Net Margin', render: (r) => formatCents(r.net_margin_cents, r.currency) },
  {
    key: 'real_margin_pct',
    header: 'Real Margin %',
    render: (r) => (
      <Chip
        label={`${r.real_margin_pct.toFixed(1)}%`}
        size="small"
        color={marginColor(r.real_margin_pct)}
      />
    ),
  },
];

export default function MarginsPage() {
  const { data, isLoading } = useMargins();

  if (isLoading) return <CircularProgress />;

  return (
    <Box>
      <Typography variant="h4" gutterBottom>Margins</Typography>
      <DataTable columns={columns} data={data?.items ?? []} keyExtractor={(r) => r.order_id} />
    </Box>
  );
}
