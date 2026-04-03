import { Marker, Popup } from 'react-leaflet'
import { getPlacementIcon } from './markerIcons'
import { AddStandForm } from '../panels/AddStandForm'

interface Props {
  lat: number
  lng: number
  onSuccess: () => void
  onCancel: () => void
}

export function AddStandPin({ lat, lng, onSuccess, onCancel }: Props) {
  return (
    <Marker position={[lat, lng]} icon={getPlacementIcon()} interactive={false}>
      <Popup
        closeButton={false}
        autoClose={false}
        closeOnClick={false}
        minWidth={280}
        maxWidth={320}
        className="add-stand-popup"
      >
        <AddStandForm lat={lat} lng={lng} onSuccess={onSuccess} onCancel={onCancel} />
      </Popup>
    </Marker>
  )
}
