'use client'

import dynamic from 'next/dynamic'
import { useEffect, useState } from 'react'

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
 * Hero 3D stage. Waits for client mount before loading the canvas so
 * `dynamic(..., { ssr: false })` does not hydrate-mismatch the loading shell.
 */
export function HeroCharacterStage({ className }: HeroCharacterStageProps) {
  const reducedMotion = useReducedMotion()
  const [mounted, setMounted] = useState(false)

  useEffect(() => {
    setMounted(true)
  }, [])

  const showCanvas = mounted && !reducedMotion

  return (
    <div
      className={cn(
        'relative mx-auto w-full overflow-hidden rounded-2xl border border-border/70 bg-background/30 shadow-sm',
        'aspect-[4/5] max-w-sm min-h-[18rem]',
        'sm:max-w-md sm:min-h-[20rem]',
        'md:max-w-lg md:min-h-[22rem]',
        'lg:aspect-auto lg:max-w-none lg:h-[min(28rem,58vh)] lg:min-h-[24rem]',
        'xl:h-[min(32rem,60vh)]',
        className,
      )}
    >
      {showCanvas ? (
        <HeroCharacterCanvas className="absolute inset-0 h-full w-full" />
      ) : (
        <HeroCharacterPlaceholder className="absolute inset-0" />
      )}
    </div>
  )
}
