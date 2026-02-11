import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import Chip from '@mui/material/Chip';
import CircularProgress from '@mui/material/CircularProgress';
import DataTable, { type Column } from '../../components/common/DataTable';
import { useAnomalies } from '../../hooks/useAnomalies';
import { formatCents } from '../../utils/formatCurrency';
import { formatDateTime } from '../../utils/dateUtils';
import type { Anomaly } from '../../types';

const severityColors: Record<string, 'error' | 'warning' | 'info' | 'default'> = {
  critical: 'error',
  high: 'error',
  medium: 'warning',
  low: 'info',
};

const columns: Column<Anomaly>[] = [
  { key: 'service', header: 'Service', render: (r) => r.service },
  { key: 'expected', header: 'Expected', render: (r) => formatCents(r.expected_cents, r.currency) },
  { key: 'actual', header: 'Actual', render: (r) => formatCents(r.actual_cents, r.currency) },
  { key: 'deviation', header: 'Deviation', render: (r) => `${r.deviation_pct.toFixed(1)}%` },
  { key: 'severity', header: 'Severity', render: (r) => <Chip label={r.severity} color={severityColors[r.severity] ?? 'default'} size="small" /> },
  { key: 'status', header: 'Status', render: (r) => <Chip label={r.status} variant={r.status === 'open' ? 'filled' : 'outlined'} size="small" /> },
  { key: 'detected', header: 'Detected', render: (r) => formatDateTime(r.detected_at) },
];

export default function AnomaliesPage() {
  const { data, isLoading } = useAnomalies();

  if (isLoading) return <CircularProgress />;

  return (
    <Box>
      <Typography variant="h4" gutterBottom>Anomalies</Typography>
      <DataTable columns={columns} data={data ?? []} keyExtractor={(r) => r.id} />
    </Box>
  );
}
