import { useState } from 'react';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import CircularProgress from '@mui/material/CircularProgress';
import Chip from '@mui/material/Chip';
import ToggleButton from '@mui/material/ToggleButton';
import ToggleButtonGroup from '@mui/material/ToggleButtonGroup';
import DataTable, { type Column } from '../../components/common/DataTable';
import { usePayments } from '../../hooks/usePayments';
import { formatCents } from '../../utils/formatCurrency';
import type { Payment } from '../../types';

const statusColorMap: Record<string, 'success' | 'warning' | 'error' | 'info' | 'default'> = {
  completed: 'success',
  pending: 'warning',
  failed: 'error',
  processing: 'info',
  refunded: 'default',
};

const columns: Column<Payment>[] = [
  {
    key: 'direction',
    header: 'Direction',
    render: (r) => (
      <Chip
        label={r.direction === 'inbound' ? 'From Customer' : 'To Seller'}
        size="small"
        color={r.direction === 'inbound' ? 'info' : 'secondary'}
      />
    ),
  },
  { key: 'order_id', header: 'Order', render: (r) => r.order_id.slice(0, 8) + '...' },
  { key: 'amount', header: 'Amount', render: (r) => formatCents(r.amount_cents, r.currency) },
  { key: 'method', header: 'Method', render: (r) => r.method },
  {
    key: 'status',
    header: 'Status',
    render: (r) => (
      <Chip
        label={r.status}
        size="small"
        color={statusColorMap[r.status] ?? 'default'}
      />
    ),
  },
  { key: 'external_ref', header: 'External Ref', render: (r) => r.external_ref || '---' },
  {
    key: 'processed_at',
    header: 'Processed At',
    render: (r) => r.processed_at ? new Date(r.processed_at).toLocaleString() : '---',
  },
];

export default function PaymentsPage() {
  const [direction, setDirection] = useState<string | null>(null);

  const { data, isLoading } = usePayments({
    direction: direction || undefined,
  });

  if (isLoading) return <CircularProgress />;

  return (
    <Box>
      <Typography variant="h4" gutterBottom>Payments</Typography>

      <ToggleButtonGroup
        value={direction}
        exclusive
        onChange={(_, val) => setDirection(val)}
        size="small"
        sx={{ mb: 3 }}
      >
        <ToggleButton value={null as any}>All</ToggleButton>
        <ToggleButton value="inbound">From Customer</ToggleButton>
        <ToggleButton value="outbound">To Seller</ToggleButton>
      </ToggleButtonGroup>

      <DataTable columns={columns} data={data?.items ?? []} keyExtractor={(r) => r.id} />
    </Box>
  );
}
