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
    <div class="bg-white rounded-lg shadow-lg p-6 max-w-xl w-full">
      <div class="flex items-center justify-between mb-4">
        <h2 class="text-2xl font-bold">Filend - File Send</h2>
        <span class="flex items-center border-2 border-purple-400 rounded-full p-2 transition cursor-pointer" title="Developer">
          <img @click="goToExternalSite" src="@/assets/developer.svg" alt="Kopyala" class="w-5 h-5" />
        </span>
      </div>
      <p class="text-sm text-gray-600 mb-2">Tek tıkla gönder, tek kodla al!</p>
      <div class="border-2 border-gray-300 rounded-lg text-center h-64 overflow-y-auto flex">
        
        <!-- Sol kısım: Upload simgesi ve arka plan rengi beyaz -->
        <div class="flex-1 flex justify-center items-center cursor-pointer bg-blue-50" @click="goToUpload">
          <div class="flex flex-col items-center">
            <img src="@/assets/upload.svg" alt="Upload" class="w-16 h-16" />
            <span class="text-gray-600 text-xl font-bold">Upload</span>
          </div>
        </div>

        <!-- Sağ kısım: Download simgesi ve arka plan rengi gri -->
        <div class="flex-1 flex justify-center items-center cursor-pointer bg-green-50" @click="goToDownload">
          <div class="flex flex-col items-center">
            <img src="@/assets/download.svg" alt="Download" class="w-16 h-16" />
            <span class="text-gray-600 text-xl font-bold">Download</span>
          </div>
        </div>
        
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>