import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import CircularProgress from '@mui/material/CircularProgress';
import Chip from '@mui/material/Chip';
import LinearProgress from '@mui/material/LinearProgress';
import DataTable, { type Column } from '../../components/common/DataTable';
import { useImports } from '../../hooks/useImports';
import type { ImportJob } from '../../types';

const statusColorMap: Record<string, 'success' | 'warning' | 'error' | 'info' | 'default'> = {
  completed: 'success',
  pending: 'warning',
  failed: 'error',
  processing: 'info',
};

const columns: Column<ImportJob>[] = [
  { key: 'name', header: 'Name', render: (r) => r.name },
  {
    key: 'source',
    header: 'Source',
    render: (r) => <Chip label={r.source} size="small" variant="outlined" />,
  },
  {
    key: 'entity_type',
    header: 'Entity Type',
    render: (r) => <Chip label={r.entity_type} size="small" />,
  },
  {
    key: 'progress',
    header: 'Progress',
    render: (r) => {
      const pct = r.total_rows > 0 ? (r.processed_rows / r.total_rows) * 100 : 0;
      return (
        <Box sx={{ minWidth: 120 }}>
          <LinearProgress variant="determinate" value={pct} />
          <Typography variant="caption">
            {r.processed_rows} / {r.total_rows}
          </Typography>
        </Box>
      );
    },
  },
  { key: 'failed_rows', header: 'Errors', render: (r) => r.failed_rows },
  {
    key: 'status',
    header: 'Status',
    render: (r) => (
      <Chip label={r.status} size="small" color={statusColorMap[r.status] ?? 'default'} />
    ),
  },
  {
    key: 'created_at',
    header: 'Created',
    render: (r) => new Date(r.created_at).toLocaleString(),
  },
];

export default function ImportsPage() {
  const { data, isLoading } = useImports();

  if (isLoading) return <CircularProgress />;

  return (
    <Box>
      <Typography variant="h4" gutterBottom>Data Imports</Typography>
      <DataTable columns={columns} data={data?.items ?? []} keyExtractor={(r) => r.id} />
    </Box>
  );
}
