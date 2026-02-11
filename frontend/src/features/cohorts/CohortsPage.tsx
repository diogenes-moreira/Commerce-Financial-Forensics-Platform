import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';
import type { RetentionCohort } from '../../types';

const placeholderData: RetentionCohort[] = [
  { cohort_month: '2025-07', total_customers: 120, retention: { '0': 100, '1': 68, '2': 52, '3': 45, '4': 41, '5': 38, '6': 35, '7': 33, '8': 31, '9': 29, '10': 28, '11': 27, '12': 26 } },
  { cohort_month: '2025-08', total_customers: 145, retention: { '0': 100, '1': 72, '2': 55, '3': 48, '4': 43, '5': 40, '6': 37, '7': 35, '8': 33, '9': 31, '10': 30, '11': 29 } },
  { cohort_month: '2025-09', total_customers: 132, retention: { '0': 100, '1': 65, '2': 49, '3': 42, '4': 38, '5': 35, '6': 32, '7': 30, '8': 28, '9': 26, '10': 25 } },
  { cohort_month: '2025-10', total_customers: 158, retention: { '0': 100, '1': 70, '2': 54, '3': 46, '4': 42, '5': 39, '6': 36, '7': 34, '8': 32, '9': 30 } },
  { cohort_month: '2025-11', total_customers: 175, retention: { '0': 100, '1': 74, '2': 58, '3': 50, '4': 45, '5': 41, '6': 38, '7': 36, '8': 34 } },
  { cohort_month: '2025-12', total_customers: 190, retention: { '0': 100, '1': 76, '2': 60, '3': 52, '4': 47, '5': 43, '6': 40, '7': 38 } },
  { cohort_month: '2026-01', total_customers: 210, retention: { '0': 100, '1': 78, '2': 62, '3': 54, '4': 49, '5': 45, '6': 42 } },
];

const maxMonths = 12;

function cellBgColor(value: number): string {
  if (value >= 80) return '#1b5e20';
  if (value >= 60) return '#2e7d32';
  if (value >= 40) return '#43a047';
  if (value >= 30) return '#66bb6a';
  if (value >= 20) return '#a5d6a7';
  if (value >= 10) return '#c8e6c9';
  return '#e8f5e9';
}

function cellTextColor(value: number): string {
  return value >= 40 ? '#ffffff' : '#1b5e20';
}

export default function CohortsPage() {
  const data = placeholderData;
  const monthHeaders = Array.from({ length: maxMonths + 1 }, (_, i) => `M${i}`);

  return (
    <Box>
      <Typography variant="h4" gutterBottom>Retention Cohorts</Typography>
      <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
        Each row represents a cohort of customers by their first purchase month. Cell values show the percentage
        of customers who made a repeat purchase in subsequent months.
      </Typography>

      <TableContainer component={Paper} sx={{ overflowX: 'auto' }}>
        <Table size="small">
          <TableHead>
            <TableRow>
              <TableCell sx={{ fontWeight: 600, minWidth: 100 }}>Cohort</TableCell>
              <TableCell sx={{ fontWeight: 600, minWidth: 80 }} align="right">Customers</TableCell>
              {monthHeaders.map((h) => (
                <TableCell key={h} sx={{ fontWeight: 600, minWidth: 60 }} align="center">{h}</TableCell>
              ))}
            </TableRow>
          </TableHead>
          <TableBody>
            {data.map((cohort) => (
              <TableRow key={cohort.cohort_month} hover>
                <TableCell sx={{ fontWeight: 600 }}>{cohort.cohort_month}</TableCell>
                <TableCell align="right">{cohort.total_customers.toLocaleString()}</TableCell>
                {monthHeaders.map((_, idx) => {
                  const val = cohort.retention[String(idx)];
                  return (
                    <TableCell
                      key={idx}
                      align="center"
                      sx={{
                        backgroundColor: val != null ? cellBgColor(val) : 'transparent',
                        color: val != null ? cellTextColor(val) : 'text.disabled',
                        fontWeight: val != null ? 600 : 400,
                        fontSize: '0.75rem',
                      }}
                    >
                      {val != null ? `${val}%` : '---'}
                    </TableCell>
                  );
                })}
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </Box>
  );
}
