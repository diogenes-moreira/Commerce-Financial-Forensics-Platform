import { useQuery } from '@tanstack/react-query';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import Chip from '@mui/material/Chip';
import CircularProgress from '@mui/material/CircularProgress';
import DataTable, { type Column } from '../../components/common/DataTable';
import { cloudAccountsApi } from '../../api/cloudAccounts';
import { formatDateTime } from '../../utils/dateUtils';
import type { CloudAccount } from '../../types';

const columns: Column<CloudAccount>[] = [
  { key: 'name', header: 'Name', render: (r) => r.name },
  { key: 'provider', header: 'Provider', render: (r) => r.provider.toUpperCase() },
  { key: 'status', header: 'Status', render: (r) => <Chip label={r.status} color={r.status === 'active' ? 'success' : 'default'} size="small" /> },
  { key: 'sync', header: 'Sync Status', render: (r) => <Chip label={r.sync_status} size="small" /> },
  { key: 'synced', header: 'Last Synced', render: (r) => r.last_synced_at ? formatDateTime(r.last_synced_at) : 'Never' },
];

export default function CloudAccountsPage() {
  const { data, isLoading } = useQuery({ queryKey: ['cloudAccounts'], queryFn: cloudAccountsApi.list });

  if (isLoading) return <CircularProgress />;

  return (
    <Box>
      <Typography variant="h4" gutterBottom>Cloud Accounts</Typography>
      <DataTable columns={columns} data={data ?? []} keyExtractor={(r) => r.id} />
    </Box>
  );
}
