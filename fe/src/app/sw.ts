/// <reference lib="webworker" />

import { defaultCache } from '@serwist/next/worker'
import type { PrecacheEntry, SerwistGlobalConfig } from 'serwist'
import { NetworkFirst, Serwist } from 'serwist'

declare global {
  interface WorkerGlobalScope extends SerwistGlobalConfig {
    __SW_MANIFEST: (PrecacheEntry | string)[] | undefined
  }
}

declare const self: ServiceWorkerGlobalScope

const serwist = new Serwist({
  precacheEntries: self.__SW_MANIFEST,
  skipWaiting: true,
  clientsClaim: true,
  navigationPreload: true,
  runtimeCaching: [
    ...defaultCache,
    {
      matcher: ({ url }) => url.pathname.startsWith('/api/'),
      handler: new NetworkFirst({
        cacheName: 'api-cache',
        networkTimeoutSeconds: 8,
        plugins: [
          {
            handlerDidError: async () => {
              return (await caches.match('/offline.html')) ?? Response.error()
            },
          },
        ],
      }),
    },
  ],
  fallbacks: {
    entries: [
      {
        url: '/offline.html',
        matcher: ({ request }) => request.mode === 'navigate',
      },
    ],
  },
})

serwist.addEventListeners()
