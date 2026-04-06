export function haversineM(lat1: number, lng1: number, lat2: number, lng2: number): number {
  const R = 6371000
  const dLat = (lat2 - lat1) * (Math.PI / 180)
  const dLng = (lng2 - lng1) * (Math.PI / 180)
  const a =
    Math.sin(dLat / 2) ** 2 +
    Math.cos(lat1 * (Math.PI / 180)) *
    Math.cos(lat2 * (Math.PI / 180)) *
    Math.sin(dLng / 2) ** 2
  return R * 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a))
}

export function bearingDeg(lat1: number, lng1: number, lat2: number, lng2: number): number {
  const toRad = Math.PI / 180
  const dLng = (lng2 - lng1) * toRad
  const y = Math.sin(dLng) * Math.cos(lat2 * toRad)
  const x =
    Math.cos(lat1 * toRad) * Math.sin(lat2 * toRad) -
    Math.sin(lat1 * toRad) * Math.cos(lat2 * toRad) * Math.cos(dLng)
  const brng = Math.atan2(y, x) * (180 / Math.PI)
  return (brng + 360) % 360
}

export type RelativeSide = 'left' | 'right' | 'ahead' | 'behind'

export function relativeSide(heading: number, bearing: number): RelativeSide {
  let diff = ((bearing - heading) + 360) % 360
  if (diff <= 45 || diff >= 315) return 'ahead'
  if (diff >= 135 && diff <= 225) return 'behind'
  if (diff > 45 && diff < 135) return 'right'
  return 'left'
}

export function compassDirection(bearing: number): string {
  const dirs = ['N', 'NE', 'E', 'SE', 'S', 'SW', 'W', 'NW']
  const index = Math.round(bearing / 45) % 8
  return dirs[index]
}
