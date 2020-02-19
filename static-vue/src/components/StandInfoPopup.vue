<template>
  <div>
    <table>
      <tr><td><b>Type of Stands</b></td><td>{{ stand.type }}</td></tr>
      <tr><td><b>Location</b></td><td>{{ stand.name }}</td></tr>
      <tr><td><b>Number of Stands</b></td><td>{{ numberOfStands() }}</td></tr>
      <tr v-if="showId"><td><b>ID</b></td><td>{{ stand.id }}</td></tr>
      <tr v-if="stand.verified">
        <td colspan=2><small>This is a user submitted stand we have not verified yet.</small></td>
      </tr>
    </table>
    <router-link class="btn btn-primary" :to="{ name: 'edit', params: { id: stand.id }}">
      <b-icon-pencil/> Edit
    </router-link>
  </div>
</template>

<script>
export default {
  name: 'StandInfoPopup',
  props: {
    stand: {
      type: Object,
    },
  },
  data() {
    return {
      showId: false,
    };
  },
  methods: {
    numberOfStands() {
      return this.stand.numberOfStands > 0 ? this.stand.numberOfStands : 'Unknown';
    },
  },
};
</script>

<style scoped>

</style>
