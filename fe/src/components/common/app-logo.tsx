import Link from 'next/link'

import { cn } from '@/lib/utils'

type AppLogoProps = {
  className?: string
}

export function AppLogo({ className }: AppLogoProps) {
  return (
    <Link
      href="/"
      className={cn(
        'flex size-9 items-center justify-center rounded-md bg-primary text-sm font-semibold text-primary-foreground',
        className,
      )}
    >
      S
    </Link>
  )
}
