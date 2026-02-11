import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import CircularProgress from '@mui/material/CircularProgress';
import LinearProgress from '@mui/material/LinearProgress';
import DataTable, { type Column } from '../../components/common/DataTable';
import { useBudgets } from '../../hooks/useBudgets';
import { formatCents } from '../../utils/formatCurrency';
import type { Budget } from '../../types';

const columns: Column<Budget>[] = [
  { key: 'name', header: 'Budget', render: (r) => r.name },
  { key: 'amount', header: 'Limit', render: (r) => formatCents(r.amount_cents, r.currency) },
  { key: 'spent', header: 'Spent', render: (r) => formatCents(r.spent_cents, r.currency) },
  {
    key: 'usage',
    header: 'Usage',
    render: (r) => (
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, minWidth: 120 }}>
        <LinearProgress
          variant="determinate"
          value={Math.min(r.usage_pct, 100)}
          color={r.usage_pct >= 100 ? 'error' : r.usage_pct >= 80 ? 'warning' : 'primary'}
          sx={{ flexGrow: 1 }}
        />
        <Typography variant="body2">{r.usage_pct}%</Typography>
      </Box>
    ),
  },
  { key: 'period', header: 'Period', render: (r) => `${r.period_start} - ${r.period_end}` },
];

export default function BudgetsPage() {
  const { data, isLoading } = useBudgets();

  if (isLoading) return <CircularProgress />;

  return (
    <Box>
      <Typography variant="h4" gutterBottom>Budgets</Typography>
      <DataTable columns={columns} data={data ?? []} keyExtractor={(r) => r.id} />
    </Box>
  );
}
