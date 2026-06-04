'use client'

import { Progress } from '@/components/ui/progress'
import { useNavLoading } from '@/hooks/use-nav-loading'
import { cn } from '@/lib/utils'

type NavLoadingBarProps = {
  className?: string
}

/**
 * App-wide top nav loading indicator.
 * Composes shadcn Progress (indeterminate); driven by useNavLoading().
 */
export function NavLoadingBar({ className }: NavLoadingBarProps) {
  const active = useNavLoading()

  return (
    <Progress
      indeterminate={active}
      className={cn(
        'pointer-events-none absolute inset-x-0 bottom-0 h-0.5 rounded-none bg-transparent transition-opacity duration-200',
        !active && 'opacity-0',
        className,
      )}
      aria-hidden={!active}
      aria-valuetext={active ? 'Loading' : undefined}
    />
  )
}
