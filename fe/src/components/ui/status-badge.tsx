import { cn } from '@/lib/utils'

export type StatusVariant = 'success' | 'warning' | 'error' | 'info' | 'neutral'

/**
 * Dot color for each semantic variant.
 * Uses semantic theme tokens from src/styles/index.css.
 */
const DOT_CLASSES: Record<StatusVariant, string> = {
  success: 'bg-success',
  warning: 'bg-warning',
  error:   'bg-destructive',
  info:    'bg-primary',
  neutral: 'bg-muted-foreground/50',
}

const PING_CLASSES: Record<StatusVariant, string> = {
  success: 'bg-success/60',
  warning: 'bg-warning/60',
  error:   'bg-destructive/60',
  info:    'bg-primary/60',
  neutral: 'bg-muted-foreground/30',
}

const BADGE_CLASSES: Record<StatusVariant, string> = {
  success: 'bg-success/10 text-success dark:bg-success/20',
  warning: 'bg-warning/10 text-warning dark:bg-warning/20',
  error:   'bg-destructive/10 text-destructive dark:bg-destructive/20',
  info:    'bg-primary/10 text-primary dark:bg-primary/20',
  neutral: 'bg-muted text-muted-foreground',
}

type StatusBadgeProps = {
  variant: StatusVariant
  label: string
  /**
   * When true, the dot shows an `animate-ping` outer ring.
   * Use for "live" / "active" states only.
   * Automatically hidden when the user has "prefer-reduced-motion" set.
   */
  pulse?: boolean
  className?: string
}

export function StatusBadge({
  variant,
  label,
  pulse = false,
  className,
}: StatusBadgeProps) {
  return (
    <span
      role="status"
      className={cn(
        'inline-flex items-center gap-1.5 rounded-full px-2.5 py-0.5 text-xs font-medium',
        BADGE_CLASSES[variant],
        className,
      )}
    >
      {/* Dot with optional pulse ring */}
      <span className="relative flex h-1.5 w-1.5 shrink-0" aria-hidden="true">
        {pulse && (
          <span
            className={cn(
              'absolute inline-flex h-full w-full rounded-full opacity-75',
              'animate-ping motion-reduce:hidden',
              PING_CLASSES[variant],
            )}
          />
        )}
        <span
          className={cn(
            'relative inline-flex h-1.5 w-1.5 rounded-full',
            DOT_CLASSES[variant],
          )}
        />
      </span>

      {label}
    </span>
  )
}
