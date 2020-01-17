import L from 'leaflet';

const StandIcon = L.Icon.extend({
  options: {
    iconSize: [22, 22],
  },
});
const standIcons = {
  'Sheffield Stand': new StandIcon({ iconUrl: 'icons/stand_icon_sheffield.png' }),
  'Wheel Only': new StandIcon({ iconUrl: 'icons/stand_icon_wheel_only.png' }),
  Hoop: new StandIcon({ iconUrl: 'icons/stand_icon_hoop.png' }),
  'Stainless Steel Curved': new StandIcon({ iconUrl: 'icons/stand_icon_ssc.png' }),
  Railing: new StandIcon({ iconUrl: 'icons/stand_icon_railing.png' }),
  DublinBikes: new StandIcon({ iconUrl: 'icons/station_icon_db.png' }),
};

export default standIcons;
