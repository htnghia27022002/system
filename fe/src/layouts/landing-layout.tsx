'use client'

import type { ReactNode } from 'react'

type LandingLayoutProps = {
  children: ReactNode
}

export function LandingLayout({ children }: LandingLayoutProps) {
  return (
    <div className="h-[100dvh] w-full overflow-hidden bg-background text-foreground">
      <main className="h-full">{children}</main>
    </div>
  )
}
