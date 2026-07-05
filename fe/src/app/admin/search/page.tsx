import type { Metadata } from 'next'
import { Suspense } from 'react'

import { DataTableSkeleton } from '@/components/common/data-table/data-table-skeleton'
import { AdminSearchPage, SearchAccessGuard } from '@/features/admin-search'

export const metadata: Metadata = {
  title: 'Search',
  robots: { index: false, follow: false },
}

export default function Page() {
  return (
    <SearchAccessGuard>
      <Suspense
        fallback={
          <div className="flex flex-col gap-4 p-4 md:p-6">
            <DataTableSkeleton columns={1} />
          </div>
        }
      >
        <AdminSearchPage />
      </Suspense>
    </SearchAccessGuard>
  )
}
