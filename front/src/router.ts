import { createRouter, createWebHistory } from 'vue-router'
import AboutPage from "./pages/About.vue"
import HelpPage from "./pages/Help.vue"
import InquiryPage from "./pages/Inquiry.vue"

const routes = [
  { path: '/about', name: 'about', component: AboutPage },
  { path: '/help', name: 'help', component: HelpPage },
  { path: '/inquiry', name: 'inquiry', component: InquiryPage },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
