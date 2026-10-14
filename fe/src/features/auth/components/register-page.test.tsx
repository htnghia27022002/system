import { render, screen } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import i18n from '@/config/i18n'

const { state } = vi.hoisted(() => ({ state: { open: true } }))

vi.mock('@/hooks/use-app-info', () => ({ useRegistrationEnabled: () => state.open }))
vi.mock('./register-form', () => ({ RegisterForm: () => <form aria-label="register" /> }))

import { RegisterPage } from './register-page'

describe('RegisterPage', () => {
  beforeEach(async () => {
    await i18n.changeLanguage('en')
  })

  it('shows the form while sign-up is open', () => {
    state.open = true
    render(<RegisterPage />)
    expect(screen.getByRole('form', { name: 'register' })).toBeInTheDocument()
  })

  it('explains that sign-up is closed', () => {
    state.open = false
    render(<RegisterPage />)
    expect(screen.queryByRole('form', { name: 'register' })).not.toBeInTheDocument()
    expect(screen.getByRole('status')).toHaveTextContent('Sign-up is closed')
  })
})
