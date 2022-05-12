import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

import PrimeVue from 'primevue/config'
import Button from 'primevue/button'

const app = createApp(App)

app.use(router)
app.use(PrimeVue)

app.mount('#app')

app.component('Button', Button)
