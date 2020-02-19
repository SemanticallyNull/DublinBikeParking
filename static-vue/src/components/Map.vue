<template>
  <l-map class='map' :zoom='zoom' :center='center' :options='mapOptions'
         @contextmenu="addStand">
    <l-tile-layer :url='url' :options="tileLayerOptions">
    </l-tile-layer>
    <l-marker-cluster ref="clusterRef" :options='clusterOptions'>
      <l-geo-json
        v-for='geojson in geojsons'
        :key='geojson.features[0].properties.id'
        :geojson='geojson'
        :options='geojsonOptions'>
      </l-geo-json>
    </l-marker-cluster>
    <l-marker ref="addMarker" v-if="addMarkerShow" @popupclose="hideAddMarker"
              :lat-lng="addMarkerPosition">
      <l-popup ref="addMarkerPopup" >FOO</l-popup>
    </l-marker>
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
import Vue from 'vue';
import {
  LMap, LTileLayer, LGeoJson, LControl, LMarker, LPopup,
} from 'vue2-leaflet';
import LMarkerCluster from 'vue2-leaflet-markercluster';
import axios from 'axios';
import StandInfoPopup from './StandInfoPopup.vue';
import standIcons from '../lib/stand-icons';

const geojsons = [];

const bikeTypes = {
  'My Bike': {
    query: 'dublinbikes=off',
    standData: null,
  },
  DublinBikes: {
    query: 'dublinbikes=only',
    standData: null,
  },
};

function standType(f, l) {
  const StandInfo = Vue.extend(StandInfoPopup);
  l.setIcon(standIcons[f.properties.type]).bindPopup(new StandInfo({
    propsData: {
      stand: f.properties,
    },
  }).$mount().$el);
}

export default {
  name: 'Map',
  components: {
    'l-map': LMap,
    'l-tile-layer': LTileLayer,
    'l-geo-json': LGeoJson,
    'l-marker': LMarker,
    'l-marker-cluster': LMarkerCluster,
    'l-control': LControl,
    'l-popup': LPopup,
  },
  data() {
    return {
      url: 'https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png',
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
      addMarkerShow: false,
      addMarkerPosition: { lat: 0, lng: 0 },
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
      if (this.bikeTypes[this.activeBikeType].standData != null) {
        this.geojsons.push(this.bikeTypes[this.activeBikeType].standData);
        return;
      }
      axios.get(`/api/v0/stand?${this.bikeTypes[this.activeBikeType].query}`).then((response) => {
        this.geojsons.push(response.data);
        this.bikeTypes[this.activeBikeType].standData = response.data;
      });
    },
    addStand(e) {
      this.addMarkerPosition = e.latlng;
      this.addMarkerShow = true;
      this.$nextTick(() => this.$refs.addMarker.mapObject.openPopup());
    },
    hideAddMarker() {
      this.addMarkerShow = false;
    },
  },
  created() {
    this.updateMap();
  },
};
</script>

<style scoped>
@import '~leaflet.markercluster/dist/MarkerCluster.css';
@import '~leaflet.markercluster/dist/MarkerCluster.Default.css';

.map {
  width: 100%;
  height: 100%;
}
</style>
