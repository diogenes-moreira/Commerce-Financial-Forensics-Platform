import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import CircularProgress from '@mui/material/CircularProgress';
import Chip from '@mui/material/Chip';
import DataTable, { type Column } from '../../components/common/DataTable';
import { useOrders } from '../../hooks/useOrders';
import { formatCents } from '../../utils/formatCurrency';
import type { Order } from '../../types';

const statusColors: Record<string, 'default' | 'primary' | 'secondary' | 'error' | 'info' | 'success' | 'warning'> = {
  pending: 'warning',
  confirmed: 'info',
  shipped: 'primary',
  delivered: 'success',
  cancelled: 'error',
  refunded: 'secondary',
};

const columns: Column<Order>[] = [
  { key: 'external_id', header: 'Order ID', render: (r) => r.external_id },
  { key: 'total', header: 'Total', render: (r) => formatCents(r.total_cents, r.currency) },
  { key: 'items', header: 'Items', render: (r) => `${r.items?.length ?? 0}` },
  {
    key: 'status',
    header: 'Status',
    render: (r) => <Chip label={r.status} size="small" color={statusColors[r.status] ?? 'default'} />,
  },
  { key: 'date', header: 'Date', render: (r) => r.order_date },
];

export default function OrdersPage() {
  const { data, isLoading } = useOrders();

  if (isLoading) return <CircularProgress />;

  return (
    <Box>
      <Typography variant="h4" gutterBottom>Orders</Typography>
      <DataTable columns={columns} data={data?.items ?? []} keyExtractor={(r) => r.id} />
    </Box>
  );
}
