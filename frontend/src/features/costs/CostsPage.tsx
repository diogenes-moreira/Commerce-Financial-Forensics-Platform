import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import CircularProgress from '@mui/material/CircularProgress';
import DataTable, { type Column } from '../../components/common/DataTable';
import { useCosts } from '../../hooks/useCosts';
import { formatCents } from '../../utils/formatCurrency';
import type { CostRecord } from '../../types';

const columns: Column<CostRecord>[] = [
  { key: 'service', header: 'Service', render: (r) => r.service },
  { key: 'category', header: 'Category', render: (r) => r.category },
  { key: 'amount', header: 'Amount', render: (r) => formatCents(r.amount_cents, r.currency) },
  { key: 'date', header: 'Usage Date', render: (r) => r.usage_date },
];

export default function CostsPage() {
  const { data, isLoading } = useCosts();

  if (isLoading) return <CircularProgress />;

  return (
    <Box>
      <Typography variant="h4" gutterBottom>Cost Explorer</Typography>
      <DataTable columns={columns} data={data?.items ?? []} keyExtractor={(r) => r.id} />
    </Box>
  );
}
