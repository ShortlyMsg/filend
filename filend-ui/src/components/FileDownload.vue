<script setup>
import { ref } from 'vue';
import FileIcon from '@/utils/FileIcon.vue';

const otp = ref('');
const files = ref([]);

const fetchFiles = async () => {
  if (!otp.value) {
      alert('Lütfen OTP kodunu girin.');
      return;
  }

  try {
    const response = await fetch('http://localhost:9091/getAllFiles', {
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
      const response = await fetch('http://localhost:9091/download', {
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
  <div class="flex justify-center items-center h-screen bg-gray-100">
    <div class="bg-white rounded-lg shadow-lg p-6 max-w-xl w-full">
      <div class="mb-4">
        <input 
          type="text" 
          v-model="otp" 
          class="border border-gray-300 p-2 rounded"
          placeholder="OTP kodunu girin"
        />
      </div>
      <div class="border-2 border-gray-300 rounded-lg p-6 text-center h-64 overflow-y-auto">
        <ul>
          <li v-for="(file, index) in files" :key="index" class="mb-4">
            <div class="flex items-center">
              <FileIcon :fileName="file.name" class="32px"/>
              <div class="flex flex-col ml-4 w-full">
                <div class="flex items-center">
                  <span class="text-sm">{{ file.name }}</span>
                  <img src="@/assets/download.svg" alt="Download" @click="downloadFile(index)" class="ml-auto cursor-pointer w-6 h-6 text-[#5D3FD3]" />
                </div>
              </div>
            </div>
          </li>
        </ul>
      </div>
      <div class="mt-4 flex justify-end space-x-3">
        <button @click="fetchFiles" class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700">
        Listele
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>
