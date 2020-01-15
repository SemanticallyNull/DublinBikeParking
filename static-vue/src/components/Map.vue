<template>
  <l-map class='map' :zoom='zoom' :center='center' :options='mapOptions'>
    <l-tile-layer :url='url' :options="tileLayerOptions"></l-tile-layer>
    <l-marker-cluster ref="clusterRef" :options='clusterOptions'>
      <l-geo-json
        v-for='geojson in geojsons'
        :key='geojson.features[0].properties.id'
        :geojson='geojson'
        :options='geojsonOptions'>
      </l-geo-json>
    </l-marker-cluster>
    <l-control>
      <div class="btn-group leaflet-selector">
        <button v-for="bikeType in Object.keys(bikeTypes)" v-bind:key="bikeType"
                class="btn btn-sm"
                :class="{
                  'btn-dark': bikeTypeIsActive(bikeType),
                  'btn-outline-dark': !bikeTypeIsActive(bikeType) }"
                @click="changeBikeType(bikeType)"
        >
          {{bikeType}}
        </button>
      </div>
    </l-control>
  </l-map>
</template>

<script>
import L from 'leaflet';
import {
  LMap, LTileLayer, LGeoJson, LControl,
} from 'vue2-leaflet';
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

const bikeTypes = {
  'My Bike': {
    query: 'dublinbikes=off',
  },
  DublinBikes: {
    query: 'dublinbikes=only',
  },
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
    'l-control': LControl,
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
      mapOptions: {},
      tileLayerOptions: {
        maxZoom: 18,
        attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors, Stand Data: <a href="https://bleeperbike.com/">Bleeper Bikes</a>, <a href="https://data.smartdublin.ie/dataset/dcc_public_cycle_parking_stands">DCC</a> DublinBikes: <a href="https://developer.jcdecaux.com/#/opendata/vls?page=getstarted&contract=Dublin">JCDecaux</a>  | <a href="/privacy">Privacy</a>',
      },
      geojsonOptions: {
        onEachFeature: standType,
      },
      bikeTypes,
      activeBikeType: 'My Bike',
      geojsons,
    };
  },
  methods: {
    bikeTypeIsActive(bikeType) {
      return bikeType === this.activeBikeType;
    },
    changeBikeType(bikeType) {
      this.activeBikeType = bikeType;
      this.geojsons.pop();
      this.updateMap();
    },
    updateMap() {
      axios.get(`/api/v0/stand?${this.bikeTypes[this.activeBikeType].query}`).then((response) => {
        this.geojsons.push(response.data);
      });
    },
  },
  created() {
    this.updateMap();
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
