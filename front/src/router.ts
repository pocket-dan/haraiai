import { createRouter, createWebHistory } from 'vue-router'
import AboutPage from "./components/pages/About.vue"
import HelpPage from "./components/pages/Help.vue"
import InquiryPage from "./components/pages/Inquiry.vue"
import NotFoundPage from "./components/pages/error/NotFound.vue"

const routes = [
  { path: '/', name: 'about', component: AboutPage },
  { path: '/help', name: 'help', component: HelpPage },
  { path: '/inquiry', name: 'inquiry', component: InquiryPage },
  { path: '/:pathMatch(.*)*', name: 'not-found', component: NotFoundPage },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
