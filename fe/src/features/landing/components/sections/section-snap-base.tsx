import { forwardRef, type ReactNode } from 'react'

import { cn } from '@/lib/utils'

type SectionSnapBaseProps = {
  id: string
  className?: string
  children: ReactNode
}

export const SectionSnapBase = forwardRef<HTMLElement, SectionSnapBaseProps>(
  function SectionSnapBase({ id, className, children }, ref) {
    return (
      <section
        ref={ref}
        id={id}
        className={cn(
          'relative flex min-h-[100dvh] w-full shrink-0 snap-start snap-always items-center justify-center scroll-mt-16 px-4 py-24 sm:px-8',
          className,
        )}
      >
        <div className="relative z-10 mx-auto w-full max-w-5xl">{children}</div>
      </section>
    )
  },
)
