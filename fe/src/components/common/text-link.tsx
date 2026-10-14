import type { ComponentProps } from 'react'
import Link from 'next/link'

import { cn } from '@/lib/utils'

type TextLinkProps = ComponentProps<typeof Link>

export function TextLink({ className, children, ...props }: TextLinkProps) {
  return (
    <Link
      className={cn(
        'text-foreground underline decoration-muted-foreground/40 underline-offset-4 transition-colors duration-300 ease-out hover:decoration-current!',
        className,
      )}
      {...props}
    >
      {children}
    </Link>
  )
}
