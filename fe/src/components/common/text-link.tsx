import type { ComponentProps } from 'react'
import { Link } from 'react-router-dom'

import { cn } from '@/lib/utils'

type TextLinkProps = ComponentProps<typeof Link>

export function TextLink({ className, children, ...props }: TextLinkProps) {
  return (
    <Link
      className={cn(
        'text-foreground underline decoration-neutral-300 underline-offset-4 transition-colors duration-300 ease-out hover:decoration-current! dark:decoration-neutral-500',
        className,
      )}
      {...props}
    >
      {children}
    </Link>
  )
}
