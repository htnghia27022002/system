import { ArrowDownIcon, ArrowUpIcon, ArrowUpDownIcon } from 'lucide-react'
import type { Column } from '@tanstack/react-table'

import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'

type DataTableColumnHeaderProps<TData, TValue> = {
  column: Column<TData, TValue>
  title: string
  className?: string
}

export function DataTableColumnHeader<TData, TValue>({
  column,
  title,
  className,
}: DataTableColumnHeaderProps<TData, TValue>) {
  if (!column.getCanSort()) {
    return (
      <span
        className={cn(
          'text-sm font-medium text-foreground',
          className,
        )}
      >
        {title}
      </span>
    )
  }

  const sorted = column.getIsSorted()

  return (
    <Button
      variant="ghost"
      size="sm"
      className={cn(
        '-ml-2 h-8 gap-1.5 px-2 font-medium text-foreground hover:bg-muted/60',
        sorted && 'text-foreground',
        className,
      )}
      onClick={() => column.toggleSorting(sorted === 'asc')}
      aria-sort={
        sorted === 'asc'
          ? 'ascending'
          : sorted === 'desc'
            ? 'descending'
            : 'none'
      }
    >
      {title}
      {sorted === 'asc' ? (
        <ArrowUpIcon className="size-3.5" aria-hidden />
      ) : sorted === 'desc' ? (
        <ArrowDownIcon className="size-3.5" aria-hidden />
      ) : (
        <ArrowUpDownIcon className="size-3.5 opacity-40" aria-hidden />
      )}
    </Button>
  )
}
