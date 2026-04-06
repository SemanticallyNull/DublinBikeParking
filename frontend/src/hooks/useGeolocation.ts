import { useState, useCallback } from 'react'
import type { StandFeature } from '../types'
import { haversineM } from '../utils/geo'

export interface GeolocationState {
  position: GeolocationPosition | null
  nearest: NearestResult | null
  loading: boolean
  error: string | null
}

export interface NearestResult {
  feature: StandFeature
  distanceM: number
}

export function useGeolocation(features: StandFeature[]) {
  const [state, setState] = useState<GeolocationState>({
    position: null,
    nearest: null,
    loading: false,
    error: null,
  })

  const locate = useCallback(() => {
    if (!navigator.geolocation) {
      setState(s => ({ ...s, error: 'Geolocation is not supported by your browser.' }))
      return
    }

    setState(s => ({ ...s, loading: true, error: null }))

    navigator.geolocation.getCurrentPosition(
      pos => {
        const { latitude, longitude } = pos.coords

        // Find nearest stand within 2km
        let nearest: NearestResult | null = null
        for (const feature of features) {
          const [lng, lat] = feature.geometry.coordinates
          const d = haversineM(latitude, longitude, lat, lng)
          if (d <= 2000 && (nearest === null || d < nearest.distanceM)) {
            nearest = { feature, distanceM: d }
          }
        }

        setState({ position: pos, nearest, loading: false, error: null })
      },
      err => {
        let msg = 'Could not get your location.'
        if (err.code === err.PERMISSION_DENIED) msg = 'Location access was denied.'
        if (err.code === err.POSITION_UNAVAILABLE) msg = 'Location unavailable.'
        if (err.code === err.TIMEOUT) msg = 'Location request timed out.'
        setState(s => ({ ...s, position: null, nearest: null, loading: false, error: msg }))
      },
      { enableHighAccuracy: true, timeout: 10000, maximumAge: 30000 }
    )
  }, [features])

  const clear = useCallback(() => {
    setState({ position: null, nearest: null, loading: false, error: null })
  }, [])

  return { ...state, locate, clear }
}
