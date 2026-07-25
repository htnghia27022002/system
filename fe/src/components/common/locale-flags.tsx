'use client'

import type { FC } from 'react'

import { cn } from '@/lib/utils'

export type LocaleCode = 'en' | 'vi'

export type LocaleOption = {
  code: LocaleCode
  label: string
  Flag: FC<{ className?: string }>
}

export function EnglishFlag({ className }: { className?: string }) {
  return (
    <svg
      viewBox="0 0 24 16"
      className={cn('h-3.5 w-5 shrink-0 overflow-hidden rounded-[2px]', className)}
      aria-hidden
    >
      <rect width="24" height="16" fill="#012169" />
      <path d="M0 0 L24 16 M24 0 L0 16" stroke="#fff" strokeWidth="2.5" />
      <path d="M0 0 L24 16 M24 0 L0 16" stroke="#C8102E" strokeWidth="1.2" />
      <path d="M12 0 V16 M0 8 H24" stroke="#fff" strokeWidth="4" />
      <path d="M12 0 V16 M0 8 H24" stroke="#C8102E" strokeWidth="2" />
    </svg>
  )
}

export function VietnamFlag({ className }: { className?: string }) {
  return (
    <svg
      viewBox="0 0 24 16"
      className={cn('h-3.5 w-5 shrink-0 overflow-hidden rounded-[2px]', className)}
      aria-hidden
    >
      <rect width="24" height="16" fill="#DA251D" />
      <polygon
        points="12,3.2 13.4,7.4 17.8,7.4 14.2,9.9 15.6,14.1 12,11.5 8.4,14.1 9.8,9.9 6.2,7.4 10.6,7.4"
        fill="#FFCD00"
      />
    </svg>
  )
}

export const LOCALE_OPTIONS: LocaleOption[] = [
  { code: 'en', label: 'English', Flag: EnglishFlag },
  { code: 'vi', label: 'Tiếng Việt', Flag: VietnamFlag },
]

export function normalizeLocale(language: string): LocaleCode {
  const base = language.toLowerCase().split('-')[0]
  return base === 'vi' ? 'vi' : 'en'
}
