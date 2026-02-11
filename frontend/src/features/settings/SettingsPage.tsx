import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';

export default function SettingsPage() {
  return (
    <Box>
      <Typography variant="h4" gutterBottom>Settings</Typography>
      <Card>
        <CardContent>
          <Typography color="text.secondary">
            Settings page — configure cloud integrations, notification preferences, and team management.
          </Typography>
        </CardContent>
      </Card>
    </Box>
  );
}
