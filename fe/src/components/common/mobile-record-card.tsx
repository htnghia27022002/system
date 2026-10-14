import type { ReactNode } from 'react'

import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'

export type MobileRecordField = {
  label: string
  value: ReactNode
}

type MobileRecordCardProps = {
  title: ReactNode
  actions?: ReactNode
  fields: MobileRecordField[]
}

export function MobileRecordCard({
  title,
  actions,
  fields,
}: MobileRecordCardProps) {
  return (
    <Card className="overflow-hidden py-0">
      <CardHeader className="flex flex-row items-start justify-between gap-3 space-y-0 border-b px-4 py-3">
        <div className="min-w-0 flex-1 font-semibold leading-snug">{title}</div>
        {actions ? <div className="shrink-0">{actions}</div> : null}
      </CardHeader>
      <CardContent className="px-4 py-3">
        <dl className="flex flex-col gap-2">
          {fields.map((field, index) => (
            <div key={field.label}>
              {index > 0 ? <Separator className="mb-2" /> : null}
              <div className="flex items-start justify-between gap-4 text-sm">
                <dt className="shrink-0 text-muted-foreground">{field.label}</dt>
                <dd className="min-w-0 text-right font-medium">{field.value}</dd>
              </div>
            </div>
          ))}
        </dl>
      </CardContent>
    </Card>
  )
}
