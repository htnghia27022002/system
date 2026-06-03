/**
 * Type-safe enum → badge config mappings for access-control entities.
 *
 * Rule: Every enum value rendered in a table or form must have an entry here.
 * Never use if/switch inside JSX for badge colors — always look up from this map.
 */
import type { StatusVariant } from '@/components/ui/status-badge'
import type { ManagedUserStatus } from './types'

export type StatusBadgeConfig = {
  variant: StatusVariant
  /**
   * Show the animate-ping outer ring.
   * Reserve for "live" / "active" states only.
   */
  pulse: boolean
}

// ---------------------------------------------------------------------------
// User status
// ---------------------------------------------------------------------------
export const USER_STATUS_MAP: Record<ManagedUserStatus, StatusBadgeConfig> = {
  active:   { variant: 'success', pulse: true  },
  inactive: { variant: 'neutral', pulse: false },
}
