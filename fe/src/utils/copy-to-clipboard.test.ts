import { afterEach, describe, expect, it, vi } from 'vitest'

import { copyToClipboard } from './copy-to-clipboard'

function stubClipboard(value: Clipboard | undefined) {
  Object.defineProperty(navigator, 'clipboard', {
    configurable: true,
    value,
  })
}

function stubExecCommand(result: boolean) {
  const exec = vi.fn().mockReturnValue(result)
  Object.defineProperty(document, 'execCommand', {
    configurable: true,
    value: exec,
    writable: true,
  })
  return exec
}

describe('copyToClipboard', () => {
  const originalClipboard = navigator.clipboard

  afterEach(() => {
    stubClipboard(originalClipboard)
    vi.restoreAllMocks()
  })

  it('writes via Clipboard API when writeText is available', async () => {
    const writeText = vi.fn().mockResolvedValue(undefined)
    stubClipboard({ writeText } as unknown as Clipboard)

    await copyToClipboard('https://example.com/hook')

    expect(writeText).toHaveBeenCalledWith('https://example.com/hook')
  })

  it('falls back to execCommand when clipboard is undefined', async () => {
    stubClipboard(undefined)
    const exec = stubExecCommand(true)

    await copyToClipboard('http://system.local:8080/tools/webhooks/abc')

    expect(exec).toHaveBeenCalledWith('copy')
  })

  it('falls back to execCommand when writeText rejects', async () => {
    const writeText = vi.fn().mockRejectedValue(new Error('denied'))
    stubClipboard({ writeText } as unknown as Clipboard)
    const exec = stubExecCommand(true)

    await copyToClipboard('fallback-text')

    expect(writeText).toHaveBeenCalled()
    expect(exec).toHaveBeenCalledWith('copy')
  })

  it('throws when neither clipboard API nor execCommand succeeds', async () => {
    stubClipboard(undefined)
    stubExecCommand(false)

    await expect(copyToClipboard('nope')).rejects.toThrow(
      'Clipboard copy is not available',
    )
  })
})
