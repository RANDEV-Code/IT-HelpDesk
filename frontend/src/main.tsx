import React from 'react';
import ReactDOM from 'react-dom/client';
import { RouterProvider } from 'react-router-dom';

import { router } from './app/router';
import { Providers } from './app/providers';
import './index.css';

// Providers() membungkus QueryClientProvider (TanStack Query); RouterProvider
// memakai route placeholder DESIGN §3. Satu origin SPA<->API lewat proxy Vite.
ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <Providers>
      <RouterProvider router={router} />
    </Providers>
  </React.StrictMode>,
);
