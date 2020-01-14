<template>
  <l-map class='map' :zoom='zoom' :center='center' :options='mapOptions'>
    <l-tile-layer :url='url'></l-tile-layer>
    <l-marker-cluster :options='clusterOptions'>
      <l-geo-json
        v-for='geojson in geojsons'
        :key='geojson.type'
        :geojson='geojson'
        :options='geojsonOptions'>
      </l-geo-json>
    </l-marker-cluster>
  </l-map>
</template>

<script>
import L from 'leaflet';
import { LMap, LTileLayer, LGeoJson } from 'vue2-leaflet';
import LMarkerCluster from 'vue2-leaflet-markercluster';
import axios from 'axios';

const geojsons = [];

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

function standType(f, l) {
  const noStands = f.properties.numberOfStands > 0 ? f.properties.numberOfStands : 'Unknown';

  let unverifiedMessage = '';
  if (!f.properties.verified) {
    unverifiedMessage = '<tr><td colspan=2><small>This is a user submitted stand we have not verified yet.</small></td></tr>';
  }

  const urlParams = new URLSearchParams(window.location.search);
  const standIDParam = urlParams.get('showIDs');
  let standID = '';
  if (standIDParam === 'true') {
    standID = `<tr><td><b>Stand ID</b></td><td>${f.properties.id}</td></tr>`;
  }

  l.setIcon(standIcons[f.properties.type]).bindPopup(`<table>
    <tr><td><b>Type of Stands</b></td><td>${f.properties.type}</td></tr>
    <tr><td><b>Location</b></td><td>${f.properties.name}</td></tr>
    <tr><td><b>Number of Stands</b></td><td>${noStands}</td></tr>
    ${standID}
    ${unverifiedMessage}
    </table>`);
}

export default {
  name: 'Map',
  components: {
    'l-map': LMap,
    'l-tile-layer': LTileLayer,
    'l-geo-json': LGeoJson,
    'l-marker-cluster': LMarkerCluster,
  },
  data() {
    return {
      url: 'https://maps.wikimedia.org/osm-intl/{z}/{x}/{y}.png',
      zoom: 12,
      center: [53.34587706183283, -6.267379465984959],
      clusterOptions: {
        maxClusterRadius: 50,
        disableClusteringAtZoom: 18,
      },
      mapOptions: {
        maxZoom: 18,
      },
      geojsonOptions: {
        onEachFeature: standType,
      },
      geojsons,
    };
  },
  created() {
    axios.get('/api/v0/stand?dublinbikes=off').then((response) => {
      this.geojsons.push(response.data);
    });
  },
};
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style scoped>
@import '~leaflet.markercluster/dist/MarkerCluster.css';
@import '~leaflet.markercluster/dist/MarkerCluster.Default.css';

h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
.map {
  width: 100%;
  height: 100%;
}
</style>
