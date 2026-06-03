type DashboardPageHeaderProps = {
  title: string
  description: string
}

export function DashboardPageHeader({
  title,
  description,
}: DashboardPageHeaderProps) {
  return (
    <header className="mb-6 space-y-1">
      <h1 className="text-2xl font-semibold tracking-tight text-foreground sm:text-3xl">
        {title}
      </h1>
      <p className="text-sm text-muted-foreground">
        {description}
      </p>
    </header>
  )
}
