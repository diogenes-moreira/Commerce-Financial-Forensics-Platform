import { useState } from 'react';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import CircularProgress from '@mui/material/CircularProgress';
import Chip from '@mui/material/Chip';
import TextField from '@mui/material/TextField';
import Stack from '@mui/material/Stack';
import DataTable, { type Column } from '../../components/common/DataTable';
import { useLedger } from '../../hooks/useLedger';
import { formatCents } from '../../utils/formatCurrency';
import type { LedgerEntry } from '../../types';

const columns: Column<LedgerEntry>[] = [
  { key: 'account_code', header: 'Account Code', render: (r) => r.account_code },
  {
    key: 'side',
    header: 'Side',
    render: (r) => (
      <Chip
        label={r.side}
        size="small"
        color={r.side === 'debit' ? 'info' : 'success'}
      />
    ),
  },
  { key: 'amount', header: 'Amount', render: (r) => formatCents(r.amount_cents, r.currency) },
  { key: 'description', header: 'Description', render: (r) => r.description },
  { key: 'reference_type', header: 'Ref Type', render: (r) => r.reference_type },
  { key: 'effective_date', header: 'Effective Date', render: (r) => new Date(r.effective_date).toLocaleDateString() },
];

export default function LedgerPage() {
  const [accountCode, setAccountCode] = useState('');
  const [dateFrom, setDateFrom] = useState('');
  const [dateTo, setDateTo] = useState('');

  const { data, isLoading } = useLedger({
    account_code: accountCode || undefined,
    date_from: dateFrom || undefined,
    date_to: dateTo || undefined,
  });

  if (isLoading) return <CircularProgress />;

  return (
    <Box>
      <Typography variant="h4" gutterBottom>Ledger</Typography>

      <Stack direction="row" spacing={2} sx={{ mb: 3 }}>
        <TextField
          label="Account Code"
          size="small"
          value={accountCode}
          onChange={(e) => setAccountCode(e.target.value)}
        />
        <TextField
          label="From"
          type="date"
          size="small"
          InputLabelProps={{ shrink: true }}
          value={dateFrom}
          onChange={(e) => setDateFrom(e.target.value)}
        />
        <TextField
          label="To"
          type="date"
          size="small"
          InputLabelProps={{ shrink: true }}
          value={dateTo}
          onChange={(e) => setDateTo(e.target.value)}
        />
      </Stack>

      <DataTable columns={columns} data={data?.items ?? []} keyExtractor={(r) => r.id} />
    </Box>
  );
}
