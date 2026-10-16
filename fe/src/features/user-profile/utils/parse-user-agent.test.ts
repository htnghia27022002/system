import { describe, expect, it } from 'vitest'

import { parseUserAgent } from './parse-user-agent'

describe('parseUserAgent', () => {
  it('detects Chrome on Windows desktop', () => {
    expect(
      parseUserAgent(
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36',
      ),
    ).toEqual({ browser: 'Chrome', os: 'Windows', deviceType: 'desktop' })
  })

  it('prefers Edge over Chrome', () => {
    expect(
      parseUserAgent(
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36 Edg/141.0.0.0',
      ).browser,
    ).toBe('Edge')
  })

  it('detects Safari on iPhone as mobile', () => {
    expect(
      parseUserAgent(
        'Mozilla/5.0 (iPhone; CPU iPhone OS 18_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.0 Mobile/15E148 Safari/604.1',
      ),
    ).toEqual({ browser: 'Safari', os: 'iOS', deviceType: 'mobile' })
  })

  it('treats Android without Mobi as tablet', () => {
    expect(
      parseUserAgent(
        'Mozilla/5.0 (Linux; Android 14; SM-X710) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36',
      ).deviceType,
    ).toBe('tablet')
  })

  it('returns nulls for missing or unknown agents', () => {
    expect(parseUserAgent(null)).toEqual({ browser: null, os: null, deviceType: 'desktop' })
    expect(parseUserAgent('curl/8.5.0')).toEqual({ browser: null, os: null, deviceType: 'desktop' })
  })
})
