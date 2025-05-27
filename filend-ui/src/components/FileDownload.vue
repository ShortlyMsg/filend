<script setup>
import { ref, onMounted } from 'vue';
import FileIcon from '@/utils/FileIcon.vue';
import { API_ENDPOINTS } from '@/utils/api';
import { messaging, getToken, onMessage } from '@/utils/firebase';

const otp = ref('');
const files = ref([]);
const message = ref('');

const subscribeToTopic = async () => {
  try {
    const currentToken = await getToken(messaging, {
      vapidKey: 'BIpvaBgdH8WYgUeslKtuJl997WM-sC7YkUXN1avHAFnHhn4n6Uh05bcJteFhbTFR27_8iioPJYbXehHSnwm0Bro' //VAPID key
    });

    if (currentToken) {
      const response = await fetch(API_ENDPOINTS.SUBSCRIBE_TOKEN, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          token: currentToken,
          topic: otp.value
        })
      })

      if (response.ok) {
        const result = await response.json()
        console.log("Yanıt:", result)
      } else {
        alert('Abonelik başarısız')
        console.error('Hata:', await response.text())
      }
    } else {
      alert('Token alınamadı. Bildirim izni verildi mi?')
    }
  } catch (err) {
    console.error('Hata:', err)
  }
}

onMounted(() => {
  onMessage(messaging, (payload) => {
    console.log("Firebase mesajı alındı:", payload)
    if (payload?.data) {
      message.value = JSON.stringify(payload.data)

      const fileName = payload.data.fileName;
      const progressValue = parseInt(payload.data.progress, 10);
      const total = payload.data.totalMB || 0;

      const fileToUpdate = files.value.find(file => file.name === fileName);
      if (fileToUpdate) {
        fileToUpdate.progress = progressValue;
        fileToUpdate.totalMB = total;

        files.value = [...files.value];
      } else {
        files.value.push({
          name: fileName,
          hash: null,
          fileSize: 0,
          isUploaded: false,
          progress: progressValue,
          totalMB: total
        });
      }
    }
  })
})

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

    if (data.files && data.files.length > 0) {
      files.value = data.files.map(f => ({
        name: f.fileName,
        hash: f.fileHash,
        totalMB: (f.fileSize / 1024 / 1024).toFixed(2),
        uploaded: f.isUploaded,
        progress: f.isUploaded ? 100 : 0,
      }));
      subscribeToTopic(); // Burada tek tuşla çağırmış oluyoruz artık
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
                <img src="@/assets/download.svg" alt="Download"
                  @click="file.progress === 100 ? downloadFile(index) : null" :disabled="file.progress !== 100"
                  class="ml-auto cursor-pointer w-6 h-6"
                  :class="{ 'opacity-50 cursor-not-allowed': file.progress !== 100 }" />
              </div>
              <div class="w-full bg-gray-200 rounded-full h-2 mt-1">
                <div class="bg-green-400 h-2 rounded-full" :style="{ width: file.progress + '%' }"></div>
              </div>
              <div class="text-xs text-right mt-1">{{ file.totalMB }} MB</div>
            </div>
          </div>
        </div>
      </div>

      <div class="mt-4 flex justify-between space-x-3">
        <div class="">
          <input type="text" v-model="otp" class="border-2 border-dashed border-gray-400 p-2 rounded"
            placeholder="OTP kodunu buraya girin" />
        </div>
        <button @click="fetchFiles"
          class="px-4 py-1 border-2 border-green-500 text-green-600 rounded hover:bg-green-500 hover:text-white cursor-pointer">
          Listele
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped></style>
