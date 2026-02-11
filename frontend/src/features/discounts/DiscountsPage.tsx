import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import CircularProgress from '@mui/material/CircularProgress';
import Chip from '@mui/material/Chip';
import DataTable, { type Column } from '../../components/common/DataTable';
import { usePromotionRules } from '../../hooks/useDiscounts';
import type { PromotionRule } from '../../types';

const columns: Column<PromotionRule>[] = [
  { key: 'name', header: 'Name', render: (r) => r.name },
  {
    key: 'rule_type',
    header: 'Type',
    render: (r) => <Chip label={r.rule_type} size="small" variant="outlined" />,
  },
  {
    key: 'value',
    header: 'Value',
    render: (r) =>
      r.rule_type === 'percentage' ? `${r.value}%` : `$${(r.value / 100).toFixed(2)}`,
  },
  {
    key: 'funding_source',
    header: 'Funding Source',
    render: (r) => (
      <Chip
        label={r.funding_source}
        size="small"
        color={r.funding_source === 'platform' ? 'primary' : r.funding_source === 'seller' ? 'secondary' : 'default'}
      />
    ),
  },
  {
    key: 'usage',
    header: 'Usage',
    render: (r) => `${r.current_usage} / ${r.max_usage_count || '\u221E'}`,
  },
  { key: 'valid_from', header: 'Valid From', render: (r) => new Date(r.valid_from).toLocaleDateString() },
  { key: 'valid_to', header: 'Valid To', render: (r) => new Date(r.valid_to).toLocaleDateString() },
  {
    key: 'status',
    header: 'Status',
    render: (r) => (
      <Chip
        label={r.status}
        size="small"
        color={r.status === 'active' ? 'success' : r.status === 'expired' ? 'error' : 'default'}
      />
    ),
  },
];

export default function DiscountsPage() {
  const { data, isLoading } = usePromotionRules();

  if (isLoading) return <CircularProgress />;

  return (
    <Box>
      <Typography variant="h4" gutterBottom>Discounts & Promotions</Typography>
      <DataTable columns={columns} data={data?.items ?? []} keyExtractor={(r) => r.id} />
    </Box>
  );
}
