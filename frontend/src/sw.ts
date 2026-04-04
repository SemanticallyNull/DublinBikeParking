/// <reference lib="webworker" />
import { cleanupOutdatedCaches, createHandlerBoundToURL, precacheAndRoute } from 'workbox-precaching'
import { NavigationRoute, registerRoute } from 'workbox-routing'
import { CacheFirst, StaleWhileRevalidate } from 'workbox-strategies'
import { CacheableResponsePlugin } from 'workbox-cacheable-response'
import { ExpirationPlugin } from 'workbox-expiration'

declare let self: ServiceWorkerGlobalScope

self.skipWaiting()

precacheAndRoute(self.__WB_MANIFEST)
cleanupOutdatedCaches()

// Warm the stands API cache during install so offline works immediately
// on the first visit, without needing a second page load.
self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open('stands-geojson').then(cache =>
      cache.add(new Request('/api/v0/stand', { credentials: 'same-origin' }))
    )
  )
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
