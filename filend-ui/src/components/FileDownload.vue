<script setup>
import { ref } from 'vue';
import FileIcon from '@/utils/FileIcon.vue';
import { API_ENDPOINTS } from '@/utils/api';

const otp = ref('');
const files = ref([]);

const fetchFiles = async () => {
  if (!otp.value) {
      alert('Lütfen OTP kodunu girin.');
      return;
  }

  try {
    const response = await fetch(API_ENDPOINTS.GET_ALL_FILES, {
      method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ otp: otp.value })
    });
    
    const data = await response.json();

      if (data.files && data.files.length > 0 && data.hashes && data.hashes.length) {
          files.value = data.files.map((fileName, index) => ({
              name: fileName,
              hash: data.hashes[index]
          }));
      } else {
          files.value = [];
          alert('Dosya bulunamadı.');
      }
    } catch (error) {
      console.error('Hata:', error);
      alert('Bir hata oluştu: ' + error.message);
    }
};

const downloadFile = async (index) => {
  const file = files.value[index];
  if (!file) return;

    try {
      const response = await fetch(API_ENDPOINTS.DOWNLOAD_FILES, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
          },
        body: JSON.stringify({ otp: otp.value, fileHash: file.hash })
      });

    const blob = await response.blob();
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = file.name;
    a.click();
    } catch (error) {
      console.error('Hata:', error);
      alert('Bir hata oluştu: ' + error.message);
    }
};
</script>

<template>
  <div class="flex justify-center items-center h-screen">
    <div class="bg-white rounded-lg shadow-lg p-6 max-w-xl w-full">
      <h2 class="text-2xl font-bold mb-4">Filend - File Send</h2>
      <p class="text-sm text-gray-600 mb-2">Tek tıkla gönder, tek kodla al!</p>
      <div class="border-2 border-gray-300 rounded-lg p-6 text-center h-64 overflow-y-auto relative">
        <span v-if="files.length === 0" class="absolute inset-0 flex justify-center items-center text-green-500 text-2xl">Lütfen Tek Seferlik şifreyi aşağıya giriniz</span>
        <img v-if="files.length === 0" src="@/assets/download.svg" alt="download" class="w-12 h-12  absolute bottom-4 left-4"/>
        <div v-for="(file, index) in files" :key="index" class="mb-4">
          <div class="flex items-center">
            <FileIcon :fileName="file.name" class="32px"/>
            <div class="flex flex-col ml-4 w-full">
              <div class="flex items-center">
                <span class="text-sm">{{ file.name }}</span>
                <img src="@/assets/download.svg" alt="Download" @click="downloadFile(index)" class="ml-auto cursor-pointer w-6 h-6" />
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <div class="mt-4 flex justify-between space-x-3">
        <div class="">
        <input 
          type="text" 
          v-model="otp" 
          class="border border-green-400 p-2 rounded"
          placeholder="OTP kodunu girin"
        />
      </div>
        <button @click="fetchFiles" class="px-4 py-0 bg-green-400 text-white rounded hover:bg-green-700">
        Listele
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>
