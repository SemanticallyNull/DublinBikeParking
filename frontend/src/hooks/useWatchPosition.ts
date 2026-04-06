import { useState, useEffect, useRef } from 'react'
import type { StandFeature } from '../types'
import { haversineM, bearingDeg, relativeSide, compassDirection } from '../utils/geo'
import type { RelativeSide } from '../utils/geo'

export interface NearestStand {
  feature: StandFeature
  distanceM: number
  bearing: number
  side: RelativeSide
  compass: string
}

export interface WatchPositionState {
  position: GeolocationPosition | null
  heading: number | null
  nearest: NearestStand | null
  error: string | null
}

export function useWatchPosition(
  features: StandFeature[],
  excludedIds: Set<string>,
  active: boolean,
) {
  const [state, setState] = useState<WatchPositionState>({
    position: null,
    heading: null,
    nearest: null,
    error: null,
  })
  const featuresRef = useRef(features)
  const excludedRef = useRef(excludedIds)
  featuresRef.current = features
  excludedRef.current = excludedIds

  useEffect(() => {
    if (!active) return
    if (!navigator.geolocation) {
      setState(s => ({ ...s, error: 'Geolocation not supported' }))
      return
    }

    const watchId = navigator.geolocation.watchPosition(
      pos => {
        const { latitude, longitude, heading } = pos.coords
        const h = heading != null && !isNaN(heading) && heading >= 0 ? heading : null

        let nearest: NearestStand | null = null
        for (const feature of featuresRef.current) {
          if (excludedRef.current.has(feature.properties.id)) continue
          const [lng, lat] = feature.geometry.coordinates
          const d = haversineM(latitude, longitude, lat, lng)
          if (d <= 2000 && (nearest === null || d < nearest.distanceM)) {
            const brng = bearingDeg(latitude, longitude, lat, lng)
            const side = h != null ? relativeSide(h, brng) : 'ahead'
            nearest = {
              feature,
              distanceM: d,
              bearing: brng,
              side,
              compass: compassDirection(brng),
            }
          }
        }

        setState({ position: pos, heading: h, nearest, error: null })
      },
      err => {
        let msg = 'Could not get your location.'
        if (err.code === err.PERMISSION_DENIED) msg = 'Location access was denied.'
        if (err.code === err.POSITION_UNAVAILABLE) msg = 'Location unavailable.'
        if (err.code === err.TIMEOUT) msg = 'Location request timed out.'
        setState(s => ({ ...s, error: msg }))
      },
      { enableHighAccuracy: true, timeout: 10000, maximumAge: 2000 },
    )

    return () => navigator.geolocation.clearWatch(watchId)
  }, [active])

  return state
}
