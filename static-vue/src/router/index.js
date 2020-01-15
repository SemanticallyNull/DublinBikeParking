import Vue from 'vue';
import VueRouter from 'vue-router';
import Home from '../views/Home.vue';
import Privacy from '../views/Privacy.vue';

Vue.use(VueRouter);

const routes = [
  {
    path: '/',
    name: 'home',
    component: Home,
  },
  {
    path: '/privacy',
    name: 'privacy',
    component: Privacy,
  },
];

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes,
});

export default router;
