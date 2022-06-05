import { createRouter, createWebHistory } from 'vue-router'
import AboutPage from "./components/pages/About.vue"
import HelpPage from "./components/pages/Help.vue"
import InquiryPage from "./components/pages/Inquiry.vue"
import PrivacyPolicyPage from "./components/pages/PrivacyPolicy.vue"
import TermsOfServicePage from "./components/pages/TermsOfService.vue"
import NotFoundPage from "./components/pages/error/NotFound.vue"

const routes = [
  { path: '/', name: 'about', component: AboutPage },
  { path: '/help', name: 'help', component: HelpPage },
  { path: '/inquiry', name: 'inquiry', component: InquiryPage },
  { path: '/terms/privacy', name: 'privacy', component: PrivacyPolicyPage },
  { path: '/terms/tos', name: 'tos', component: TermsOfServicePage },
  { path: '/:pathMatch(.*)*', name: 'not-found', component: NotFoundPage },
]

const router = createRouter({
  history: createWebHistory(),
  scrollBehavior(to, from, savedPosition) {
    // always scroll to top
    return { top: 0 }
  },
  routes,
})

export default router
