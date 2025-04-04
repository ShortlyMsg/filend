<script setup>
import { ref, onMounted } from 'vue';
import FileIcon from '@/utils/FileIcon.vue';
import { API_ENDPOINTS } from '@/utils/api';
import { messaging, onMessage } from '@/utils/firebase';

const otp = ref('74cmze');
const files = ref([]);
const progress = ref(0);

const listenToFirebaseMessages = () => {
  onMessage(messaging, (payload) => {
    console.log("Tüm mesajlar (filtrelenmeden):", payload);
    // ASYNC İŞLEM YAPMAYIN! Direkt işleyin.
    if (payload.topic === otp.value) {
      try {
        const data = JSON.parse(payload.notification?.body || "{}");
        console.log("Alınan veri:", {
          fileName: payload.notification?.title,
          ...data
        });
      } catch (error) {
        console.error("JSON parse hatası:", error);
      }
    }
  });
};

onMounted(() => {
  listenToFirebaseMessages();
});

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
        <div v-if="files.length === 0" class="absolute inset-0 flex flex-col items-center justify-center">
          <img src="@/assets/download-file.png" alt="download" class="w-40 h-40" />
          <span class="text-green-500 text-lg font-extrabold mt-2">Lütfen Tek Seferlik şifreyi aşağıya giriniz</span>
        </div>
        <div v-for="(file, index) in files" :key="index" class="mb-4">
          <div class="flex items-center">
            <FileIcon :fileName="file.name" class="32px" />
            <div class="flex flex-col ml-4 w-full">
              <div class="flex items-center">
                <span class="text-sm">{{ file.name }}</span>
                <img src="@/assets/download.svg" alt="Download" @click="downloadFile(index)"
                  class="ml-auto cursor-pointer w-6 h-6" />
              </div>
              <div class="w-full bg-gray-200 rounded-full h-2 mt-1">
                <div class="bg-green-400 h-2 rounded-full" :style="{ width: progress + '%' }"></div>
              </div>
              <div class="text-xs text-right mt-1">{{ progress }} MB</div>
            </div>
          </div>
        </div>
      </div>

      <div class="mt-4 flex justify-between space-x-3">
        <div class="">
          <input type="text" v-model="otp" class="border-2 border-dashed border-gray-400 p-2 rounded"
            placeholder="OTP kodunu buraya girin" />
        </div>
        <button @click="fetchFiles" class="px-4 py-1 border-2 border-green-500 text-green-600 rounded hover:bg-green-500 hover:text-white cursor-pointer">
          Listele
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped></style>
