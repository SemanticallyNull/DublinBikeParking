import { useState, useEffect } from 'react'
import { fetchStands } from '../api/stands'
import type { StandFeature, QueryParams } from '../types'

interface UseStandsResult {
  features: StandFeature[]
  loading: boolean
  error: string | null
  reload: () => void
}

export function useStands(params: Partial<QueryParams> = {}): UseStandsResult {
  const [features, setFeatures] = useState<StandFeature[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [tick, setTick] = useState(0)

  useEffect(() => {
    let cancelled = false
    setLoading(true)
    setError(null)

    fetchStands({ verified: params.verified, pendingReview: params.pendingReview })
      .then(col => {
        if (!cancelled) {
          setFeatures(col.features ?? [])
          setLoading(false)
        }
      })
      .catch(err => {
        if (!cancelled) {
          setError(err.message)
          setLoading(false)
        }
      })

    return () => { cancelled = true }
  }, [params.verified, params.pendingReview, tick])

  return { features, loading, error, reload: () => setTick(t => t + 1) }
}
