import MarkerClusterGroup from 'react-leaflet-cluster'
import { StandMarker } from './StandMarker'
import type { StandFeature } from '../../types'

interface Props {
  features: StandFeature[]
  selectedId: string | null
  onSelect: (feature: StandFeature) => void
}

export function StandsLayer({ features, selectedId, onSelect }: Props) {
  return (
    <MarkerClusterGroup maxClusterRadius={50} disableClusteringAtZoom={18} chunkedLoading>
      {features.map(f => (
        <StandMarker
          key={f.properties.id}
          feature={f}
          isActive={f.properties.id === selectedId}
          onSelect={onSelect}
        />
      ))}
    </MarkerClusterGroup>
  )
}
