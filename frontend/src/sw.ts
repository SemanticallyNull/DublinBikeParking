/// <reference lib="webworker" />
import { clientsClaim, skipWaiting } from 'workbox-core'
import { cleanupOutdatedCaches, createHandlerBoundToURL, precacheAndRoute } from 'workbox-precaching'
import { NavigationRoute, registerRoute } from 'workbox-routing'
import { CacheFirst, StaleWhileRevalidate } from 'workbox-strategies'
import { CacheableResponsePlugin } from 'workbox-cacheable-response'
import { ExpirationPlugin } from 'workbox-expiration'

declare let self: ServiceWorkerGlobalScope

skipWaiting()
clientsClaim()

precacheAndRoute(self.__WB_MANIFEST)
cleanupOutdatedCaches()

self.addEventListener('install', (event) => {
  console.log('[SW] install')
  event.waitUntil(
    caches.open('stands-geojson')
      .then(cache => {
        console.log('[SW] warming /api/v0/stand...')
        return cache.add(new Request('/api/v0/stand', { credentials: 'same-origin' }))
      })
      .then(() => console.log('[SW] /api/v0/stand cached ok'))
      .catch(err => console.warn('[SW] API warm failed (non-fatal):', err))
  )
})

self.addEventListener('activate', () => {
  console.log('[SW] activate — in control')
})

self.addEventListener('fetch', (event: FetchEvent) => {
  console.log('[SW] fetch', event.request.method, event.request.url)
})

// SPA navigation fallback
registerRoute(new NavigationRoute(createHandlerBoundToURL('index.html')))

// Stands API — serve from cache while revalidating in background
registerRoute(
  /\/api\/v0\/stand$/,
  new StaleWhileRevalidate({
    cacheName: 'stands-geojson',
    plugins: [
      new ExpirationPlugin({ maxAgeSeconds: 7 * 24 * 60 * 60 }),
      new CacheableResponsePlugin({ statuses: [0, 200] }),
    ],
  }),
  'GET'
)

// Map tiles — cache-first, keep up to 500 tiles for 7 days
registerRoute(
  /basemaps\.cartocdn\.com/,
  new CacheFirst({
    cacheName: 'map-tiles',
    plugins: [
      new ExpirationPlugin({ maxEntries: 500, maxAgeSeconds: 7 * 24 * 60 * 60 }),
      new CacheableResponsePlugin({ statuses: [0, 200] }),
    ],
  }),
  'GET'
)
