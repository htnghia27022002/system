'use client'

import type { ReactNode } from 'react'

import { cn } from '@/lib/utils'

import { useScrolly } from '../context/scrolly-context'

type ScrollyShellProps = {
  children: ReactNode
  /** proximity after GSAP takes over pinned sections (Task 2+) */
  snapMode?: 'mandatory' | 'proximity'
}

export function ScrollyShell({ children, snapMode = 'proximity' }: ScrollyShellProps) {
  const { scrollerRef } = useScrolly()

  return (
    <div
      ref={scrollerRef}
      id="scrolly-container"
      data-snap-mode={snapMode}
      className={cn(
        'dark h-[100dvh] w-full overflow-x-hidden overflow-y-auto overscroll-y-contain',
        snapMode === 'mandatory' ? 'snap-y snap-mandatory' : 'snap-y snap-proximity',
      )}
      style={{ scrollBehavior: 'auto' }}
    >
      {children}
    </div>
  )
}
