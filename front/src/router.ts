import 'vue-router';
import { createRouter, createWebHistory } from 'vue-router';
import AboutPage from './components/pages/About.vue';
import HelpPage from './components/pages/Help.vue';
import InquiryPage from './components/pages/Inquiry.vue';
import PrivacyPolicyPage from './components/pages/PrivacyPolicy.vue';
import TermsOfServicePage from './components/pages/TermsOfService.vue';
import NotFoundPage from './components/pages/error/NotFound.vue';

declare module 'vue-router' {
  interface RouteMeta {
    title: string;
    desc: string;
  }
}

const routes = [
  {
    name: 'about',
    path: '/',
    component: AboutPage,
    meta: {
      title: 'haraiai - 払い合い',
      desc: 'LINEアプリ上で2人で使える精算しない割り勘サービス ',
    },
  },
  {
    name: 'help',
    path: '/help',
    component: HelpPage,
    meta: {
      title: 'よくある質問',
      desc: 'haraiai のよくある質問とその回答',
    },
  },
  {
    name: 'inquiry',
    path: '/inquiry',
    component: InquiryPage,
    meta: {
      title: 'フィードバック',
      desc: 'haraiai に意見や要望を送る',
    },
  },
  {
    path: '/terms/privacy',
    name: 'privacy',
    component: PrivacyPolicyPage,
    meta: {
      title: 'プライバシーポリシー',
      desc: 'プライバシーポリシーについて',
    },
  },
  {
    path: '/terms/tos',
    name: 'tos',
    component: TermsOfServicePage,
    meta: {
      title: '利用規約',
      desc: '利用規約について',
    },
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: NotFoundPage,
    meta: {
      title: 'そのページは見つかりませんでした',
      desc: 'このページは存在しないか、削除されています',
    },
  },
];

const router = createRouter({
  history: createWebHistory(),
  scrollBehavior() {
    // always scroll to top
    return { top: 0 };
  },
  routes,
});

export default router;
