// frontend/src/App.tsx
import { useEffect, useState } from 'react'
import { AppShell } from './components/layout/AppShell'
import { useStands } from './hooks/useStands'
import { useQueryParams } from './hooks/useQueryParams'
import { useGeolocation } from './hooks/useGeolocation'
import type { StandFeature } from './types'

const GUIDE_SEEN_KEY = 'guide_seen'

export default function App() {
  const params = useQueryParams()
  const { features, loading, reload } = useStands({ verified: params.verified, pendingReview: params.pendingReview })
  const [selectedId, setSelectedId] = useState<string | null>(null)
  const [placementMode, setPlacementMode] = useState(false)
  const [addPin, setAddPin] = useState<{ lat: number; lng: number } | null>(null)
  const [guideOpen, setGuideOpen] = useState(false)

  const geo = useGeolocation(features)

  // Auto-open on first visit
  useEffect(() => {
    if (!localStorage.getItem(GUIDE_SEEN_KEY)) {
      localStorage.setItem(GUIDE_SEEN_KEY, 'true')
      setGuideOpen(true)
    }
  }, [])

  const selectedStand = selectedId
    ? (features.find(f => f.properties.id === selectedId) ?? null)
    : null

  function handleSelect(feature: StandFeature) {
    setSelectedId(feature.properties.id)
    setPlacementMode(false)
  }

  function handleClose() {
    setSelectedId(null)
  }

  function handleAddStand() {
    const next = !placementMode
    setPlacementMode(next)
    setSelectedId(null)
    if (!next) setAddPin(null)
  }

  function handleMapClick(lat: number, lng: number) {
    if (placementMode) {
      setAddPin({ lat, lng })
      setPlacementMode(false)
    }
  }

  function handleAddSuccess() {
    setAddPin(null)
    reload()
  }

  function handleAddCancel() {
    setAddPin(null)
  }

  return (
    <AppShell
      features={features}
      loading={loading}
      selectedStand={selectedStand}
      placementMode={placementMode}
      queryParams={params}
      userPosition={geo.position}
      geoLoading={geo.loading}
      nearestResult={geo.nearest}
      addPin={addPin}
      guideOpen={guideOpen}
      onSelect={handleSelect}
      onClose={handleClose}
      onAddStand={handleAddStand}
      onLocate={geo.locate}
      onDismissNearest={geo.clear}
      onMapClick={handleMapClick}
      onAddSuccess={handleAddSuccess}
      onAddCancel={handleAddCancel}
      onOpenGuide={() => setGuideOpen(true)}
      onCloseGuide={() => setGuideOpen(false)}
    />
  )
}
