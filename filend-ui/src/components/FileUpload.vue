<script setup>
import { ref } from 'vue';
import axios from 'axios';
import FileIcon from '@/utils/FileIcon.vue';

const currentStep = ref(2); // 1: Dosya Seçimi, 2: Önizleme, 3: OTP Gösterimi
const selectedFiles = ref([]);
const otpMessage = ref("");
const copied = ref(false);
const showOptions = ref(false);

function copyToClipboard(text) {
  navigator.clipboard.writeText(text).then(() => {
    copied.value = true;
    setTimeout(() => {
      copied.value = false;
    }, 2000);
  });
}
function showShareOptions(){
  showOptions.value = !showOptions.value;
}
function shareViaMail() {
  const mailtoLink = `mailto:?subject=OTP&body=${otpMessage.value}`;
  window.location.href = mailtoLink;
}
function shareViaWhatsapp() {
  const whatsappLink = `https://wa.me/?text=${otpMessage.value}`;
  window.open(whatsappLink, '_blank');
}
function shareViaTelegram() {
  const telegramLink = `https://t.me/share/url?url=${encodeURIComponent(otpMessage.value)}&text=OTP%20Kodunuz`;
  window.open(telegramLink, '_blank');
}

function goBackStep() {
  if (currentStep.value > 1) {
    currentStep.value--;
  }
}

function onFileChange(event) {
  const newFiles = Array.from(event.target.files)
  // .filter(file => {
  //   const maxSize = 2 * 1024 * 1024 * 1024; // 2 GB
  //   //const allowedTypes = ["image/png", "image/jpeg", "application/pdf"];
  //   if (file.size > maxSize && file > 20) {
  //     alert(`${file.name} boyutu sınırı aşıyor.`);
  //     return false;
  //   }
  //   return true;
  // }).map(file => ({
  //   file,
  //   uploadProgress: 0,
  // }));
  selectedFiles.value = [...selectedFiles.value, ...newFiles];
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


  try {
    const response = await axios.post("http://localhost:9091/upload", formData, {
      headers: { "Content-Type": "multipart/form-data" },
      onUploadProgress: progressEvent => {
        const totalBytes = progressEvent.total || 0;
        const uploadedBytes = progressEvent.loaded || 0;

        selectedFiles.value.forEach((file, index) => {
          const uploadProgress = Math.round((uploadedBytes / totalBytes) * 100);
          selectedFiles.value[index].uploadProgress = uploadProgress;
          console.log(`Dosya: ${file.name} - Yüzde: ${uploadProgress }%`);
        });
      },
    });

    if (response.data.otp) {
      otpMessage.value = response.data.otp;
      currentStep.value = 3; // OTP gösterim adımına geç
    } else {
      otpMessage.value = "Hata: OTP alınamadı.";
    }
  } catch (error) {
    console.error("Hata:", error);
    otpMessage.value = "Sunucu ile iletişimde bir sorun var.";
  }
}
</script>

<template>
  <div class="flex justify-center items-center h-screen">
    <div class="bg-white rounded-lg shadow-lg p-6 max-w-xl w-full">
      <h2 class="text-2xl font-bold mb-4">Filend - File Send</h2>
      <p class="text-sm text-gray-600 mb-2">Tek tıkla gönder, tek kodla al!</p>
      <div v-if="currentStep === 1 || selectedFiles.length === 0">
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
          <input id="file-upload" type="file" class="hidden" @change="handleFileUpload" multiple />
        </div>
        <div class="mt-4 text-xs text-gray-600 flex justify-between">
          <p>Accepted file types: All Types</p>
          <p>Max files: 20 | Max file size: 2GB</p>
        </div>
      </div>

      <div v-else-if="currentStep === 2">
        <!-- CS 2 Önizleme -->
        <div class="border-2 border-gray-300 rounded-lg p-6 text-center h-80 relative">
          <div class="overflow-y-auto scrollbar-hidden h-60 p-2">
            <div v-for="(file, index) in selectedFiles" :key="index" class="mb-4">
              <div class="flex items-center">
                <FileIcon :fileName="file.name" class="32px"/>
                <div class="flex flex-col ml-4 w-full">
                  <div class="flex items-center">
                    <span class="text-sm">{{ file.name }}</span>
                    <button @click="removeFile(index)" class="ml-auto font-extrabold text-red-500 hover:text-red-700">✕</button>
                  </div>
                  <div class="w-full bg-gray-200 rounded-full h-2 mt-1">
                    <div :style="{ width: `${file.uploadProgress || 0}%` }" class="bg-blue-600 h-2 rounded-full"></div>
                  </div>
                  <span class="text-xs text-left mt-1">
                    {{ file.uploadProgress ? `${(file.uploadProgress || 0).toFixed(2)} MB / ${(file.size / (1024 * 1024))
                    .toFixed(2)} MB` : `0.00 MB / ${(file.size / (1024 * 1024)).toFixed(2)} MB` }}
                  </span>
                </div>
              </div>
            </div>
          </div>
          <div class="absolute bottom-2 left-2 flex">
            <label class="px-2 py-1 border-2 border-purple-600 text-purple-600 rounded hover:bg-purple-600 hover:text-white cursor-pointer">
            <input type="file" multiple hidden @change="onFileChange"/>
            Upload More Files
            </label>
          </div>
        </div>


        <div class="mt-4 text-xs text-gray-600 flex justify-between">
          <p>Accepted file types: All Types</p>
          <p>Max files: 20 | Max file size: 2GB</p>
        </div> 
        <div class="mt-4 flex justify-end space-x-3">
          <button @click="goBackStep" class="px-4 py-2 text-gray-600 hover:text-gray-800">
            Geri
          </button>
          <button @click="uploadFiles" class="px-4 py-2 bg-blue-600 text-white rounded "
          :disabled="selectedFiles.length === 0" :class="{'bg-gray-400': selectedFiles.length === 0}" >
            Gönder
          </button>
        </div>
      </div>

      <div v-else-if="currentStep === 3">
        <!-- CS 3 OTP-->
        <div v-if="copied" class="fixed top-16 right-8 bg-blue-500 text-white px-6 py-4 rounded shadow-md transition">
          Kopyalandı!
        </div>
        <div class="border-2 border-gray-300 rounded-lg p-24 text-center relative flex flex-col justify-between h-64">
          <div v-if="showOptions" class="absolute top-14 right-4 flex-col">
            <button @click="shareViaMail" class="flex items-center justify-center border-2 border-[#e2abab] rounded-full p-2 transition">
              <img src="@/assets/mail-icon.svg" alt="Mail" class="w-5 h-5" />
            </button>
            <button @click="shareViaWhatsapp" class="flex items-center justify-center border-2 border-[#e2abab] rounded-full p-2 transition">
              <img src="@/assets/whatsapp-icon.svg" alt="WhatsApp" class="w-5 h-5" />
            </button>
            <button @click="shareViaTelegram" class="flex items-center justify-center border-2 border-[#e2abab] rounded-full p-2 transition">
              <img src="@/assets/telegram-icon.svg" alt="WhatsApp" class="w-5 h-5" />
            </button>
          </div>
          <div>
            <p id="otpMessage" class="text-6xl text-green-600 font-extrabold">{{ otpMessage }}</p>
            <button @click="showShareOptions(otpMessage)" class="absolute top-4 right-4 flex items-center">
              <span class="flex items-center border-2 border-gray-300 rounded-full p-2 transition">
                <img src="@/assets/share-icon.svg" alt="Paylaş" class="w-5 h-5" />
              </span>
            </button>

            <button @click="copyToClipboard(otpMessage)" class="absolute top-4 right-16 flex items-center">
              <span class="flex items-center border-2 border-gray-300 rounded-full p-2 transition">
                <img v-if="!copied" src="@/assets/copy-icon.svg" alt="Kopyala" class="w-5 h-5" />
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

.scrollbar-hidden {
  scrollbar-width: none; /* Firefox için */
  -ms-overflow-style: none; /* IE için */
}
.scrollbar-hidden::-webkit-scrollbar {
  display: none; /* Chrome ve Safari için */
}
.fixed-size {
  width: 656px;  
  height: 531px;  /*w 492px h 398px 1080p*/
}
.fixed-size-border2 {
  width: 600px;
  height: 319px;  /*w 450px h 239px 1080p*/
}
</style>