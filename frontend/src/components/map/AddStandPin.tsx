import { Marker } from 'react-leaflet'
import { getPlacementIcon } from './markerIcons'

interface Props {
  lat: number
  lng: number
}

export function AddStandPin({ lat, lng }: Props) {
  return <Marker position={[lat, lng]} icon={getPlacementIcon()} interactive={false} />
}
