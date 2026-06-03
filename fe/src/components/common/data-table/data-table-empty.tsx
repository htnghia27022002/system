import { InboxIcon } from 'lucide-react'
import type { ReactNode } from 'react'

type DataTableEmptyProps = {
  message?: string
  children?: ReactNode
}

export function DataTableEmpty({
  message = 'No results.',
  children,
}: DataTableEmptyProps) {
  return (
    <div className="flex flex-col items-center justify-center gap-3 py-12 text-center">
      <InboxIcon
        className="size-10 text-muted-foreground/40"
        strokeWidth={1.5}
        aria-hidden
      />
      <p className="text-sm text-muted-foreground">{message}</p>
      {children}
    </div>
  )
}
