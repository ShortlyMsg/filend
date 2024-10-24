import { createRouter, createWebHistory } from 'vue-router';
import HomePage from '@/views/HomePage.vue';
import FileUpload from '@/components/FileUpload.vue';
import FileDownload from '@/components/FileDownload.vue';

const routes = [
    {
    path: '/',
    name: 'HomePage',
    component: HomePage,
    },
    {
    path: '/upload',
    name: 'FileUpload',
    component: FileUpload,
    },
    {
    path: '/download',
    name: 'FileDownload',
    component: FileDownload,
    },
];

const router = createRouter({
    history: createWebHistory(process.env.BASE_URL),
    routes,
});

export default router;
