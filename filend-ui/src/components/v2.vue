<script setup>
import { ref } from 'vue';

const currentStep = ref(1); // 1: Dosya Seçimi, 2: Önizleme, 3: OTP Gösterimi
const selectedFiles = ref([]);
const otpMessage = ref("");
//const fileListMessage = ref("");
const copyNotification = ref('');
const copied = ref(false);
const allFilesUploaded = ref(false);

function copyToClipboard(text) {
  navigator.clipboard.writeText(text).then(() => {
    copyNotification.value = '';
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
  if (selectedFiles.value.length > 0) {
    currentStep.value = 2; // İkinci aşamaya geç
  }
}

function removeFile(index) {
  selectedFiles.value.splice(index, 1);
}

async function calculateSHA256(file) {
  const arrayBuffer = await file.arrayBuffer();
  const hashBuffer = await crypto.subtle.digest("SHA-256", arrayBuffer);
  const hashArray = Array.from(new Uint8Array(hashBuffer));
  const hashHex = hashArray.map(b => b.toString(16).padStart(2, "0")).join("");
  return hashHex;
}

// function handleFileUpload(event) {
//   const files = event.target.files;
//   selectedFiles.value = Array.from(files);
//   currentStep.value = 2; // Dosya önizleme adımına geç
// }

// function handleDrop(event) {
//   const files = event.dataTransfer.files;
//   handleFileUpload({ target: { files } });
// }

async function uploadFiles() {
  allFilesUploaded.value = false;
  
  const fileHashes = [];
  for (const file of selectedFiles.value) {
    const fileHash = await calculateSHA256(file);
    fileHashes.push(fileHash);

    file.totalSize = (file.size / (1024 * 1024)).toFixed(2);
    file.uploadProgress = 0;
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

  selectedFiles.value.forEach((file, i) => {
    const formData = new FormData();
    const fileHash = fileHashes[i];
    
    if (hashCheckData.fileStatus[fileHash]) {
      formData.append("files", file);
      formData.append("fileHashes", fileHash);
    } else {
      formData.append("fileNames[]", file.name);
      formData.append("fileHashes[]", fileHash);
    }


  const xhr = new XMLHttpRequest();
  xhr.upload.onprogress = (event) => {
    if (event.lengthComputable) {
      const progressValue = (event.loaded / event.total) * 100;
      progress.value = progressValue;

      const totalSizeMB = (event.total / (1024 * 1024)).toFixed(2);
      const loadedSizeMB = (event.loaded / (1024 * 1024)).toFixed(2);
        
      file.uploadProgress = loadedSizeMB;
      file.totalSize = totalSizeMB;
      }
    };

  xhr.onload = () => {
    const data = JSON.parse(xhr.responseText);
    if (data.success) {
      allFilesUploaded.value = true;
      otpMessage.value = data.otp;
      currentStep.value = 3; 
    } else {
      otpMessage.value = "Hata: OTP alınamadı.";
    }
  };

  xhr.onerror = () => {
    console.error("Hata: Sunucu ile iletişimde bir sorun var.");
    otpMessage.value = "Sunucu ile iletişimde bir sorun var.";
  };

  xhr.open("POST", "http://localhost:9091/upload");
  xhr.send(formData);
});
}

watch(selectedFiles, () => {
  allFilesUploaded.value = selectedFiles.value.length === 0; // Tüm dosyalar yüklendiğinde güncelle
});
</script>
<!-- <script setup>
//import { computed } from 'vue'


//const iconSrc = computed(() => {
  //if (.png)
//})
</script> -->

<template>
  <div class="flex justify-center items-center h-screen bg-gray-100">
    <div class="bg-white rounded-lg shadow-lg p-6 max-w-xl w-full">
      <h2 class="text-2xl font-bold mb-4">Filend - File Send</h2>
      <p class="text-sm text-gray-600 mb-2">Tek tıkla gönder, tek kodla al!</p>
      <div v-if="currentStep === 1">
        <!-- CS 1 Upload -->
        <div class="border-2 border-dashed border-gray-300 rounded-lg p-24 text-center"
          @dragenter.prevent="dragEnter"
          @dragover.prevent="dragOver"
          @dragleave.prevent="dragLeave"
          @drop.prevent="handleDrop"
        >
          <p class="text-gray-600">Dosyayı Buraya Sürükle Bırak</p>
          <p class="text-gray-600">Yada</p>
          <label
            for="file-upload"
            class="cursor-pointer text-blue-600 border border-blue-600 rounded-md px-4 py-2 inline-block hover:bg-blue-600 hover:text-white transition"
          >
            Dosya Seç
          </label>
          <input id="file-upload" type="file" class="hidden" @change="onFileChange" multiple />
        </div>
        <div class="mt-4 text-xs text-gray-600 flex justify-between">
          <p>Accepted file types: All Types</p>
          <p>Max files: 20 | Max file size: 2GB</p>
        </div>
      </div>

      <div v-else-if="currentStep === 2">
        <!-- CS 2 Önizleme -->
        <div class="border-2 border-gray-300 rounded-lg p-6 text-center h-64 overflow-y-auto">
          <ul>
            <li v-for="(file, index) in selectedFiles" :key="index" class="mb-4">
              <div class="flex items-center">
                <img src="@/assets/file.svg" alt="file icon" class="32px" />
                <span class="text-sm ml-4">{{ file.name }}</span>
                <div class="ml-auto flex items-center">
                  <button @click="removeFile(index)" class="text-red-500 hover:text-red-700 ml-2">✕</button>
                </div>
              </div>
              <div class="mt-1">
                <div class="flex flex-col">
                  <div class="w-full bg-gray-200 rounded-full h-2 mr-2">
                    <div :style="{ width: `${(file.uploadProgress || 0)}%` }" class="bg-blue-600 h-2 rounded-full"></div>
                  </div>
                  <span class="text-sm text-left mt-1">
                    {{ file.uploadProgress ? `${(file.uploadProgress || 0).toFixed(2)} MB / ${(file.size / (1024 * 1024)).toFixed(2)} MB` : `0.00 MB / ${(file.size / (1024 * 1024)).toFixed(2)} MB` }}
                  </span>
                </div>
              </div>
            </li>
          </ul>
        </div>
        <div class="mt-4 text-xs text-gray-600 flex justify-between">
          <p>Accepted file types: All Types</p>
          <p>Max files: 20 | Max file size: 2GB</p>
        </div> 
        <div class="mt-4 flex justify-end space-x-3">
          <label class="px-4 py-2 bg-purple-600 text-white rounded hover:bg-purple-700 cursor-pointer">
            <input type="file" multiple hidden @change="onFileChange"/>
          Upload More Files
          </label>
          <button @click="goBackStep" class="px-4 py-2 text-gray-600 hover:text-gray-800">
            Geri
          </button>
          <button @click="uploadFiles" class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700">
            Gönder
          </button>
        </div>
      </div>

      <div v-else-if="currentStep === 3">
        <!-- CS 3 OTP-->
        <div class="border-2 border-gray-300 rounded-lg p-24 text-center relative flex flex-col justify-between h-64">
          <div>
            <p id="otpMessage" class="text-6xl text-green-600 font-extrabold">{{ otpMessage }}</p>
            <button @click="copyToClipboard(otpMessage)" class="absolute top-4 right-4 flex items-center">
              <span class="flex items-center border-2 border-gray-300 rounded-full p-2 transition">
                <img v-if="!copied" src="@/assets/copy.svg" alt="Kopyala" class="w-5 h-5" />
                <img v-else src="@/assets/ok.svg" class="w-5 h-5">
              </span>
            </button>
          </div>
          <p class="mt-12">Yukardaki kodu alıcıya gönderiniz.</p>
        </div>
        <div class="mt-4 text-xs text-gray-600 flex justify-between">
          <p>Accepted file types: All Types</p>
          <p>Max files: 20 | Max file size: 2GB</p>
      </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.fixed-size {
  width: 656px;  
  height: 531px;  /*w 492px h 398px 1080p*/
}
.fixed-size-border2 {
  width: 600px;
  height: 319px;  /*w 450px h 239px 1080p*/
}
</style>