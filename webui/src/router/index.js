import {createRouter, createWebHistory} from 'vue-router';
import AppLayout from '@/layout/AppLayout.vue';

const router = createRouter({
    history: createWebHistory(),
    routes: [
        {
            path: '/',
            component: AppLayout,
            children: [
                {
                    path: '/',
                    name: 'dashboard',
                    component: () => import('@/views/Tickets.vue')
                },
                {
                    path: "/new",
                    name: "new",
                    component: () => import('@/views/NewTicket.vue')
                },
                {
                    path: "/edit/:id",
                    name: "edit",
                    component: () => import('@/views/EditTicket.vue')
                },
                {
                    path: '/docs',
                    name: 'docs',
                    component: () => import('@/views/pages/Docs.vue')
                },
                {
                    path: '/about',
                    name: 'about',
                    component: () => import('@/views/pages/About.vue')
                }
            ]
        },
        {
            path: '/:pathMatch(.*)*',
            name: 'NotFound',
            component: import('@/views/pages/NotFound.vue'),
        },
    ]
});

export default router;
