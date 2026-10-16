import { describe, expect, it } from 'vitest'

import { evaluatePassword } from './password-strength'
import { generateSecurePassword } from './generate-secure-password'

describe('evaluatePassword', () => {
  it('scores empty as 0 and short passwords as weak', () => {
    expect(evaluatePassword('').score).toBe(0)
    expect(evaluatePassword('Ab1!').score).toBe(1)
  })

  it('rewards variety and length', () => {
    expect(evaluatePassword('abcdefgh').score).toBe(1)
    expect(evaluatePassword('abcdefg1').score).toBe(1)
    expect(evaluatePassword('Abcdefg1').score).toBe(2)
    expect(evaluatePassword('Abcdefg1!').score).toBe(3)
    expect(evaluatePassword('Abcdefg1!xyz').score).toBe(4)
  })

  it('reports individual checks', () => {
    expect(evaluatePassword('Abcdefg1').checks).toEqual({
      length: true,
      mixedCase: true,
      number: true,
      symbol: false,
    })
  })

  it('rates generated passwords as strong', () => {
    expect(evaluatePassword(generateSecurePassword()).score).toBe(4)
  })
})
