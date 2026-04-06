/// <reference lib="webworker" />
import { clientsClaim, skipWaiting } from 'workbox-core'
import { cleanupOutdatedCaches, precacheAndRoute } from 'workbox-precaching'
import { NavigationRoute, registerRoute, setCatchHandler } from 'workbox-routing'
import { StaleWhileRevalidate } from 'workbox-strategies'
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

// SPA navigation — look up index.html directly from the precache.
// We must strip the `redirected` flag from cached responses because
// Firefox on iOS rejects service-worker responses that carry it.
registerRoute(
  new NavigationRoute(async () => {
    console.log('[SW] navigation route handler')
    const cached = await caches.match('/index.html')
    if (cached) {
      console.log('[SW] serving index.html from cache')
      if (cached.redirected) {
        const body = await cached.blob()
        return new Response(body, {
          status: cached.status,
          statusText: cached.statusText,
          headers: cached.headers,
        })
      }
      return cached
    }
    console.log('[SW] index.html not in cache, fetching from network')
    return fetch('/index.html')
  })
)

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

// Map tiles — manual cache-first so we can log exactly what's happening.
// Fetch with mode:'cors' to avoid opaque responses (Chrome inflates their
// quota cost to ~7MB each, causing cache.put() to fail silently).
registerRoute(
  /basemaps\.cartocdn\.com/,
  async ({ request }) => {
    const cache = await caches.open('map-tiles')
    const cached = await cache.match(request)
    if (cached) {
      console.log('[SW] tile HIT:', request.url)
      return cached
    }
    console.log('[SW] tile MISS, fetching:', request.url)
    const response = await fetch(request.url, { mode: 'cors', credentials: 'omit' })
    console.log('[SW] tile fetch status:', response.status, response.type)
    if (response.ok) {
      cache.put(request, response.clone())
        .then(() => console.log('[SW] tile cached ok'))
        .catch(err => console.warn('[SW] tile cache.put failed:', err))
    }
    return response
  },
  'GET'
)

// Fallback: if any route handler throws, try the cache then give up gracefully
setCatchHandler(async ({ request }) => {
  console.warn('[SW] catch handler fired for', request.url)
  if (request.destination === 'document') {
    const cached = await caches.match('/index.html')
    if (cached) {
      if (cached.redirected) {
        const body = await cached.blob()
        return new Response(body, {
          status: cached.status,
          statusText: cached.statusText,
          headers: cached.headers,
        })
      }
      return cached
    }
  }
  return Response.error()
})
