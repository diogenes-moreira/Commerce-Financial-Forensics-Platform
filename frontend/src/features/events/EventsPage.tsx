import { useState } from 'react';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import CircularProgress from '@mui/material/CircularProgress';
import Chip from '@mui/material/Chip';
import TextField from '@mui/material/TextField';
import Stack from '@mui/material/Stack';
import Tooltip from '@mui/material/Tooltip';
import DataTable, { type Column } from '../../components/common/DataTable';
import { useEvents } from '../../hooks/useEvents';
import type { ForensicEvent } from '../../types';

const columns: Column<ForensicEvent>[] = [
  {
    key: 'entity_type',
    header: 'Entity Type',
    render: (r) => <Chip label={r.entity_type} size="small" variant="outlined" />,
  },
  { key: 'entity_id', header: 'Entity ID', render: (r) => r.entity_id.slice(0, 8) + '...' },
  {
    key: 'event_type',
    header: 'Event Type',
    render: (r) => (
      <Chip
        label={r.event_type}
        size="small"
        color={
          r.event_type.includes('created') ? 'success' :
          r.event_type.includes('deleted') ? 'error' :
          r.event_type.includes('updated') ? 'info' : 'default'
        }
      />
    ),
  },
  { key: 'actor_id', header: 'Actor', render: (r) => r.actor_id || '---' },
  { key: 'event_time', header: 'Event Time', render: (r) => new Date(r.event_time).toLocaleString() },
  {
    key: 'hash',
    header: 'Hash',
    render: (r) => (
      <Tooltip title={r.hash_integrity}>
        <Typography variant="body2" sx={{ fontFamily: 'monospace', fontSize: '0.75rem' }}>
          {r.hash_integrity.slice(0, 12)}...
        </Typography>
      </Tooltip>
    ),
  },
];

export default function EventsPage() {
  const [entityType, setEntityType] = useState('');
  const [eventType, setEventType] = useState('');

  const { data, isLoading } = useEvents({
    entity_type: entityType || undefined,
    event_type: eventType || undefined,
  });

  if (isLoading) return <CircularProgress />;

  return (
    <Box>
      <Typography variant="h4" gutterBottom>Forensic Events</Typography>

      <Stack direction="row" spacing={2} sx={{ mb: 3 }}>
        <TextField
          label="Entity Type"
          size="small"
          value={entityType}
          onChange={(e) => setEntityType(e.target.value)}
        />
        <TextField
          label="Event Type"
          size="small"
          value={eventType}
          onChange={(e) => setEventType(e.target.value)}
        />
      </Stack>

      <DataTable columns={columns} data={data?.items ?? []} keyExtractor={(r) => r.id} />
    </Box>
  );
}
