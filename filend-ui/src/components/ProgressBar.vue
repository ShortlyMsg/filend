<script setup>
import { defineProps } from 'vue';
const props = defineProps({
  progress: {
    type: Number,
    required: true,
    default: 0,
  },
});
</script>

<template>
  <div class="w-full bg-gray-200 rounded-full h-2">
    <div
      class="bg-blue-600 h-2 rounded-full"
      :style="{ width: `${props.progress}%` }"
    ></div>
  </div>
                <!-- <ProgressBar :progress="progress"/>
                <span class="text-sm text-left mt-1" style="display: block; text-align: left;">
                  {{ file.uploadProgress ? `${file.uploadProgress} MB / ${file.totalSize} MB` : `0.00 MB / ${file.totalSize} MB` }}
                </span> -->
</template>