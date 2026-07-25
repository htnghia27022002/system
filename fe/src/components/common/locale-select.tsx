'use client'

import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'

import {
  EnglishFlag,
  LOCALE_OPTIONS,
  normalizeLocale,
} from '@/components/common/locale-flags'

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
  const [mounted, setMounted] = useState(false)

  useEffect(() => {
    setMounted(true)
  }, [])

  const current = normalizeLocale(i18n.language)
  const currentLocale =
    LOCALE_OPTIONS.find((locale) => locale.code === current) ?? LOCALE_OPTIONS[0]
  const CurrentFlag = currentLocale.Flag

  // Stable SSR/first-paint shell — Radix Select ids + language can differ on hydrate.
  if (!mounted) {
    return (
      <div className={cn('inline-flex', className)}>
        <Button
          variant="outline"
          size="sm"
          className={cn('min-w-[9.5rem] justify-start gap-2', triggerClassName)}
          aria-label={t('locale.label')}
          disabled
        >
          <EnglishFlag />
          <span>English</span>
        </Button>
      </div>
    )
  }

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
          {LOCALE_OPTIONS.map(({ code, label, Flag }) => (
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
