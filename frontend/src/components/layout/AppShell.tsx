// frontend/src/components/layout/AppShell.tsx
import { useEffect } from 'react'
import { AppHeader } from './AppHeader'
import { MapView } from '../map/MapView'
import { StandDetailPanel } from '../panels/StandDetailPanel'
import { AddStandForm } from '../panels/AddStandForm'
import { BottomSheet } from '../ui/BottomSheet'
import { FloatingButton } from '../ui/FloatingButton'
import { NearestStandResult } from '../panels/NearestStandResult'
import { LockingGuide } from '../guide/LockingGuide'
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
  guideOpen: boolean
  onSelect: (feature: StandFeature) => void
  onClose: () => void
  onAddStand: () => void
  onLocate: () => void
  onDismissNearest: () => void
  onMapClick: (lat: number, lng: number) => void
  onAddSuccess: () => void
  onAddCancel: () => void
  onOpenGuide: () => void
  onCloseGuide: () => void
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
  guideOpen,
  onSelect,
  onClose,
  onAddStand,
  onLocate,
  onDismissNearest,
  onMapClick,
  onAddSuccess,
  onAddCancel,
  onOpenGuide,
  onCloseGuide,
}: Props) {
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
        onOpenGuide={onOpenGuide}
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

          {nearestResult && !selectedStand && (
            <NearestStandResult
              result={nearestResult}
              onSelect={() => onSelect(nearestResult.feature)}
              onDismiss={onDismissNearest}
            />
          )}

          <FloatingButton loading={geoLoading} onClick={onLocate} />
        </div>

        {(selectedStand || addPin) && (
          <aside className={styles.sidePanel}>
            {addPin
              ? <AddStandForm lat={addPin.lat} lng={addPin.lng} onSuccess={onAddSuccess} onCancel={onAddCancel} />
              : <StandDetailPanel feature={selectedStand!} queryParams={queryParams} onClose={onClose} onOpenGuide={onOpenGuide} />
            }
          </aside>
        )}
      </div>

      <div className={styles.mobileOnly}>
        <BottomSheet open={!!(selectedStand || addPin)} onClose={addPin ? onAddCancel : onClose}>
          {addPin
            ? <div className="bottom-sheet-content"><AddStandForm lat={addPin.lat} lng={addPin.lng} onSuccess={onAddSuccess} onCancel={onAddCancel} /></div>
            : selectedStand
              ? <div className="bottom-sheet-content"><StandDetailPanel feature={selectedStand} queryParams={queryParams} onClose={onClose} onOpenGuide={onOpenGuide} /></div>
              : null
          }
        </BottomSheet>
      </div>

      <LockingGuide open={guideOpen} onClose={onCloseGuide} />
    </div>
  )
}
