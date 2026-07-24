import { type ReactNode } from 'react'

import { cn } from '@/lib/utils'

type SectionBaseProps = {
  id: string
  className?: string
  /** Extra classes for the inner max-width content shell. */
  contentClassName?: string
  children: ReactNode
  /** Full-bleed decorations behind content (grids, glows). Placed outside the max-width column. */
  backdrop?: ReactNode
}

/** Normal document-flow section with standard padding (no snap / full-viewport lock). */
export function SectionBase({
  id,
  className,
  contentClassName,
  children,
  backdrop,
}: SectionBaseProps) {
  return (
    <section id={id} className={cn('relative scroll-mt-20 px-4 py-14 sm:px-8 sm:py-20', className)}>
      {backdrop}
      <div className={cn('relative z-10 mx-auto w-full max-w-5xl', contentClassName)}>
        {children}
      </div>
    </section>
  )
}
