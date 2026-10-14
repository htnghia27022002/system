import type { ComponentProps, ReactNode } from 'react'
import { Slot } from 'radix-ui'

import { cn } from '@/lib/utils'

/**
 * Standard admin / tool page shell. Every page under `/admin` renders inside
 * one of these so outer padding and vertical rhythm match across features.
 */
export function PageContainer({ className, ...props }: ComponentProps<'div'>) {
  return (
    <div
      data-slot="page-container"
      className={cn('flex min-w-0 flex-1 flex-col gap-4 p-4', className)}
      {...props}
    />
  )
}

type PageHeaderProps = {
  title: ReactNode
  description?: ReactNode
  /** Primary page actions (buttons, filters). Stack under the title on mobile. */
  actions?: ReactNode
  className?: string
}

export function PageHeader({
  title,
  description,
  actions,
  className,
}: PageHeaderProps) {
  return (
    <header
      data-slot="page-header"
      className={cn(
        'flex flex-col gap-4 md:flex-row md:items-start md:justify-between',
        className,
      )}
    >
      <div className="min-w-0 space-y-1">
        <h1 className="text-2xl font-semibold tracking-tight">{title}</h1>
        {description ? (
          <p className="text-sm text-muted-foreground">{description}</p>
        ) : null}
      </div>
      {actions ? (
        <div className="flex w-full shrink-0 flex-col gap-2 md:w-auto md:flex-row md:items-center">
          {actions}
        </div>
      ) : null}
    </header>
  )
}

type PanelProps = ComponentProps<'div'> & {
  /** Render as the child element (e.g. a `<form>` or `<section>`). */
  asChild?: boolean
}

/**
 * Bordered surface for grouping content on a page (forms, toolbars, lists).
 * Same radius and border as `Card`; use `Card` when you need header/footer slots.
 */
export function Panel({ className, asChild = false, ...props }: PanelProps) {
  const Comp = asChild ? Slot.Root : 'div'
  return (
    <Comp
      data-slot="panel"
      className={cn(
        'rounded-xl border border-border bg-card p-4 text-card-foreground',
        className,
      )}
      {...props}
    />
  )
}
