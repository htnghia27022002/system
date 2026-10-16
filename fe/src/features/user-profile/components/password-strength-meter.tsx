'use client'

import { CheckIcon, CircleIcon } from 'lucide-react'
import { useTranslation } from 'react-i18next'

import { cn } from '@/lib/utils'

import type { PasswordStrength } from '../utils/password-strength'

const LEVELS = ['', 'weak', 'fair', 'good', 'strong'] as const

const BAR_CLASSES: Record<number, string> = {
  1: 'bg-destructive',
  2: 'bg-warning',
  3: 'bg-primary',
  4: 'bg-success',
}

const LABEL_CLASSES: Record<number, string> = {
  1: 'text-destructive',
  2: 'text-warning',
  3: 'text-primary',
  4: 'text-success',
}

type PasswordStrengthMeterProps = {
  strength: PasswordStrength
  id?: string
}

export function PasswordStrengthMeter({ strength, id }: PasswordStrengthMeterProps) {
  const { t } = useTranslation('admin')
  const { score, checks } = strength

  const requirements = [
    ['length', checks.length],
    ['mixedCase', checks.mixedCase],
    ['number', checks.number],
    ['symbol', checks.symbol],
  ] as const

  return (
    <div id={id} className="space-y-3">
      <div className="flex items-center gap-3">
        <div className="grid flex-1 grid-cols-4 gap-1" aria-hidden>
          {[1, 2, 3, 4].map((step) => (
            <span
              key={step}
              className={cn(
                'h-1.5 rounded-full bg-muted transition-colors motion-reduce:transition-none',
                score >= step && BAR_CLASSES[score],
              )}
            />
          ))}
        </div>
        <span
          className={cn(
            'min-w-14 text-right text-xs font-medium',
            score > 0 ? LABEL_CLASSES[score] : 'text-muted-foreground',
          )}
          aria-live="polite"
        >
          {score > 0
            ? t(`profile.password.strength.${LEVELS[score]}`)
            : t('profile.password.strength.label')}
        </span>
      </div>
      <ul className="grid gap-1.5 sm:grid-cols-2">
        {requirements.map(([key, met]) => (
          <li
            key={key}
            className={cn(
              'flex items-center gap-2 text-xs',
              met ? 'text-foreground' : 'text-muted-foreground',
            )}
          >
            {met ? (
              <CheckIcon className="size-3.5 text-success" aria-hidden />
            ) : (
              <CircleIcon className="size-3.5" aria-hidden />
            )}
            <span>{t(`profile.password.requirements.${key}`)}</span>
            <span className="sr-only">
              {met ? t('profile.password.requirementMet') : t('profile.password.requirementUnmet')}
            </span>
          </li>
        ))}
      </ul>
    </div>
  )
}
