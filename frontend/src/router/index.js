import { createRouter, createWebHashHistory } from 'vue-router';
import LoginView from '../views/LoginView.vue';
import DashboardView from '../views/DashboardView.vue';
import CalendarView from '../views/CalendarView.vue';
import SettingsView from '../views/SettingsView.vue';
import AssetDetailView from '../views/AssetDetailView.vue'; // Yeni dosya

const router = createRouter({
    history: createWebHashHistory(),
    routes: [
        { path: '/', redirect: '/dashboard' },
        { path: '/login', name: 'login', component: LoginView },
        {
            path: '/dashboard',
            name: 'dashboard',
            component: DashboardView,
            meta: { requiresAuth: true }
        },
        {
            path: '/dashboard/details/:type',
            name: 'asset-detail',
            component: AssetDetailView,
            meta: { requiresAuth: true },
            props: true
        },
        {
            path: '/calendar',
            name: 'calendar',
            component: CalendarView,
            meta: { requiresAuth: true }
        },
        {
            path: '/settings',
            name: 'settings',
            component: SettingsView,
            meta: { requiresAuth: true }
        }
    ]
});

router.beforeEach((to, from, next) => {
    const isAuthenticated = localStorage.getItem('token');
    if (to.meta.requiresAuth && !isAuthenticated) {
        next('/login');
    } else if (to.path === '/login' && isAuthenticated) {
        next('/dashboard');
    } else {
        next();
    }
});

export default router;