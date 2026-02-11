import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import CircularProgress from '@mui/material/CircularProgress';
import DataTable, { type Column } from '../../components/common/DataTable';
import { useExchangeRates } from '../../hooks/useExchangeRates';
import type { ExchangeRate } from '../../types';

const columns: Column<ExchangeRate>[] = [
  { key: 'base_currency', header: 'Base', render: (r) => r.base_currency },
  { key: 'quote_currency', header: 'Quote', render: (r) => r.quote_currency },
  { key: 'rate', header: 'Rate', render: (r) => r.rate.toFixed(6) },
  { key: 'inverse_rate', header: 'Inverse Rate', render: (r) => r.inverse_rate.toFixed(6) },
  { key: 'source', header: 'Source', render: (r) => r.source },
  { key: 'effective_date', header: 'Effective Date', render: (r) => new Date(r.effective_date).toLocaleDateString() },
];

export default function ExchangeRatesPage() {
  const { data, isLoading } = useExchangeRates();

  if (isLoading) return <CircularProgress />;

  return (
    <Box>
      <Typography variant="h4" gutterBottom>Exchange Rates</Typography>
      <DataTable columns={columns} data={data?.items ?? []} keyExtractor={(r) => r.id} />
    </Box>
  );
}
