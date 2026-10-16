import { useId, type ReactNode } from 'react'

type SettingsSectionProps = {
  title: ReactNode
  description?: ReactNode
  children: ReactNode
}

/**
 * Settings row: heading + description in a side rail on `lg+`, content
 * (usually a `Card`) on the right. Stacks on smaller screens.
 */
export function SettingsSection({ title, description, children }: SettingsSectionProps) {
  const headingId = useId()
  return (
    <section
      aria-labelledby={headingId}
      className="grid gap-4 lg:grid-cols-[minmax(0,16rem)_minmax(0,1fr)] lg:gap-10"
    >
      <div className="space-y-1">
        <h2 id={headingId} className="text-base font-semibold">
          {title}
        </h2>
        {description ? (
          <p className="text-sm text-muted-foreground">{description}</p>
        ) : null}
      </div>
      <div className="min-w-0">{children}</div>
    </section>
  )
}
