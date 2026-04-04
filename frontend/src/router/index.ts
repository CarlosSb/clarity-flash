import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/HomeView.vue'),
    meta: { title: 'AulaFlash — Início' },
  },
  {
    path: '/session/:id',
    name: 'SessionDetail',
    component: () => import('@/views/SessionDetailView.vue'),
    props: true,
    meta: { title: 'Detalhes da Sessão' },
  },
  {
    path: '/quiz/:id',
    name: 'Quiz',
    component: () => import('@/views/QuizView.vue'),
    props: true,
    meta: { title: 'Quiz' },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  },
})

router.beforeEach((to) => {
  if (to.meta.title) {
    document.title = to.meta.title as string
  }
})

export default router
