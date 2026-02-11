import Grid from '@mui/material/Grid2';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import Typography from '@mui/material/Typography';
import Box from '@mui/material/Box';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';

const placeholderData = [
  { name: 'Compute', cost: 4200 },
  { name: 'Storage', cost: 1800 },
  { name: 'Network', cost: 900 },
  { name: 'Database', cost: 3100 },
  { name: 'Other', cost: 600 },
];

const summaryCards = [
  { title: 'Total Spend (MTD)', value: '$10,600', color: '#1976d2' },
  { title: 'Active Budgets', value: '3', color: '#2e7d32' },
  { title: 'Open Anomalies', value: '2', color: '#ed6c02' },
  { title: 'Cloud Accounts', value: '4', color: '#9c27b0' },
];

export default function DashboardPage() {
  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Dashboard
      </Typography>

      <Grid container spacing={3} sx={{ mb: 4 }}>
        {summaryCards.map((card) => (
          <Grid size={{ xs: 12, sm: 6, md: 3 }} key={card.title}>
            <Card>
              <CardContent>
                <Typography color="text.secondary" gutterBottom>
                  {card.title}
                </Typography>
                <Typography variant="h4" sx={{ color: card.color }}>
                  {card.value}
                </Typography>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>

      <Card>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            Cost by Category
          </Typography>
          <ResponsiveContainer width="100%" height={300}>
            <BarChart data={placeholderData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="name" />
              <YAxis />
              <Tooltip formatter={(value: number) => `$${value}`} />
              <Bar dataKey="cost" fill="#1976d2" radius={[4, 4, 0, 0]} />
            </BarChart>
          </ResponsiveContainer>
        </CardContent>
      </Card>
    </Box>
  );
}
