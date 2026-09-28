import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useState } from 'react';
import type { ReactNode } from 'react';

import { mutationRetryPolicy, queryRetryPolicy } from '../lib/api/client';

// ARCHITECTURE §7: TanStack Query untuk server state; cache dikosongkan saat
// logout/perubahan sesi (diwire bersama auth pada milestone berikutnya).
// Kebijakan retry API_SPEC §10: read diulang terbatas, mutasi tidak diulang.
export function Providers({ children }: { children: ReactNode }) {
  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: {
            retry: queryRetryPolicy,
            retryDelay: (attempt) => Math.min(1000 * 2 ** attempt, 8000),
            refetchOnWindowFocus: true,
            staleTime: 10_000,
          },
          mutations: {
            retry: mutationRetryPolicy,
          },
        },
      }),
  );

  return <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>;
}
