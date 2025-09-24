import { createWebHistory, createRouter, type RouteRecordRaw } from 'vue-router'
import HomeView from './pages/HomeView.vue'
import SettingsView from './pages/SettingsView.vue'
import AudiobooksView from './pages/HistoryView.vue'

const routes: RouteRecordRaw[] = [
    {
        path: '/',
        component: HomeView,
        name: 'Create an audiobook',
        meta: {
            id: 'page-home'
        }
    },
    {
        path: '/settings',
        component: SettingsView,
        name: 'Audiobook Settings',
        meta: {
            id: 'page-settings'
        }
    },
    {
        path: '/history',
        component: AudiobooksView,
        name: 'Previously created audiobooks',
        meta: {
            id: 'page-history'
        }
    }
]

export const router = createRouter({
    history: createWebHistory(),
    routes
})