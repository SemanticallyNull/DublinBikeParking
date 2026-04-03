import { useEffect } from 'react'
import { AppHeader } from './AppHeader'
import { MapView } from '../map/MapView'
import { StandDetailPanel } from '../panels/StandDetailPanel'
import { AddStandForm } from '../panels/AddStandForm'
import { BottomSheet } from '../ui/BottomSheet'
import { FloatingButton } from '../ui/FloatingButton'
import { NearestStandResult } from '../panels/NearestStandResult'
import type { StandFeature, QueryParams } from '../../types'
import type { NearestResult } from '../../hooks/useGeolocation'
import styles from './AppShell.module.css'

interface Props {
  features: StandFeature[]
  loading: boolean
  selectedStand: StandFeature | null
  placementMode: boolean
  queryParams: QueryParams
  userPosition: GeolocationPosition | null
  geoLoading: boolean
  nearestResult: NearestResult | null
  addPin: { lat: number; lng: number } | null
  onSelect: (feature: StandFeature) => void
  onClose: () => void
  onAddStand: () => void
  onLocate: () => void
  onDismissNearest: () => void
  onMapClick: (lat: number, lng: number) => void
  onAddSuccess: () => void
  onAddCancel: () => void
}

export function AppShell({
  features,
  loading,
  selectedStand,
  placementMode,
  queryParams,
  userPosition,
  geoLoading,
  nearestResult,
  addPin,
  onSelect,
  onClose,
  onAddStand,
  onLocate,
  onDismissNearest,
  onMapClick,
  onAddSuccess,
  onAddCancel,
}: Props) {
  // Prevent body scroll when panel is open on desktop
  useEffect(() => {
    if (selectedStand || addPin) {
      document.body.classList.add('panel-open')
    } else {
      document.body.classList.remove('panel-open')
    }
    return () => document.body.classList.remove('panel-open')
  }, [selectedStand, addPin])

  return (
    <div className={styles.shell}>
      {/* Placement mode banner */}
      {placementMode && (
        <div className={styles.placementBanner}>
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round">
            <circle cx="12" cy="12" r="9"/><circle cx="12" cy="12" r="3"/>
            <line x1="12" y1="2" x2="12" y2="5"/><line x1="12" y1="19" x2="12" y2="22"/>
            <line x1="2" y1="12" x2="5" y2="12"/><line x1="19" y1="12" x2="22" y2="12"/>
          </svg>
          Tap the map where the stand is
          <button className={styles.placementCancel} onClick={onAddStand}>Cancel</button>
        </div>
      )}

      <AppHeader
        standCount={features.length}
        placementMode={placementMode}
        onAddStand={onAddStand}
        onFindNearest={onLocate}
      />

      <div className={styles.body}>
        {loading && (
          <div className={styles.loadingPill}>Loading stands…</div>
        )}

        <div className={styles.mapArea}>
          <MapView
            features={features}
            selectedId={selectedStand?.properties.id ?? null}
            placementMode={placementMode}
            userPosition={userPosition}
            addPin={addPin}
            onSelect={onSelect}
            onMapClick={onMapClick}
          />

          {/* Nearest stand result card — overlays map */}
          {nearestResult && !selectedStand && (
            <NearestStandResult
              result={nearestResult}
              onSelect={() => onSelect(nearestResult.feature)}
              onDismiss={onDismissNearest}
            />
          )}

          {/* Geolocation FAB */}
          <FloatingButton loading={geoLoading} onClick={onLocate} />
        </div>

        {/* Desktop side panel */}
        {(selectedStand || addPin) && (
          <aside className={styles.sidePanel}>
            {addPin
              ? <AddStandForm lat={addPin.lat} lng={addPin.lng} onSuccess={onAddSuccess} onCancel={onAddCancel} />
              : <StandDetailPanel feature={selectedStand!} queryParams={queryParams} onClose={onClose} />
            }
          </aside>
        )}
      </div>

      {/* Mobile bottom sheet — hidden on desktop via CSS */}
      <div className={styles.mobileOnly}>
        <BottomSheet open={!!(selectedStand || addPin)} onClose={addPin ? onAddCancel : onClose}>
          {addPin
            ? <div className="bottom-sheet-content"><AddStandForm lat={addPin.lat} lng={addPin.lng} onSuccess={onAddSuccess} onCancel={onAddCancel} /></div>
            : selectedStand
              ? <div className="bottom-sheet-content"><StandDetailPanel feature={selectedStand} queryParams={queryParams} onClose={onClose} /></div>
              : null
          }
        </BottomSheet>
      </div>
    </div>
  )
}
