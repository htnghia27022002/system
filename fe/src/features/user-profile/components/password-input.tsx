'use client'

import { EyeIcon, EyeOffIcon } from 'lucide-react'
import type { ComponentProps, ReactNode } from 'react'
import { useTranslation } from 'react-i18next'

import {
  InputGroup,
  InputGroupAddon,
  InputGroupButton,
  InputGroupInput,
} from '@/components/ui/input-group'

type PasswordInputProps = ComponentProps<typeof InputGroupInput> & {
  revealed: boolean
  onToggleReveal: () => void
  /** Extra inline buttons shown before the show/hide toggle (e.g. copy). */
  extraActions?: ReactNode
}

/** Password field with a show/hide toggle; shared by profile and admin user forms. */
export function PasswordInput({
  revealed,
  onToggleReveal,
  extraActions,
  ...props
}: PasswordInputProps) {
  const { t } = useTranslation('admin')
  return (
    <InputGroup className="h-9">
      <InputGroupInput type={revealed ? 'text' : 'password'} spellCheck={false} {...props} />
      <InputGroupAddon align="inline-end">
        {extraActions}
        <InputGroupButton
          size="icon-xs"
          aria-label={
            revealed ? t('profile.actions.hidePassword') : t('profile.actions.showPassword')
          }
          aria-pressed={revealed}
          onClick={onToggleReveal}
        >
          {revealed ? <EyeOffIcon /> : <EyeIcon />}
        </InputGroupButton>
      </InputGroupAddon>
    </InputGroup>
  )
}
