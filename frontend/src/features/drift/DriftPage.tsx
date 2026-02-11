import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import CircularProgress from '@mui/material/CircularProgress';
import Chip from '@mui/material/Chip';
import DataTable, { type Column } from '../../components/common/DataTable';
import { useDrift } from '../../hooks/useDrift';
import { formatCents } from '../../utils/formatCurrency';
import type { DriftResult } from '../../types';

const severityColorMap: Record<string, 'success' | 'warning' | 'error' | 'info'> = {
  low: 'success',
  medium: 'warning',
  high: 'error',
  critical: 'error',
};

const columns: Column<DriftResult>[] = [
  { key: 'order_id', header: 'Order ID', render: (r) => r.order_id.slice(0, 8) + '...' },
  { key: 'field_name', header: 'Field', render: (r) => r.field_name },
  { key: 'expected_value', header: 'Expected', render: (r) => formatCents(r.expected_value) },
  { key: 'actual_value', header: 'Actual', render: (r) => formatCents(r.actual_value) },
  {
    key: 'drift_pct',
    header: 'Drift %',
    render: (r) => (
      <Typography
        variant="body2"
        sx={{ color: Math.abs(r.drift_pct) > 10 ? 'error.main' : Math.abs(r.drift_pct) > 5 ? 'warning.main' : 'text.primary' }}
      >
        {r.drift_pct > 0 ? '+' : ''}{r.drift_pct.toFixed(2)}%
      </Typography>
    ),
  },
  {
    key: 'severity',
    header: 'Severity',
    render: (r) => (
      <Chip
        label={r.severity}
        size="small"
        color={severityColorMap[r.severity] ?? 'default'}
        variant={r.severity === 'critical' ? 'filled' : 'outlined'}
      />
    ),
  },
  { key: 'detected_at', header: 'Detected At', render: (r) => new Date(r.detected_at).toLocaleString() },
];

export default function DriftPage() {
  const { data, isLoading } = useDrift();

  if (isLoading) return <CircularProgress />;

  return (
    <Box>
      <Typography variant="h4" gutterBottom>Drift Detection</Typography>
      <DataTable columns={columns} data={data?.items ?? []} keyExtractor={(r) => `${r.order_id}-${r.field_name}`} />
    </Box>
  );
}
