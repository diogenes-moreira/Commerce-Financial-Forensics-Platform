import Drawer from '@mui/material/Drawer';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import Toolbar from '@mui/material/Toolbar';
import Typography from '@mui/material/Typography';
import Box from '@mui/material/Box';
import Divider from '@mui/material/Divider';
import DashboardIcon from '@mui/icons-material/Dashboard';
import CloudIcon from '@mui/icons-material/Cloud';
import AttachMoneyIcon from '@mui/icons-material/AttachMoney';
import AccountBalanceIcon from '@mui/icons-material/AccountBalance';
import WarningIcon from '@mui/icons-material/Warning';
import AssessmentIcon from '@mui/icons-material/Assessment';
import SettingsIcon from '@mui/icons-material/Settings';
import Inventory2Icon from '@mui/icons-material/Inventory2';
import StorefrontIcon from '@mui/icons-material/Storefront';
import PeopleIcon from '@mui/icons-material/People';
import ShoppingCartIcon from '@mui/icons-material/ShoppingCart';
import MenuBookIcon from '@mui/icons-material/MenuBook';
import PaymentIcon from '@mui/icons-material/Payment';
import CurrencyExchangeIcon from '@mui/icons-material/CurrencyExchange';
import LocalOfferIcon from '@mui/icons-material/LocalOffer';
import BarChartIcon from '@mui/icons-material/BarChart';
import ShowChartIcon from '@mui/icons-material/ShowChart';
import CompareArrowsIcon from '@mui/icons-material/CompareArrows';
import TimelineIcon from '@mui/icons-material/Timeline';
import GroupWorkIcon from '@mui/icons-material/GroupWork';
import FileUploadIcon from '@mui/icons-material/FileUpload';
import IntegrationInstructionsIcon from '@mui/icons-material/IntegrationInstructions';
import { useNavigate, useLocation } from 'react-router-dom';

const DRAWER_WIDTH = 260;

type NavEntry =
  | { kind: 'link'; label: string; path: string; icon: React.ReactNode }
  | { kind: 'divider'; label: string };

const navItems: NavEntry[] = [
  { kind: 'link', label: 'Dashboard', path: '/', icon: <DashboardIcon /> },
  { kind: 'link', label: 'Cloud Accounts', path: '/cloud-accounts', icon: <CloudIcon /> },
  { kind: 'link', label: 'Cost Explorer', path: '/costs', icon: <AttachMoneyIcon /> },
  { kind: 'link', label: 'Budgets', path: '/budgets', icon: <AccountBalanceIcon /> },
  { kind: 'link', label: 'Anomalies', path: '/anomalies', icon: <WarningIcon /> },
  { kind: 'link', label: 'Reports', path: '/reports', icon: <AssessmentIcon /> },
  { kind: 'divider', label: 'Commerce' },
  { kind: 'link', label: 'Products', path: '/products', icon: <Inventory2Icon /> },
  { kind: 'link', label: 'Sellers', path: '/sellers', icon: <StorefrontIcon /> },
  { kind: 'link', label: 'Customers', path: '/customers', icon: <PeopleIcon /> },
  { kind: 'link', label: 'Orders', path: '/orders', icon: <ShoppingCartIcon /> },
  { kind: 'divider', label: 'Financial' },
  { kind: 'link', label: 'Ledger', path: '/ledger', icon: <MenuBookIcon /> },
  { kind: 'link', label: 'Payments', path: '/payments', icon: <PaymentIcon /> },
  { kind: 'link', label: 'Exchange Rates', path: '/exchange-rates', icon: <CurrencyExchangeIcon /> },
  { kind: 'link', label: 'Discounts', path: '/discounts', icon: <LocalOfferIcon /> },
  { kind: 'divider', label: 'Analytics' },
  { kind: 'link', label: 'P&L Report', path: '/pnl', icon: <BarChartIcon /> },
  { kind: 'link', label: 'Margins', path: '/margins', icon: <ShowChartIcon /> },
  { kind: 'link', label: 'Drift', path: '/drift', icon: <CompareArrowsIcon /> },
  { kind: 'link', label: 'Events', path: '/events', icon: <TimelineIcon /> },
  { kind: 'link', label: 'Cohorts', path: '/cohorts', icon: <GroupWorkIcon /> },
  { kind: 'divider', label: 'Data' },
  { kind: 'link', label: 'Data Import', path: '/imports', icon: <FileUploadIcon /> },
  { kind: 'link', label: 'Integrations', path: '/integrations', icon: <IntegrationInstructionsIcon /> },
  { kind: 'divider', label: 'Settings' },
  { kind: 'link', label: 'Settings', path: '/settings', icon: <SettingsIcon /> },
];

export default function Sidebar() {
  const navigate = useNavigate();
  const location = useLocation();

  return (
    <Drawer
      variant="permanent"
      sx={{
        width: DRAWER_WIDTH,
        flexShrink: 0,
        '& .MuiDrawer-paper': { width: DRAWER_WIDTH, boxSizing: 'border-box' },
      }}
    >
      <Toolbar>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          <AttachMoneyIcon color="primary" />
          <Typography variant="h6" noWrap color="primary">
            costForensics
          </Typography>
        </Box>
      </Toolbar>
      <List>
        {navItems.map((item) =>
          item.kind === 'divider' ? (
            <Box key={item.label}>
              <Divider sx={{ mt: 1, mb: 0.5 }} />
              <Typography variant="overline" sx={{ px: 2, color: 'text.secondary' }}>
                {item.label}
              </Typography>
            </Box>
          ) : (
            <ListItem key={item.path} disablePadding>
              <ListItemButton
                selected={location.pathname === item.path}
                onClick={() => navigate(item.path)}
              >
                <ListItemIcon>{item.icon}</ListItemIcon>
                <ListItemText primary={item.label} />
              </ListItemButton>
            </ListItem>
          )
        )}
      </List>
    </Drawer>
  );
}
