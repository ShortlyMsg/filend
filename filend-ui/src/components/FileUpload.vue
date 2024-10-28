<script setup>
import { ref } from 'vue';

const currentStep = ref(1); // 1: Dosya Seçimi, 2: Önizleme, 3: OTP Gösterimi
const selectedFiles = ref([]);
const otpMessage = ref("");
const fileListMessage = ref("");
const copyNotification = ref('');
const copied = ref(false);

function copyToClipboard(text) {
  navigator.clipboard.writeText(text).then(() => {
    copyNotification.value = '✓';
    copied.value = true;
    setTimeout(() => {
      copyNotification.value = '';
      copied.value = false;
    }, 2000);
  });
}

function goBackStep() {
  if (currentStep.value > 1) {
    currentStep.value--;
  }
}

function onFileChange(event) {
  const newFiles = Array.from(event.target.files);
  selectedFiles.value = [...selectedFiles.value, ...newFiles];
}

// const removeFile = (index: number) => {
//   files.value.splice(index, 1)
// }

async function calculateSHA256(file) {
  const arrayBuffer = await file.arrayBuffer();
  const hashBuffer = await crypto.subtle.digest("SHA-256", arrayBuffer);
  const hashArray = Array.from(new Uint8Array(hashBuffer));
  const hashHex = hashArray.map(b => b.toString(16).padStart(2, "0")).join("");
  return hashHex;
}

function handleFileUpload(event) {
  const files = event.target.files;
  selectedFiles.value = Array.from(files);
  currentStep.value = 2; // Dosya önizleme adımına geç
}

function handleDrop(event) {
  const files = event.dataTransfer.files;
  handleFileUpload({ target: { files } });
}

async function uploadFiles() {
  // Dosyaları yüklemeye başla
  const fileHashes = [];
  for (const file of selectedFiles.value) {
    const fileHash = await calculateSHA256(file);
    fileHashes.push(fileHash);
  }

  // Hash kontrol isteği
  const hashCheckResponse = await fetch("http://localhost:9091/checkFileHash", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ fileHashes }),
  });
  const hashCheckData = await hashCheckResponse.json();

  // Dosya yükleme işlemi için formData oluştur
  const formData = new FormData();
  selectedFiles.value.forEach((file, i) => {
    const fileHash = fileHashes[i];
    if (hashCheckData.fileStatus[fileHash]) {
      formData.append("files", file);
      formData.append("fileHashes", fileHash);
    } else {
      formData.append("fileNames[]", file.name);
      formData.append("fileHashes[]", fileHash);
    }
  });

  // Dosya yükleme isteği
  fetch("http://localhost:9091/upload", {
    method: "POST",
    body: formData,
  })
    .then(response => response.json())
    .then(data => {
      if (data.otp) {
        otpMessage.value = data.otp;
        fileListMessage.value = "Yüklenen Dosyalar: " + data.fileNames.join(", ");
        currentStep.value = 3; // OTP gösterim adımına geç
      } else {
        otpMessage.value = "Hata: OTP alınamadı.";
      }
    })
    .catch(error => {
      console.error("Hata:", error);
      otpMessage.value = "Sunucu ile iletişimde bir sorun var.";
    });
}
</script>
<!-- <script setup>
//import { computed } from 'vue'


//const iconSrc = computed(() => {
  //if (.png)
//})
</script> -->

<template>
  <div class="flex justify-center items-center h-screen bg-gray-100">
    <div class="bg-white rounded-lg shadow-lg p-6 fixed-size-1080p">
      <div v-if="currentStep === 1">
        <!-- CS 1 Upload -->
        <h2 class="text-2xl font-bold mb-4">Filend - File Send</h2>
        <p class="text-sm text-gray-600 mb-2">Tek tıkla gönder, tek kodla al!</p>
        <div class="border-2 border-dashed border-gray-300 rounded-lg p-14 text-center fixed-size-border2"
          @dragenter.prevent="dragEnter"
          @dragover.prevent="dragOver"
          @dragleave.prevent="dragLeave"
          @drop.prevent="handleDrop"
        >
          <p class="text-gray-600">Drag and drop a files here</p>
          <p class="text-gray-600 my-2">or</p>
          <label
            for="file-upload"
            class="cursor-pointer text-blue-600 border border-blue-600 rounded-md px-4 py-2 inline-block hover:bg-blue-600 hover:text-white transition"
          >
            Dosya Seç
          </label>
          <input id="file-upload" type="file" class="hidden" @change="handleFileUpload" multiple />
        </div>
        <div class="mt-4 text-sm text-gray-600 flex justify-between">
          <p>Accepted file types: All Types</p>
          <p>Max files: 20 | Max file size: 2GB</p>
        </div>
      </div>

      <div v-else-if="currentStep === 2">
        <!-- CS 2 Önizleme -->
        <h2 class="text-2xl font-bold mb-4">Yüklenecek Dosyalar</h2>
        <div class="border-2 border-gray-300 rounded-lg p-14 text-center h-64 overflow-y-auto fixed-size-border2">
        <ul>
          <li v-for="(file, index) in selectedFiles" :key="index" class="flex justify-between mb-2">
            <img src="@/assets/file.svg" alt="file icon" class="w-4 h-4" />
            <span>{{ file.name }}</span>
            <span>{{ (file.size / (1024 * 1024)).toFixed(2) }} MB</span>
          </li>
        </ul>
        </div>
        <div class="mt-4 flex justify-end space-x-3">
          <input type="file" 
                multiple
                class="block text-sm text-gray-500
                        file:mr-4 file:py-2 file:px-4
                        file:rounded-full file:border-0
                        file:text-sm file:font-semibold
                        file:bg-blue-50 file:text-blue-700
                        hover:file:bg-blue-100"
                @change="onFileChange"
          />
          <button @click="goBackStep" class="px-4 py-2 text-gray-600 hover:text-gray-800">
            Cancel
          </button>
          <button @click="uploadFiles" class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700">
          Gönder
          </button>
        </div>
      </div>

      <div v-else-if="currentStep === 3">
        <!-- CS 3 OTP-->
        <h2 class="text-2xl font-bold mb-4">Dosyalar Yüklendi!</h2>
        <div class="border-2 border-gray-300 rounded-lg p-14 text-center relative fixed-size-border2 flex flex-col justify-between h-full">
          <div>
            <p id="otpMessage" class="text-3xl text-green-600 font-semibold">{{ otpMessage }}</p>
            <button @click="copyToClipboard(otpMessage)" class="absolute top-4 right-4 flex items-center text-blue-500 hover:text-blue-700 transition">
              <img v-if="!copied" src="@/assets/copy.svg" alt="Kopyala" class="w-5 h-5 mr-1" />
              <p v-else class="text-m text-green-500 mt-2">✓</p>
            </button>
          </div>
          <p class="mt-4">Yukardaki kodu alıcıya gönderiniz.</p>
        </div>
        <p id="fileList" class="mt-4 text-gray-600">{{ fileListMessage }}</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.fixed-size-1080p {
  width: 492px;
  height: 398px;
}
.fixed-size-border2 {
  width: 444px;
  height: 250px;
}
</style>