import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import CircularProgress from '@mui/material/CircularProgress';
import Chip from '@mui/material/Chip';
import DataTable, { type Column } from '../../components/common/DataTable';
import { useSellers } from '../../hooks/useSellers';
import type { Seller } from '../../types';

const columns: Column<Seller>[] = [
  { key: 'code', header: 'Code', render: (r) => r.code || '—' },
  { key: 'name', header: 'Name', render: (r) => r.name },
  { key: 'email', header: 'Email', render: (r) => r.email || '—' },
  { key: 'commission', header: 'Commission %', render: (r) => `${r.commission_pct}%` },
  {
    key: 'status',
    header: 'Status',
    render: (r) => (
      <Chip label={r.status} size="small" color={r.status === 'active' ? 'success' : 'default'} />
    ),
  },
];

export default function SellersPage() {
  const { data, isLoading } = useSellers();

  if (isLoading) return <CircularProgress />;

  return (
    <Box>
      <Typography variant="h4" gutterBottom>Sellers</Typography>
      <DataTable columns={columns} data={data?.items ?? []} keyExtractor={(r) => r.id} />
    </Box>
  );
}
