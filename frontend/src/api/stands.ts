import type { StandCollection } from '../types'

interface FetchStandsParams {
  verified?: boolean
  pendingReview?: boolean
}

export async function fetchStands(params: FetchStandsParams = {}): Promise<StandCollection> {
  const query = new URLSearchParams()
  if (!params.verified) query.set('checked', 'unchecked')
  if (params.pendingReview) query.set('review', 'true')

  const url = `/api/v0/stand${query.size > 0 ? '?' + query.toString() : ''}`
  const res = await fetch(url)
  if (!res.ok) throw new Error(`Failed to fetch stands: ${res.status}`)
  return res.json()
}

export async function checkImageAvailability(): Promise<boolean> {
  try {
    const res = await fetch('/api/v0/image', { method: 'OPTIONS' })
    return res.status === 200
  } catch {
    return false
  }
}
