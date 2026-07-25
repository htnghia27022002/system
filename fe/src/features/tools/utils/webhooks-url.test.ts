import { describe, expect, it } from 'vitest'

import {
  buildWebhookPublicPath,
  buildWebhookPublicUrl,
  captureApiDestination,
  isWebhookPublicUuid,
  normalizeMethodFilter,
} from './webhooks-url'

describe('webhooks URL helpers', () => {
  const uuid = '550e8400-e29b-41d4-a716-446655440000'

  it('validates UUID-shaped public ids', () => {
    expect(isWebhookPublicUuid(uuid)).toBe(true)
    expect(isWebhookPublicUuid('not-a-uuid')).toBe(false)
    expect(isWebhookPublicUuid('')).toBe(false)
  })

  it('builds product public path', () => {
    expect(buildWebhookPublicPath(uuid)).toBe(`/tools/webhooks/${uuid}`)
  })

  it('builds absolute display URL from origin + publicPath', () => {
    expect(
      buildWebhookPublicUrl('http://system.local:8080', `/tools/webhooks/${uuid}`),
    ).toBe(`http://system.local:8080/tools/webhooks/${uuid}`)
  })

  it('builds absolute display URL from origin + raw uuid', () => {
    expect(buildWebhookPublicUrl('https://example.com/', uuid)).toBe(
      `https://example.com/tools/webhooks/${uuid}`,
    )
  })

  it('normalizes method filters for the list API', () => {
    expect(normalizeMethodFilter('ALL')).toBe('')
    expect(normalizeMethodFilter('all')).toBe('')
    expect(normalizeMethodFilter(' post ')).toBe('POST')
    expect(normalizeMethodFilter(undefined)).toBe('')
  })

  it('derives capture rewrite destination from API base', () => {
    expect(captureApiDestination('http://localhost:8080/api')).toBe(
      'http://localhost:8080/api/webhooks/capture/:uuid',
    )
    expect(captureApiDestination('http://localhost:8080/api/', 'x')).toBe(
      'http://localhost:8080/api/webhooks/capture/x',
    )
  })
})
