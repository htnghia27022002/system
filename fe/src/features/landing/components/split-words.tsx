'use client'

import { cn } from '@/lib/utils'

type SplitWordsProps = {
  text: string
  className?: string
  as?: 'h1' | 'h2' | 'h3' | 'p'
  wordClassName?: string
}

export function SplitWords({
  text,
  className,
  as: Tag = 'h2',
  wordClassName,
}: SplitWordsProps) {
  const words = text.split(/\s+/).filter(Boolean)

  return (
    <Tag className={className}>
      {words.map((word, index) => (
        <span
          key={`${word}-${index}`}
          className={cn('split-word inline-block', wordClassName)}
          data-word-index={index}
        >
          {word}
          {index < words.length - 1 ? '\u00A0' : null}
        </span>
      ))}
    </Tag>
  )
}
