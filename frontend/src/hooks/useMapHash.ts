import { useEffect } from 'react'
import { useMap, useMapEvents } from 'react-leaflet'

// Preserves the same #zoom/lat/lng format as the existing leaflet-hash library
// so existing bookmarks continue to work

function parseHash(): { zoom: number; lat: number; lng: number } | null {
  const hash = window.location.hash.slice(1)
  const parts = hash.split('/')
  if (parts.length !== 3) return null
  const zoom = parseInt(parts[0], 10)
  const lat = parseFloat(parts[1])
  const lng = parseFloat(parts[2])
  if (isNaN(zoom) || isNaN(lat) || isNaN(lng)) return null
  return { zoom, lat, lng }
}

function formatHash(zoom: number, lat: number, lng: number): string {
  return `#${zoom}/${lat.toFixed(5)}/${lng.toFixed(5)}`
}

export function MapHashSync() {
  const map = useMap()

  // On mount, if there's a hash in the URL, jump to that position
  useEffect(() => {
    const pos = parseHash()
    if (pos) map.setView([pos.lat, pos.lng], pos.zoom)
  }, [map])

  // On every map move, update the hash
  useMapEvents({
    moveend() {
      const c = map.getCenter()
      const z = map.getZoom()
      window.location.replace(formatHash(z, c.lat, c.lng))
    },
  })

  return null
}
