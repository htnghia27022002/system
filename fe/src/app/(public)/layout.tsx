import type { ReactNode } from 'react'

/**
 * Thin public shell — `/`, `/tools`, and nested tool pages use normal document scroll.
 */
export default function PublicGroupLayout({ children }: { children: ReactNode }) {
  return children
}
