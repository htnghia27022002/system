import type { ReactNode } from 'react'

type AccessControlPageHeaderProps = {
  title: string
  description: string
  actions?: ReactNode
}

export function AccessControlPageHeader({
  title,
  description,
  actions,
}: AccessControlPageHeaderProps) {
  return (
    <div className="flex flex-col gap-4 md:flex-row md:items-start md:justify-between">
      <div className="min-w-0 space-y-1">
        <h1 className="text-2xl font-semibold tracking-tight">{title}</h1>
        <p className="text-sm text-muted-foreground">{description}</p>
      </div>
      {actions ? (
        <div className="flex w-full shrink-0 flex-col gap-2 md:w-auto md:flex-row md:items-center">
          {actions}
        </div>
      ) : null}
    </div>
  )
}
