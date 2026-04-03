import type { QueryParams } from '../types'

export function useQueryParams(): QueryParams {
  const p = new URLSearchParams(window.location.search)
  return {
    showIDs:       p.get('showIDs') === 'true',
    verified:      p.get('checked') !== 'unchecked',
    pendingReview: p.get('review') === 'true',
  }
}
