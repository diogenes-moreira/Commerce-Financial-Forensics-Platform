import { useState } from 'react';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import Grid from '@mui/material/Grid2';
import ToggleButton from '@mui/material/ToggleButton';
import ToggleButtonGroup from '@mui/material/ToggleButtonGroup';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';
import IconButton from '@mui/material/IconButton';
import Collapse from '@mui/material/Collapse';
import KeyboardArrowDownIcon from '@mui/icons-material/KeyboardArrowDown';
import KeyboardArrowRightIcon from '@mui/icons-material/KeyboardArrowRight';
import { formatCents } from '../../utils/formatCurrency';
import type { PLReport } from '../../types';

const placeholderData: PLReport[] = [
  {
    period: '2026-01', granularity: 'month', revenue_cents: 1520000, cogs_cents: 890000,
    gross_profit_cents: 630000, discount_cents: 45000, commission_cents: 152000,
    shipping_cents: 38000, platform_fees_cents: 22000, net_profit_cents: 373000,
    gross_margin_pct: 41.4, net_margin_pct: 24.5, order_count: 342, currency: 'USD',
    children: [
      {
        period: '2026-01-W1', granularity: 'week', revenue_cents: 380000, cogs_cents: 222500,
        gross_profit_cents: 157500, discount_cents: 11250, commission_cents: 38000,
        shipping_cents: 9500, platform_fees_cents: 5500, net_profit_cents: 93250,
        gross_margin_pct: 41.4, net_margin_pct: 24.5, order_count: 86, currency: 'USD',
      },
      {
        period: '2026-01-W2', granularity: 'week', revenue_cents: 410000, cogs_cents: 240000,
        gross_profit_cents: 170000, discount_cents: 12500, commission_cents: 41000,
        shipping_cents: 10250, platform_fees_cents: 5950, net_profit_cents: 100300,
        gross_margin_pct: 41.5, net_margin_pct: 24.5, order_count: 92, currency: 'USD',
      },
      {
        period: '2026-01-W3', granularity: 'week', revenue_cents: 365000, cogs_cents: 213750,
        gross_profit_cents: 151250, discount_cents: 10625, commission_cents: 36500,
        shipping_cents: 9125, platform_fees_cents: 5275, net_profit_cents: 89725,
        gross_margin_pct: 41.4, net_margin_pct: 24.6, order_count: 82, currency: 'USD',
      },
      {
        period: '2026-01-W4', granularity: 'week', revenue_cents: 365000, cogs_cents: 213750,
        gross_profit_cents: 151250, discount_cents: 10625, commission_cents: 36500,
        shipping_cents: 9125, platform_fees_cents: 5275, net_profit_cents: 89725,
        gross_margin_pct: 41.4, net_margin_pct: 24.6, order_count: 82, currency: 'USD',
      },
    ],
  },
  {
    period: '2025-12', granularity: 'month', revenue_cents: 1380000, cogs_cents: 820000,
    gross_profit_cents: 560000, discount_cents: 52000, commission_cents: 138000,
    shipping_cents: 34500, platform_fees_cents: 20100, net_profit_cents: 315400,
    gross_margin_pct: 40.6, net_margin_pct: 22.9, order_count: 310, currency: 'USD',
    children: [
      {
        period: '2025-12-W1', granularity: 'week', revenue_cents: 345000, cogs_cents: 205000,
        gross_profit_cents: 140000, discount_cents: 13000, commission_cents: 34500,
        shipping_cents: 8625, platform_fees_cents: 5025, net_profit_cents: 78850,
        gross_margin_pct: 40.6, net_margin_pct: 22.9, order_count: 78, currency: 'USD',
      },
    ],
  },
  {
    period: '2025-11', granularity: 'month', revenue_cents: 1250000, cogs_cents: 780000,
    gross_profit_cents: 470000, discount_cents: 38000, commission_cents: 125000,
    shipping_cents: 31250, platform_fees_cents: 18200, net_profit_cents: 257550,
    gross_margin_pct: 37.6, net_margin_pct: 20.6, order_count: 285, currency: 'USD',
  },
  {
    period: '2025-10', granularity: 'month', revenue_cents: 1100000, cogs_cents: 710000,
    gross_profit_cents: 390000, discount_cents: 33000, commission_cents: 110000,
    shipping_cents: 27500, platform_fees_cents: 16050, net_profit_cents: 203450,
    gross_margin_pct: 35.5, net_margin_pct: 18.5, order_count: 258, currency: 'USD',
  },
];

function marginColor(pct: number): string {
  if (pct >= 20) return '#2e7d32';
  if (pct >= 5) return '#ed6c02';
  return '#d32f2f';
}

function PLRow({ row, depth = 0 }: { row: PLReport; depth?: number }) {
  const [open, setOpen] = useState(false);
  const hasChildren = row.children && row.children.length > 0;

  return (
    <>
      <TableRow hover>
        <TableCell sx={{ pl: 2 + depth * 3 }}>
          <Box sx={{ display: 'flex', alignItems: 'center' }}>
            {hasChildren && (
              <IconButton size="small" onClick={() => setOpen(!open)} sx={{ mr: 0.5 }}>
                {open ? <KeyboardArrowDownIcon /> : <KeyboardArrowRightIcon />}
              </IconButton>
            )}
            {!hasChildren && <Box sx={{ width: 34 }} />}
            <Typography variant="body2" sx={{ fontWeight: depth === 0 ? 600 : 400 }}>
              {row.period}
            </Typography>
          </Box>
        </TableCell>
        <TableCell align="right">{formatCents(row.revenue_cents, row.currency)}</TableCell>
        <TableCell align="right">{formatCents(row.cogs_cents, row.currency)}</TableCell>
        <TableCell align="right">{formatCents(row.gross_profit_cents, row.currency)}</TableCell>
        <TableCell align="right">{formatCents(row.discount_cents, row.currency)}</TableCell>
        <TableCell align="right">{formatCents(row.commission_cents, row.currency)}</TableCell>
        <TableCell align="right">{formatCents(row.net_profit_cents, row.currency)}</TableCell>
        <TableCell align="right">
          <Typography variant="body2" sx={{ color: marginColor(row.gross_margin_pct), fontWeight: 600 }}>
            {row.gross_margin_pct.toFixed(1)}%
          </Typography>
        </TableCell>
        <TableCell align="right">
          <Typography variant="body2" sx={{ color: marginColor(row.net_margin_pct), fontWeight: 600 }}>
            {row.net_margin_pct.toFixed(1)}%
          </Typography>
        </TableCell>
        <TableCell align="right">{row.order_count.toLocaleString()}</TableCell>
      </TableRow>
      {hasChildren && (
        <TableRow>
          <TableCell colSpan={10} sx={{ p: 0, border: 0 }}>
            <Collapse in={open} timeout="auto" unmountOnExit>
              <Table size="small">
                <TableBody>
                  {row.children!.map((child) => (
                    <PLRow key={child.period} row={child} depth={depth + 1} />
                  ))}
                </TableBody>
              </Table>
            </Collapse>
          </TableCell>
        </TableRow>
      )}
    </>
  );
}

export default function PnLPage() {
  const [granularity, setGranularity] = useState('month');

  const data = placeholderData;

  const totals = data.reduce(
    (acc, r) => ({
      revenue: acc.revenue + r.revenue_cents,
      cogs: acc.cogs + r.cogs_cents,
      gross: acc.gross + r.gross_profit_cents,
      net: acc.net + r.net_profit_cents,
    }),
    { revenue: 0, cogs: 0, gross: 0, net: 0 },
  );

  const summaryCards = [
    { title: 'Total Revenue', value: formatCents(totals.revenue), color: '#1976d2' },
    { title: 'Total COGS', value: formatCents(totals.cogs), color: '#d32f2f' },
    { title: 'Gross Profit', value: formatCents(totals.gross), color: '#2e7d32' },
    { title: 'Net Profit', value: formatCents(totals.net), color: '#7b1fa2' },
  ];

  return (
    <Box>
      <Typography variant="h4" gutterBottom>P&L Report</Typography>

      <Grid container spacing={3} sx={{ mb: 3 }}>
        {summaryCards.map((card) => (
          <Grid size={{ xs: 12, sm: 6, md: 3 }} key={card.title}>
            <Card>
              <CardContent>
                <Typography color="text.secondary" gutterBottom>{card.title}</Typography>
                <Typography variant="h4" sx={{ color: card.color }}>{card.value}</Typography>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>

      <Box sx={{ mb: 3 }}>
        <ToggleButtonGroup
          value={granularity}
          exclusive
          onChange={(_, val) => val && setGranularity(val)}
          size="small"
        >
          <ToggleButton value="day">Day</ToggleButton>
          <ToggleButton value="month">Month</ToggleButton>
          <ToggleButton value="quarter">Quarter</ToggleButton>
          <ToggleButton value="year">Year</ToggleButton>
        </ToggleButtonGroup>
      </Box>

      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell sx={{ fontWeight: 600 }}>Period</TableCell>
              <TableCell sx={{ fontWeight: 600 }} align="right">Revenue</TableCell>
              <TableCell sx={{ fontWeight: 600 }} align="right">COGS</TableCell>
              <TableCell sx={{ fontWeight: 600 }} align="right">Gross Profit</TableCell>
              <TableCell sx={{ fontWeight: 600 }} align="right">Discounts</TableCell>
              <TableCell sx={{ fontWeight: 600 }} align="right">Commissions</TableCell>
              <TableCell sx={{ fontWeight: 600 }} align="right">Net Profit</TableCell>
              <TableCell sx={{ fontWeight: 600 }} align="right">Gross Margin %</TableCell>
              <TableCell sx={{ fontWeight: 600 }} align="right">Net Margin %</TableCell>
              <TableCell sx={{ fontWeight: 600 }} align="right">Orders</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {data.map((row) => (
              <PLRow key={row.period} row={row} />
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </Box>
  );
}
