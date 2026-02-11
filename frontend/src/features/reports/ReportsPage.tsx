import { useQuery } from '@tanstack/react-query';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import CircularProgress from '@mui/material/CircularProgress';
import DataTable, { type Column } from '../../components/common/DataTable';
import { reportsApi } from '../../api/reports';
import { formatCents } from '../../utils/formatCurrency';
import { formatDateTime } from '../../utils/dateUtils';
import type { CostReport } from '../../types';

const columns: Column<CostReport>[] = [
  { key: 'name', header: 'Report', render: (r) => r.name },
  { key: 'type', header: 'Type', render: (r) => r.report_type },
  { key: 'period', header: 'Period', render: (r) => `${r.period_start} - ${r.period_end}` },
  { key: 'total', header: 'Total', render: (r) => formatCents(r.total_cents, r.currency) },
  { key: 'generated', header: 'Generated', render: (r) => formatDateTime(r.generated_at) },
];

export default function ReportsPage() {
  const { data, isLoading } = useQuery({ queryKey: ['reports'], queryFn: reportsApi.list });

  if (isLoading) return <CircularProgress />;

  return (
    <Box>
      <Typography variant="h4" gutterBottom>Cost Reports</Typography>
      <DataTable columns={columns} data={data ?? []} keyExtractor={(r) => r.id} />
    </Box>
  );
}
