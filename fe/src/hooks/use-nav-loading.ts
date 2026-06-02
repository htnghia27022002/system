import { useIsFetching } from '@tanstack/react-query'

import { useGlobalLoadingStore } from '@/store/global-loading-store'

/**
 * Top nav indeterminate bar — page/data loads (React Query + tracked GET axios).
 * Form mutations (login/register) use button-level loading instead.
 */
export function useNavLoading(): boolean {
  const queryFetchingCount = useIsFetching()
  const httpPendingCount = useGlobalLoadingStore((state) => state.httpPendingCount)

  return queryFetchingCount > 0 || httpPendingCount > 0
}
