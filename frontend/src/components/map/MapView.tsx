import { useRef } from 'react'
import { MapContainer, TileLayer, useMapEvents } from 'react-leaflet'
import type { Map as LeafletMap, LeafletMouseEvent } from 'leaflet'
import { StandsLayer } from './StandsLayer'
import { GeolocationMarker } from './GeolocationMarker'
import { AddStandPin } from './AddStandPin'
import { MapHashSync } from '../../hooks/useMapHash'
import type { StandFeature } from '../../types'
import '../../styles/map.css'

const DUBLIN_CENTER: [number, number] = [53.3441, -6.2675]
const CARTO_VOYAGER = 'https://{s}.basemaps.cartocdn.com/rastertiles/voyager/{z}/{x}/{y}{r}.png'
const ATTRIBUTION = '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> &copy; <a href="https://carto.com/attributions">CartoDB</a>'

interface Props {
  features: StandFeature[]
  selectedId: string | null
  placementMode: boolean
  userPosition?: GeolocationPosition | null
  addPin?: { lat: number; lng: number } | null
  onSelect: (feature: StandFeature) => void
  onMapClick: (lat: number, lng: number) => void
}

function MapEvents({ placementMode, onMapClick }: { placementMode: boolean; onMapClick: (lat: number, lng: number) => void }) {
  useMapEvents({
    click(e: LeafletMouseEvent) {
      if (placementMode) onMapClick(e.latlng.lat, e.latlng.lng)
    },
  })
  return null
}

export function MapView({ features, selectedId, placementMode, userPosition, addPin, onSelect, onMapClick }: Props) {
  const mapRef = useRef<LeafletMap | null>(null)

  return (
    <MapContainer
      center={DUBLIN_CENTER}
      zoom={14}
      zoomControl={false}
      style={{ width: '100%', height: '100%' }}
      className={placementMode ? 'placement-mode' : ''}
      ref={mapRef}
    >
      <TileLayer url={CARTO_VOYAGER} attribution={ATTRIBUTION} subdomains="abcd" maxZoom={20} />
      <StandsLayer features={features} selectedId={selectedId} onSelect={onSelect} />
      {userPosition && <GeolocationMarker position={userPosition} />}
      {addPin && <AddStandPin lat={addPin.lat} lng={addPin.lng} />}
      <MapEvents placementMode={placementMode} onMapClick={onMapClick} />
      <MapHashSync />
    </MapContainer>
  )
}
