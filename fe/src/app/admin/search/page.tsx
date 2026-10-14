import { notFound } from 'next/navigation'

/**
 * Admin unified search is temporarily hidden. The feature module
 * (`@/features/admin-search`) and `GET /api/admin/search` are kept; to restore,
 * render `<SearchAccessGuard><AdminSearchPage /></SearchAccessGuard>` here again
 * and re-add its `metadata` (see git history).
 */
export default function Page() {
  notFound()
}
