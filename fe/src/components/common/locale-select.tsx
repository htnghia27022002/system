'use client'

import type { FC } from 'react'
import { useTranslation } from 'react-i18next'

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { cn } from '@/lib/utils'

type LocaleCode = 'en' | 'vi'

type LocaleOption = {
  code: LocaleCode
  label: string
  Flag: FC<{ className?: string }>
}

function EnglishFlag({ className }: { className?: string }) {
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

function VietnamFlag({ className }: { className?: string }) {
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

const LOCALES: LocaleOption[] = [
  { code: 'en', label: 'English', Flag: EnglishFlag },
  { code: 'vi', label: 'Tiếng Việt', Flag: VietnamFlag },
]

function normalizeLocale(language: string): LocaleCode {
  const base = language.toLowerCase().split('-')[0]
  return base === 'vi' ? 'vi' : 'en'
}

type LocaleSelectProps = {
  className?: string
  triggerClassName?: string
}

/**
 * Locale Select with flag icon + language label for EN/VI.
 * Replaces binary EN↔VI toggle on public chrome.
 */
export function LocaleSelect({ className, triggerClassName }: LocaleSelectProps) {
  const { i18n, t } = useTranslation('common')
  const current = normalizeLocale(i18n.language)
  const currentLocale = LOCALES.find((locale) => locale.code === current) ?? LOCALES[0]
  const CurrentFlag = currentLocale.Flag

  return (
    <div className={cn('inline-flex', className)}>
      <Select
        value={current}
        onValueChange={(value) => {
          if (value === 'en' || value === 'vi') {
            void i18n.changeLanguage(value)
          }
        }}
      >
        <SelectTrigger
          size="sm"
          aria-label={t('locale.label')}
          className={cn('min-w-[9.5rem] gap-2', triggerClassName)}
        >
          <CurrentFlag />
          <SelectValue>{currentLocale.label}</SelectValue>
        </SelectTrigger>
        <SelectContent align="end" position="popper">
          {LOCALES.map(({ code, label, Flag }) => (
            <SelectItem key={code} value={code} textValue={label}>
              <Flag />
              <span>{label}</span>
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
    </div>
  )
}
