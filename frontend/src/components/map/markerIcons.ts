import L from 'leaflet'

const TOASTER_TYPES = ['Wheel Only']

function buildIcon(type: string, isActive: boolean): L.DivIcon {
  const size = isActive ? 32 : 26
  const r = size / 2
  const isToaster = TOASTER_TYPES.includes(type)

  const svg = isToaster
    ? `<svg width="${size}" height="${size}" viewBox="0 0 ${size} ${size}" xmlns="http://www.w3.org/2000/svg">
        <circle cx="${r}" cy="${r}" r="${r - 2}" fill="#8C98B0" stroke="#C4CEDF" stroke-width="2"/>
        <circle cx="${r}" cy="${r}" r="${r * 0.32}" fill="white" opacity="0.7"/>
       </svg>`
    : `<svg width="${size}" height="${size}" viewBox="0 0 ${size} ${size}" xmlns="http://www.w3.org/2000/svg">
        <circle cx="${r}" cy="${r}" r="${r - 2}" fill="#0E2052" stroke="#5A94C8" stroke-width="2"/>
        <circle cx="${r}" cy="${r}" r="${r * 0.32}" fill="#5A94C8"/>
       </svg>`

  const className = ['stand-marker', isToaster ? 'toaster' : '', isActive ? 'active' : '']
    .filter(Boolean)
    .join(' ')

  return L.divIcon({ html: svg, className, iconSize: [size, size], iconAnchor: [r, r] })
}

// Cache icons to avoid recreating on every render
const cache = new Map<string, L.DivIcon>()

export function getStandIcon(type: string, isActive = false): L.DivIcon {
  const key = `${type}:${isActive}`
  if (!cache.has(key)) cache.set(key, buildIcon(type, isActive))
  return cache.get(key)!
}

export function getPlacementIcon(): L.DivIcon {
  return L.divIcon({
    html: `<svg width="32" height="32" viewBox="0 0 32 32" xmlns="http://www.w3.org/2000/svg">
      <circle cx="16" cy="16" r="14" fill="#0E2052" stroke="#5A94C8" stroke-width="2"/>
      <line x1="16" y1="10" x2="16" y2="22" stroke="white" stroke-width="2.5" stroke-linecap="round"/>
      <line x1="10" y1="16" x2="22" y2="16" stroke="white" stroke-width="2.5" stroke-linecap="round"/>
    </svg>`,
    className: 'stand-marker placement-pin',
    iconSize: [32, 32],
    iconAnchor: [16, 16],
  })
}
