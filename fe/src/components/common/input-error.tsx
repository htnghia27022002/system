import type { HTMLAttributes } from 'react'

import { cn } from '@/lib/utils'

type InputErrorProps = HTMLAttributes<HTMLParagraphElement> & {
  message?: string
}

export function InputError({ message, className, ...props }: InputErrorProps) {
  if (!message) {
    return null
  }

  return (
    <p
      {...props}
      className={cn('text-sm text-destructive', className)}
    >
      {message}
    </p>
  )
}
