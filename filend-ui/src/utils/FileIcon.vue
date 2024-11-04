<script setup>
import { computed, toRefs, defineProps } from 'vue';

const props = defineProps({
  fileName: {
    type: String,
    required: true,
  },
});

const { fileName } = toRefs(props);

const iconSrc = computed(() => {
  const extension = fileName.value.split('.').pop().toLowerCase();
  switch (extension) {
    case 'png':
    case 'jpg':
    case 'jpeg':
      return new URL('@/assets/image-icon.svg', import.meta.url).href;
    case 'pdf':
      return new URL('@/assets/pdf-icon.svg', import.meta.url).href;
    // case 'doc':
    // case 'docx':
    //   return new URL('@/assets/doc-icon.svg', import.meta.url).href;
    case 'xls':
    case 'xlsx':
      return new URL('@/assets/excel-icon.svg', import.meta.url).href;
    case 'mp4':
    case 'avi':
    case 'mkv':
    case 'mov':
    case 'wmv':
    case 'flv':
      return new URL('@/assets/video-icon.svg', import.meta.url).href;
    case 'mp3':
    case 'wav':
    case 'flac':
    case 'aac':
    case 'ogg':
    case 'wma':
      return new URL('@/assets/music-icon.svg', import.meta.url).href;
    default:
      return new URL('@/assets/file-icon.svg', import.meta.url).href;
  }
});
</script>

<template>
  <img :src="iconSrc" alt="file icon"/>
</template>