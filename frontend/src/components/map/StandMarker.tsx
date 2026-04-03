import { Marker } from 'react-leaflet'
import { getStandIcon } from './markerIcons'
import type { StandFeature } from '../../types'

interface Props {
  feature: StandFeature
  isActive: boolean
  onSelect: (feature: StandFeature) => void
}

export function StandMarker({ feature, isActive, onSelect }: Props) {
  const [lng, lat] = feature.geometry.coordinates
  const icon = getStandIcon(feature.properties.type, isActive)

  return (
    <Marker
      position={[lat, lng]}
      icon={icon}
      eventHandlers={{ click: () => onSelect(feature) }}
    />
  )
}
