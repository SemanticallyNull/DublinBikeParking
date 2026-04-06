import { useState, useCallback } from 'react'
import { verifyStand, reportMissingAuth } from '../api/actions'

export interface SessionStats {
  verified: number
  missing: number
  skipped: number
}

export function useVerifySession(password: string) {
  const [verifiedIds, setVerifiedIds] = useState<Set<string>>(new Set())
  const [missingIds, setMissingIds] = useState<Set<string>>(new Set())
  const [skippedIds, setSkippedIds] = useState<Set<string>>(new Set())

  const allExcluded = new Set([...verifiedIds, ...missingIds, ...skippedIds])

  const verify = useCallback(async (id: string) => {
    setVerifiedIds(s => new Set(s).add(id))
    try {
      await verifyStand(id, password)
    } catch (err) {
      console.error('Failed to verify stand:', err)
    }
  }, [password])

  const markMissing = useCallback(async (id: string) => {
    setMissingIds(s => new Set(s).add(id))
    try {
      await reportMissingAuth(id, password)
    } catch (err) {
      console.error('Failed to report missing:', err)
    }
  }, [password])

  const skip = useCallback((id: string) => {
    setSkippedIds(s => new Set(s).add(id))
  }, [])

  const sessionStats: SessionStats = {
    verified: verifiedIds.size,
    missing: missingIds.size,
    skipped: skippedIds.size,
  }

  return { verifiedIds, missingIds, skippedIds, allExcluded, verify, markMissing, skip, sessionStats }
}
