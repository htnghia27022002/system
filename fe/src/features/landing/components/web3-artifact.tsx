'use client'

import { cn } from '@/lib/utils'

type Web3ArtifactProps = {
  artifactId?: string
  variant?: 'hex' | 'diamond' | 'ring'
  className?: string
  label?: string
}

export function Web3Artifact({
  artifactId,
  variant = 'hex',
  className,
  label,
}: Web3ArtifactProps) {
  return (
    <div
      data-artifact-id={artifactId}
      data-artifact={variant}
      className={cn(
        'web3-artifact flex items-center justify-center rounded-2xl border border-primary/40 bg-card/30 p-6 shadow-[0_0_40px_-8px_var(--landing-glow-primary)] backdrop-blur-md will-change-transform',
        className,
      )}
      aria-hidden={!label}
    >
      {variant === 'hex' && (
        <svg viewBox="0 0 120 120" className="h-24 w-24 text-primary" fill="none">
          <polygon
            points="60,8 108,32 108,88 60,112 12,88 12,32"
            stroke="currentColor"
            strokeWidth="2"
            className="drop-shadow-[0_0_12px_var(--landing-glow-primary)]"
          />
          <polygon
            points="60,28 88,44 88,76 60,92 32,76 32,44"
            stroke="currentColor"
            strokeWidth="1"
            opacity="0.5"
          />
        </svg>
      )}
      {variant === 'diamond' && (
        <svg viewBox="0 0 120 120" className="h-24 w-24 text-[var(--chart-2)]" fill="none">
          <path
            d="M60 12 L108 60 L60 108 L12 60 Z"
            stroke="currentColor"
            strokeWidth="2"
            className="drop-shadow-[0_0_12px_var(--landing-glow-accent)]"
          />
        </svg>
      )}
      {variant === 'ring' && (
        <svg viewBox="0 0 120 120" className="h-24 w-24 text-[var(--chart-5)]" fill="none">
          <circle cx="60" cy="60" r="44" stroke="currentColor" strokeWidth="2" opacity="0.9" />
          <circle cx="60" cy="60" r="28" stroke="currentColor" strokeWidth="1" opacity="0.45" />
          <circle cx="60" cy="60" r="8" fill="currentColor" opacity="0.8" />
        </svg>
      )}
      {label ? <span className="sr-only">{label}</span> : null}
    </div>
  )
}
