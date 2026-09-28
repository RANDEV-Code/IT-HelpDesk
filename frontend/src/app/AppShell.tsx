import { useState } from 'react';
import { NavLink, Outlet } from 'react-router-dom';

// DESIGN §5: sidebar 240 px (w-60) pada >=1024 px; drawer di bawahnya;
// topbar 64 px (h-16); konten max-width 1440 px; padding 16/20/24 px.

const navItems: { to: string; label: string }[] = [
  { to: '/dashboard', label: 'Dashboard' },
  { to: '/tickets', label: 'Tiket' },
  { to: '/tickets/new', label: 'Buat tiket' },
  { to: '/notifications', label: 'Notifikasi' },
  { to: '/admin/users', label: 'Akun (admin)' },
  { to: '/admin/categories', label: 'Kategori (admin)' },
  { to: '/admin/departments', label: 'Departemen (admin)' },
  { to: '/admin/audit', label: 'Audit (admin)' },
];

function SidebarContent({ onNavigate }: { onNavigate?: () => void }) {
  return (
    <nav aria-label="Navigasi utama" className="flex flex-col gap-1 p-3">
      <p className="px-2 py-3 text-section font-bold text-ink">RANDesk</p>
      {navItems.map((item) => (
        <NavLink
          key={item.to}
          to={item.to}
          onClick={onNavigate}
          className={({ isActive }) =>
            `min-h-11 rounded-control px-3 py-2 text-body ${
              isActive
                ? 'bg-status-open-bg font-semibold text-status-open'
                : 'text-muted hover:bg-background hover:text-ink'
            }`
          }
        >
          {item.label}
        </NavLink>
      ))}
    </nav>
  );
}

export function AppShell() {
  const [drawerOpen, setDrawerOpen] = useState(false);

  return (
    <div className="flex min-h-screen flex-col bg-background">
      <header className="sticky top-0 z-30 flex h-16 items-center gap-3 border-b border-line bg-surface px-4 md:px-5 lg:px-6">
        <button
          type="button"
          aria-label="Buka navigasi"
          aria-expanded={drawerOpen}
          onClick={() => setDrawerOpen(true)}
          className="flex min-h-11 min-w-11 items-center justify-center rounded-control border border-control text-ink lg:hidden"
        >
          <svg aria-hidden="true" width="20" height="20" viewBox="0 0 20 20">
            <path d="M2 4h16M2 10h16M2 16h16" stroke="currentColor" strokeWidth="2" strokeLinecap="round" />
          </svg>
        </button>
        <p className="text-section font-semibold">IT Helpdesk</p>
        <div className="ml-auto text-label text-muted">Sesi &amp; notifikasi: milestone berikutnya</div>
      </header>

      <div className="flex flex-1">
        {/* Drawer < 1024 px */}
        {drawerOpen && (
          <div
            aria-hidden="true"
            className="fixed inset-0 z-40 bg-black/40 lg:hidden"
            onClick={() => setDrawerOpen(false)}
          />
        )}
        <aside
          className={`fixed inset-y-0 left-0 z-50 w-60 overflow-y-auto border-r border-line bg-surface transition-transform duration-150 lg:static lg:z-auto lg:translate-x-0 ${
            drawerOpen ? 'translate-x-0' : '-translate-x-full'
          }`}
        >
          <SidebarContent onNavigate={() => setDrawerOpen(false)} />
        </aside>

        <main className="flex-1 p-4 md:p-5 lg:p-6">
          <div className="mx-auto w-full max-w-[1440px]">
            <Outlet />
          </div>
        </main>
      </div>
    </div>
  );
}
