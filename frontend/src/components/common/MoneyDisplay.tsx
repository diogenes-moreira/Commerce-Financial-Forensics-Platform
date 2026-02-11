import Typography from '@mui/material/Typography';
import { formatCents } from '../../utils/formatCurrency';

interface MoneyDisplayProps {
  cents: number;
  currency?: string;
  variant?: 'h4' | 'h5' | 'h6' | 'body1' | 'body2';
}

export default function MoneyDisplay({ cents, currency = 'USD', variant = 'body1' }: MoneyDisplayProps) {
  return <Typography variant={variant}>{formatCents(cents, currency)}</Typography>;
}
