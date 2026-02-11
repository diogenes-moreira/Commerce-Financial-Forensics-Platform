import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import CircularProgress from '@mui/material/CircularProgress';
import Chip from '@mui/material/Chip';
import DataTable, { type Column } from '../../components/common/DataTable';
import { useProducts } from '../../hooks/useProducts';
import { formatCents } from '../../utils/formatCurrency';
import type { Product } from '../../types';

const columns: Column<Product>[] = [
  { key: 'sku', header: 'SKU', render: (r) => r.sku },
  { key: 'ean', header: 'EAN', render: (r) => r.ean || '—' },
  { key: 'upc', header: 'UPC', render: (r) => r.upc || '—' },
  { key: 'name', header: 'Name', render: (r) => r.name },
  { key: 'category', header: 'Category', render: (r) => r.category || '—' },
  { key: 'cost', header: 'Cost', render: (r) => formatCents(r.unit_cost_cents, r.currency) },
  { key: 'price', header: 'Price', render: (r) => formatCents(r.unit_price_cents, r.currency) },
  {
    key: 'margin',
    header: 'Margin %',
    render: (r) => (
      <Chip
        label={`${r.gross_margin_pct.toFixed(1)}%`}
        size="small"
        color={r.gross_margin_pct >= 30 ? 'success' : r.gross_margin_pct >= 15 ? 'warning' : 'error'}
      />
    ),
  },
  {
    key: 'status',
    header: 'Status',
    render: (r) => (
      <Chip label={r.status} size="small" color={r.status === 'active' ? 'success' : 'default'} />
    ),
  },
];

export default function ProductsPage() {
  const { data, isLoading } = useProducts();

  if (isLoading) return <CircularProgress />;

  return (
    <Box>
      <Typography variant="h4" gutterBottom>Products</Typography>
      <DataTable columns={columns} data={data?.items ?? []} keyExtractor={(r) => r.id} />
    </Box>
  );
}
