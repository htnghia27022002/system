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
      <h1 className="font-[family-name:var(--ds-font-display)] text-[32px] font-semibold leading-[1.1] tracking-normal text-foreground sm:text-[40px]">
        {title}
      </h1>
      <p className="text-[17px] leading-[1.47] tracking-[-0.374px] text-muted-foreground">
        {description}
      </p>
    </header>
  )
}
