import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import CircularProgress from '@mui/material/CircularProgress';
import Chip from '@mui/material/Chip';
import DataTable, { type Column } from '../../components/common/DataTable';
import { useIntegrations } from '../../hooks/useIntegrations';
import type { Integration } from '../../types';

const statusColorMap: Record<string, 'success' | 'warning' | 'error' | 'default'> = {
  active: 'success',
  inactive: 'warning',
  error: 'error',
};

const syncColorMap: Record<string, 'success' | 'info' | 'error' | 'default'> = {
  idle: 'default',
  running: 'info',
  failed: 'error',
};

const columns: Column<Integration>[] = [
  { key: 'name', header: 'Name', render: (r) => r.name },
  {
    key: 'platform',
    header: 'Platform',
    render: (r) => <Chip label={r.platform} size="small" variant="outlined" />,
  },
  {
    key: 'status',
    header: 'Status',
    render: (r) => <Chip label={r.status} size="small" color={statusColorMap[r.status] ?? 'default'} />,
  },
  {
    key: 'sync_status',
    header: 'Sync Status',
    render: (r) => <Chip label={r.sync_status} size="small" color={syncColorMap[r.sync_status] ?? 'default'} />,
  },
  {
    key: 'last_synced_at',
    header: 'Last Synced',
    render: (r) => r.last_synced_at ? new Date(r.last_synced_at).toLocaleString() : 'Never',
  },
  {
    key: 'created_at',
    header: 'Created',
    render: (r) => new Date(r.created_at).toLocaleDateString(),
  },
];

export default function IntegrationsPage() {
  const { data, isLoading } = useIntegrations();

  if (isLoading) return <CircularProgress />;

  return (
    <Box>
      <Typography variant="h4" gutterBottom>Integrations</Typography>
      <DataTable columns={columns} data={data?.items ?? []} keyExtractor={(r) => r.id} />
    </Box>
  );
}
