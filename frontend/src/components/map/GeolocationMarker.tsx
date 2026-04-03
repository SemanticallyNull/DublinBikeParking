import { Circle, Marker } from 'react-leaflet'
import L from 'leaflet'

interface Props {
  position: GeolocationPosition
}

const userIcon = L.divIcon({
  className: '',
  html: `<div style="
    width: 16px;
    height: 16px;
    background: #6C96BB;
    border: 3px solid white;
    border-radius: 50%;
    box-shadow: 0 2px 8px rgba(14,32,82,0.4);
  "></div>`,
  iconSize: [16, 16],
  iconAnchor: [8, 8],
})

export function GeolocationMarker({ position }: Props) {
  const { latitude, longitude, accuracy } = position.coords
  return (
    <>
      <Circle
        center={[latitude, longitude]}
        radius={Math.min(accuracy, 500)}
        pathOptions={{ color: '#6C96BB', fillColor: '#6C96BB', fillOpacity: 0.08, weight: 1.5, opacity: 0.5 }}
      />
      <Marker position={[latitude, longitude]} icon={userIcon} />
    </>
  )
}
