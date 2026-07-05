import { InboxIcon, SearchXIcon, type LucideIcon } from 'lucide-react'
import type { ReactNode } from 'react'

import { cn } from '@/lib/utils'

type DataTableEmptyProps = {
  message?: string
  title?: string
  description?: string
  icon?: LucideIcon
  action?: ReactNode
  className?: string
  children?: ReactNode
}

export function DataTableEmpty({
  message,
  title,
  description,
  icon: Icon = InboxIcon,
  action,
  className,
  children,
}: DataTableEmptyProps) {
  const resolvedTitle = title ?? message ?? 'No results'
  const resolvedDescription =
    description ??
    (title && message ? message : undefined)

  return (
    <div
      className={cn(
        'flex flex-col items-center justify-center gap-4 px-6 py-14 text-center',
        className,
      )}
    >
      <div className="flex size-14 items-center justify-center rounded-2xl border border-dashed bg-muted/30">
        <Icon
          className="size-7 text-muted-foreground/70"
          strokeWidth={1.5}
          aria-hidden
        />
      </div>
      <div className="flex max-w-sm flex-col gap-1.5">
        <p className="text-base font-medium text-foreground">{resolvedTitle}</p>
        {resolvedDescription ? (
          <p className="text-sm leading-relaxed text-muted-foreground">
            {resolvedDescription}
          </p>
        ) : null}
      </div>
      {action ? <div className="pt-1">{action}</div> : null}
      {children}
    </div>
  )
}

export { SearchXIcon as DataTableEmptySearchIcon }
