<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()

const goToUpload = () => {
  router.push({ name: 'FileUpload' })
}

const goToDownload = () => {
  router.push({ name: 'FileDownload' })
}
const goToExternalSite = () => {
  window.open("https://github.com/shortlymsg", "_blank");
}
</script>

<template>
  <div class="flex justify-center items-center h-screen">
    <div class="grid grid-rows-4 items-center">
      <button @click="goToExternalSite" class="hexagon bg-purple-500 text-white font-bold text-2xl focus:outline-none 
      hover:bg-purple-700 hover:scale-90 active:scale-75 transition-all duration-200 mx-auto -mb-22">
        Developer
      </button>
      <div class="flex justify-center space-x-1">
        <button @click="goToUpload" class="hexagon text-white font-bold text-2xl focus:outline-none
          hover:scale-90 active:scale-75 transition-all duration-200"
          style="background: radial-gradient(circle, #3B82F6 0%, #000000 75%);">
          Upload
        </button>
        <button @click="goToDownload" class="hexagon text-white font-bold text-2xl focus:outline-none 
          hover:scale-90 active:scale-75 transition-all duration-200"
          style="background: radial-gradient(circle, #70FF88 0%, #000000 75%);">
          Download
        </button>
      </div>
      <!-- <button @click="goToExternalSite" 
        class="relative w-32 h-32 bg-red-500 text-white font-bold text-lg flex items-center justify-center 
        hover:bg-red-700 hover:scale-90 active:scale-75 transition-all duration-200 mx-auto -mt-14" 
        style="clip-path: polygon(50% 0%, 100% 25%, 100% 75%, 50% 100%, 0% 75%, 0% 25%);">
        MSG
      </button> -->
    </div>
  </div>
</template>

<style scoped>
.hexagon {
  position: relative;
  width: 12rem;
  height: 12rem;
  clip-path: polygon(50% 0%, 100% 25%, 100% 75%, 50% 100%, 0% 75%, 0% 25%);
}
</style>