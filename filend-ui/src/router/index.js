import { createRouter, createWebHistory } from 'vue-router';
import FileUpload from '@/components/FileUpload.vue';
import FileDownload from '@/components/FileDownload.vue';

const routes = [
    { path: '/upload', component: FileUpload },
    { path: '/download', component: FileDownload },
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

export default router;
