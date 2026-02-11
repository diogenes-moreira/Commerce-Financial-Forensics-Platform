import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import CircularProgress from '@mui/material/CircularProgress';
import Chip from '@mui/material/Chip';
import DataTable, { type Column } from '../../components/common/DataTable';
import { useCustomers } from '../../hooks/useCustomers';
import type { Customer } from '../../types';

const columns: Column<Customer>[] = [
  { key: 'name', header: 'Name', render: (r) => r.name },
  { key: 'email', header: 'Email', render: (r) => r.email || '—' },
  { key: 'segment', header: 'Segment', render: (r) => (
    <Chip label={r.segment || 'unknown'} size="small" variant="outlined" />
  )},
  { key: 'external', header: 'External ID', render: (r) => r.external_id || '—' },
];

export default function CustomersPage() {
  const { data, isLoading } = useCustomers();

  if (isLoading) return <CircularProgress />;

  return (
    <Box>
      <Typography variant="h4" gutterBottom>Customers</Typography>
      <DataTable columns={columns} data={data?.items ?? []} keyExtractor={(r) => r.id} />
    </Box>
  );
}
