<template>
  <div class="container">
    <div>
      <b-navbar toggleable="lg">
        <b-breadcrumb>
          <b-breadcrumb-item :to='{ name: "home" }'>
            Home
          </b-breadcrumb-item>
          <b-breadcrumb-item active>Edit {{id}}</b-breadcrumb-item>
        </b-breadcrumb>
      </b-navbar>
    </div>

    <div class="form">
      <h1>Edit {{id}}</h1>

      <b-form @submit="onSubmit" @reset="onReset">
        <b-form-group
          id="input-group-checked"
        >
          <b-form-checkbox
            id="input-checked"
            v-model="form.checked"
          >Checked</b-form-checkbox>
        </b-form-group>

        <b-form-group
          id="input-group-verified"
        >
          <b-form-checkbox
            id="input-verified"
            v-model="form.verified"
          >Verified</b-form-checkbox>
        </b-form-group>

        <b-form-group
          id="input-group-name"
          label="Name"
          label-for="input-name"
        >
          <b-form-input
            id="input-name"
            v-model="form.name"
            type="text"
          ></b-form-input>
        </b-form-group>

        <b-form-group
          id="input-group-numberofstands"
          label="Number of Stands"
          label-for="input-numberofstands"
        >
          <b-form-input
            id="input-numberofstands"
            v-model="form.numberOfStands"
            type="number"
          ></b-form-input>
        </b-form-group>

        <b-form-group
          id="input-group-notes"
          label="Notes"
          label-for="input-notes"
        >
          <b-form-textarea
            id="input-notes"
            v-model="form.notes"
          ></b-form-textarea>
        </b-form-group>

        <b-form-group
          id="input-group-source"
          label="Source"
          label-for="input-source"
        >
          <b-form-input
            id="input-source"
            v-model="form.source"
            type="text"
          ></b-form-input>
        </b-form-group>

        <b-form-group
          id="input-group-type"
          label="Type"
          label-for="input-type"
        >
          <b-form-select
            id="input-type"
            v-model="form.type"
            :options="standTypes"
          ></b-form-select>
        </b-form-group>

        <div>
          <h4>Thefts</h4>
          <ul v-if="form.thefts.length !== 0">
            <li v-for="theft in form.thefts" :key="theft.ID">{{theft.ID}}</li>
          </ul>
          <p v-else>There have been no thefts reported at this stand</p>
        </div>

        <b-button type="submit" variant="primary">Submit</b-button>
        <b-button type="reset" variant="danger">Reset</b-button>
      </b-form>
    </div>
  </div>
</template>

<script>
import axios from 'axios';
import standIcons from '../lib/stand-icons';


export default {
  name: 'edit',
  data() {
    return {
      id: this.$route.params.id,
      standTypes: Object.keys(standIcons),
      form: {
        id: '123',
      },
    };
  },
  methods: {
    onSubmit: () => {},
    onReset: () => {},
  },
  created() {
    axios.get(`/api/v0/stand/${this.id}`).then((response) => {
      this.form = response.data;
    });
  },
};
</script>

<style scoped>
  .form {
    margin-top: 1rem;
  }
</style>
