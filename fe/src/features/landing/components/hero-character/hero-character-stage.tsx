'use client'

import dynamic from 'next/dynamic'

import { cn } from '@/lib/utils'

import { useReducedMotion } from '../../hooks/use-reduced-motion'

function HeroCharacterPlaceholder({ className }: { className?: string }) {
  return <div className={cn('h-full w-full bg-transparent', className)} aria-hidden />
}

const HeroCharacterCanvas = dynamic(
  () =>
    import('./hero-character-canvas').then((mod) => mod.HeroCharacterCanvas),
  { ssr: false, loading: () => <HeroCharacterPlaceholder /> },
)

type HeroCharacterStageProps = {
  className?: string
}

/**
 * Hero 3D stage — large, fluid size across breakpoints (no tiny capped frame).
 */
export function HeroCharacterStage({ className }: HeroCharacterStageProps) {
  const reducedMotion = useReducedMotion()

  return (
    <div
      className={cn(
        'relative mx-auto w-full overflow-hidden rounded-2xl border border-border/70 bg-background/30 shadow-sm',
        // Mobile: wide + tall enough; desktop: fill column with a solid viewport-based height
        'aspect-[4/5] max-w-md min-h-[20rem]',
        'sm:max-w-lg sm:min-h-[24rem]',
        'md:max-w-xl md:min-h-[28rem]',
        'lg:aspect-auto lg:max-w-none lg:h-[min(38rem,72vh)] lg:min-h-[32rem]',
        'xl:h-[min(44rem,76vh)]',
        className,
      )}
    >
      {reducedMotion ? (
        <HeroCharacterPlaceholder className="absolute inset-0" />
      ) : (
        <HeroCharacterCanvas className="absolute inset-0 h-full w-full" />
      )}
    </div>
  )
}
