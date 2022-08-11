import { createApp } from 'vue';
import App from './App.vue';
import router from './router';
import VueGtag from 'vue-gtag';

const app = createApp(App);

app.use(router);

if (import.meta.env.PROD) {
  app.use(
    VueGtag,
    {
      config: { id: 'G-FMS9YJCVM8' },
    },
    router
  );
}

app.mount('#app');
