'use client'

import { create } from 'zustand'

type GlobalLoadingState = {
  /** In-flight HTTP requests tracked for the top nav bar (typically GET). */
  httpPendingCount: number
  incrementHttp: () => void
  decrementHttp: () => void
}

export const useGlobalLoadingStore = create<GlobalLoadingState>((set) => ({
  httpPendingCount: 0,
  incrementHttp: () =>
    set((state) => ({ httpPendingCount: state.httpPendingCount + 1 })),
  decrementHttp: () =>
    set((state) => ({
      httpPendingCount: Math.max(0, state.httpPendingCount - 1),
    })),
}))
