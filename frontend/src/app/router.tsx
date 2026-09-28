import { createBrowserRouter, Navigate } from 'react-router-dom';

import { AppShell } from './AppShell';
import { DashboardPage } from '../pages/DashboardPage';
import { PlaceholderPage } from '../pages/PlaceholderPage';
import { NotFoundPage } from '../pages/NotFoundPage';

// Route placeholder sesuai DESIGN §3. Logika akses/route guard menyusul
// setelah auth (TASK-010/011); semua route dulu ditampilkan.
export const router = createBrowserRouter([
  {
    path: '/login',
    element: (
      <PlaceholderPage
        title="Login"
        note="Form login diimplementasikan pada TASK-010 (setelah API auth tersedia)."
      />
    ),
  },
  {
    path: '/change-password',
    element: (
      <PlaceholderPage
        title="Ganti password"
        note="Form ganti password diimplementasikan pada TASK-011."
      />
    ),
  },
  {
    path: '/',
    element: <AppShell />,
    children: [
      { index: true, element: <Navigate to="/dashboard" replace /> },
      { path: 'dashboard', element: <DashboardPage /> },
      { path: 'tickets', element: <PlaceholderPage title="Daftar tiket" note="TASK-016/017." /> },
      { path: 'tickets/new', element: <PlaceholderPage title="Buat tiket" note="TASK-018." /> },
      { path: 'tickets/:id', element: <PlaceholderPage title="Detail tiket" note="TASK-020." /> },
      { path: 'notifications', element: <PlaceholderPage title="Notifikasi" note="TASK-026." /> },
      { path: 'admin/users', element: <PlaceholderPage title="Akun pengguna" note="TASK-025." /> },
      { path: 'admin/categories', element: <PlaceholderPage title="Kategori" note="TASK-025." /> },
      { path: 'admin/departments', element: <PlaceholderPage title="Departemen" note="TASK-025." /> },
      { path: 'admin/audit', element: <PlaceholderPage title="Audit" note="TASK-032." /> },
      { path: '*', element: <NotFoundPage /> },
    ],
  },
]);
