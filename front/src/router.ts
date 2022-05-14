import { createRouter, createWebHistory } from 'vue-router'
import AboutPage from "./pages/About.vue"
import HelpPage from "./pages/Help.vue"
import InquiryPage from "./pages/Inquiry.vue"
import NotFoundPage from "./pages/error/NotFound.vue"

const routes = [
  { path: '/about', name: 'about', component: AboutPage },
  { path: '/help', name: 'help', component: HelpPage },
  { path: '/inquiry', name: 'inquiry', component: InquiryPage },
  { path: '/:pathMatch(.*)*', name: 'not-found', component: NotFoundPage },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
