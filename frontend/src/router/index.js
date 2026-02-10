import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import DashboardView from '../views/DashboardView.vue'
import CalendarView from '../views/CalendarView.vue' // Yeni sayfamız

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        { path: '/', redirect: '/login' },
        { path: '/login', component: LoginView },
        { path: '/dashboard', component: DashboardView },
        { path: '/calendar', component: CalendarView } // Takvim rotası
    ]
})

export default router