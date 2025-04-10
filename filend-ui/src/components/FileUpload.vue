<script setup>
import { ref } from 'vue';
import axios from 'axios';
import FileIcon from '@/utils/FileIcon.vue';
import { API_ENDPOINTS } from '@/utils/api';

const currentStep = ref(1); // 1: Dosya Seçimi, 2: Önizleme, 3: OTP Gösterimi
const selectedFiles = ref([]);
const percentage = ref({});
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
function showShareOptions() {
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

// function validateFiles(files, selectedFiles) {
//   const maxFileLimit = 20;
//   const maxTotalSize = 2 * 1024 * 1024 * 1024; // 2GB

//   const filesArray = Array.from(files);

//   if (selectedFiles.length + filesArray.length > maxFileLimit) {
//     return "Maksimum dosya limiti aşıldı.";
//   }

//   const totalSize = selectedFiles.reduce((acc, file) => acc + file.size, 0)
//     + files.reduce((acc, file) => acc + file.size, 0);
//   if (totalSize > maxTotalSize) {
//     return "Toplam dosya boyutu 2 GB'ı geçemez.";
//   }

//   return null; // Geçerli dosyalar
// }

function handleFileUpload(event) {
  const files = event.target.files;
  const newFiles = Array.from(files);

  // const errorMessage = validateFiles(newFiles, selectedFiles.value);
  // if (errorMessage) {
  //   alert(errorMessage);
  //   return;
  // }
  // newFiles.forEach(file => {
  //   if (!selectedFiles.value.some(selectedFile => selectedFile.name === file.name)) {
  //     selectedFiles.value.push({ ...file, isUpload: false });
  //   }
  // });
  selectedFiles.value = [...selectedFiles.value, ...newFiles];
  currentStep.value = 2; // Dosya önizleme adımına geç
  uploadFiles(newFiles);
}

function handleDrop(event) {
  event.preventDefault()
  const files = event.dataTransfer.files;
  handleFileUpload({ target: { files } });
}

const chunkSize = 1024 * 1024;

async function uploadFiles(filesToUpload) {
  if (filesToUpload.length === 0) return;

  let otp = otpMessage.value;

  if (!otp) {
    try {
      const otpResponse = await fetch(API_ENDPOINTS.GENERATE_OTP, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (!otpResponse.ok) throw new Error("OTP oluşturulamadı.");

      const otpData = await otpResponse.json();
      otp = otpData.otp;

      if (!otp) {
        otpMessage.value = "OTP alınamadı.";
        return;
      }

      otpMessage.value = otp; // OTP'yi göster
    } catch (error) {
      console.error("OTP alma hatası:", error);
      otpMessage.value = "OTP catch err";
      return;
    }
  }

  // Dosyaları yüklemeye başla
  for (const file of filesToUpload) {
    const fileHash = await calculateSHA256(file);
    const totalChunks = Math.ceil(file.size / chunkSize)

    // Hash kontrol isteği
    const hashCheckResponse = await fetch(API_ENDPOINTS.CHECK_FILE_HASH, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ fileHash }),
    });
    const hashCheckData = await hashCheckResponse.json();

    for (let chunkIndex = 0; chunkIndex < totalChunks; chunkIndex++) {
      const chunkStart = chunkIndex * chunkSize;
      const chunkEnd = Math.min(chunkStart + chunkSize, file.size);
      const chunk = file.slice(chunkStart, chunkEnd);


      // Dosya yükleme işlemi için formData oluştur
      const formData = new FormData();
      if (hashCheckData.fileStatus[fileHash]) {
        formData.append("chunk", chunk);
        formData.append("files", file);
        formData.append("fileHash", fileHash);
        formData.append("chunkIndex", chunkIndex.toString());
        formData.append("totalChunks", totalChunks.toString());
        formData.append("otp", otp);
      } else {
        formData.append("fileName", file.name);
        formData.append("fileHash", fileHash);
      }

      try {
        const response = await axios.post(`${API_ENDPOINTS.UPLOAD_FILES}?otp=${otp}`, formData, {
          headers: { "Content-Type": "multipart/form-data" },
          onUploadProgress: progressEvent => {
            const totalBytes = progressEvent.total || 0;
            const uploadedBytes = progressEvent.loaded || 0;

            const uploadProgress = Math.round((uploadedBytes / totalBytes) * 100);
            const uploadedMB = (uploadedBytes / (1024 * 1024)).toFixed(2);
            const totalMB = (file.size / (1024 * 1024)).toFixed(2);

            percentage.value[file.name] = {
              uploadProgress,
              uploadedMB,
              totalMB,
            };

            console.log(`Dosya: ${file.name} - Yüzde: ${uploadProgress}%`);
            try {
              fetch(`${API_ENDPOINTS.SEND_PROGRESS}`, {
                method: "POST",
                headers: {
                  "Content-Type": "application/json",
                },
                body: JSON.stringify({
                  otp,
                  fileName: file.name,
                  uploadedMB,
                  totalMB,
                  progress: uploadProgress,
                })
              })
            } catch (error) {
              console.error("Firebase'e yükleme durumu gönderme hatası:", error);
            }
          },
        });

        if (response.data.success) {
          console.log(`Dosya ${file.name} başarıyla yüklendi.`);
        }
      } catch (error) {
        console.error(`Dosya ${file.name} yüklenirken hata oluştu:`, error);
        otpMessage.value = "Sunucu ile iletişimde bir sorun var.";
      }
    }
  }
}

</script>

<template>
  <div class="flex justify-center items-center h-screen">
    <div class="bg-white rounded-lg shadow-lg p-6 max-w-xl w-full">

      <div v-if="currentStep === 1"> <!--  || selectedFiles.length === 0 -->
        <h2 class="text-2xl font-bold mb-4">Filend - File Send</h2>
        <p class="text-sm text-gray-600 mb-2">Tek tıkla gönder, tek kodla al!</p>
        <!-- CS 1 Upload -->
        <div class="border-2 border-dashed border-gray-300 rounded-lg p-24 text-center" @dragenter.prevent="dragEnter"
          @dragover.prevent="dragOver" @dragleave.prevent="dragLeave" @drop.prevent="handleDrop">
          <p class="text-gray-600">Dosyayı Buraya Sürükle Bırak</p>
          <p class="text-gray-600">Yada</p>
          <label for="file-upload"
            class="cursor-pointer text-blue-600 border border-blue-600 rounded-md px-4 py-2 inline-block hover:bg-blue-600 hover:text-white transition">
            Dosya Seç
          </label>
          <input id="file-upload" type="file" class="hidden" @change="handleFileUpload" multiple />
        </div>
        <div class="mt-4 text-xs text-gray-600 flex justify-between">
          <p>Accepted file types: All Types</p>
          <p>Max files: 20 | Max file size: 2GB</p>
        </div>
      </div>

      <div class="relative" v-else-if="currentStep === 2">
        <!-- CS 2 Önizleme -->
        <div class="flex items-center justify-between mb-1">
          <h2 class="text-2xl font-bold">Filend - File Send</h2>
          <p id="otpMessage" class="flex items-center justify-end border-2 border-dashed px-1 pb-1
          border-gray-400 rounded text-3xl text-green-600 font-extrabold">{{ otpMessage }}</p>
        </div>

        <!-- Kopyalanan Bildirimi -->
        <div v-if="copied" class="fixed top-16 right-8 bg-blue-500 text-white px-6 py-4 rounded shadow-md transition">
          Kopyalandı!
        </div>

        <!-- Paylaş Butonları -->
        <div class="flex items-center mb-1">
          <p class="text-sm text-gray-600 mb-2">Tek tıkla gönder, tek kodla al!</p>

          <button @click="copyToClipboard(otpMessage)"
            class="ml-auto flex items-center border-2 border-gray-300 rounded-full p-2 transition">
            <img v-if="!copied" src="@/assets/copy-icon.svg" alt="Kopyala" class="w-5 h-5" />
            <img v-else src="@/assets/ok.svg" class="w-5 h-5">
          </button>

          <button @click="showShareOptions(otpMessage)"
            class="flex items-center border-2 border-gray-300 rounded-full p-2 transition">
            <img src="@/assets/share-icon.svg" alt="Paylaş" class="w-5 h-5" />
          </button>

          <div v-if="showOptions"
            class="absolute bg-white right-0 top-22 flex-col space-y-1 rounded-lg p-0 shadow-lg z-10">
            <button @click="shareViaMail" class="flex items-center justify-center  p-2 transition">
              <img src="@/assets/mail-icon.svg" alt="Mail" class="w-5 h-5" />
            </button>
            <button @click="shareViaWhatsapp" class="flex items-center justify-center p-2 transition">
              <img src="@/assets/whatsapp-icon.svg" alt="WhatsApp" class="w-5 h-5" />
            </button>
            <button @click="shareViaTelegram" class="flex items-center justify-center p-2 transition">
              <img src="@/assets/telegram-icon.svg" alt="WhatsApp" class="w-5 h-5" />
            </button>
          </div>
        </div>

        <div class="border-2 border-gray-300 rounded-lg p-6 text-center h-80 relative">
          <div class="overflow-y-auto scrollbar-hidden h-60 p-2">
            <div v-for="(file, index) in selectedFiles" :key="index" class="mb-4">
              <div class="flex items-center">
                <FileIcon :fileName="file.name || 'file-icon.svg'" class="32px" />
                <div class="flex flex-col ml-4 w-full">
                  <div class="flex items-center">
                    <span class="text-sm">{{ file.name }}</span>
                    <button @click="removeFile(index)"
                      class="ml-auto font-extrabold text-red-500 hover:text-red-700">✕</button>
                  </div>
                  <div class="w-full bg-gray-200 rounded-full h-2 mt-1">
                    <div :style="{ width: `${percentage[file.name]?.uploadProgress || 0}%` }"
                      class="bg-blue-600 h-2 rounded-full">
                    </div>
                  </div>
                  <span class="text-xs text-left mt-1">
                    {{
                      percentage[file.name]?.uploadProgress
                        ? `${percentage[file.name]?.uploadedMB || "0.00"} MB / ${percentage[file.name]?.totalMB || "0.00"
                        } MB`
                        : `0.00 MB / ${(file.size / (1024 * 1024)).toFixed(2)} MB`
                    }}
                  </span>
                </div>
              </div>
            </div>
          </div>
          <div class="absolute bottom-2 left-2 flex">
            <label
              class="px-2 py-1 border-2 border-purple-600 text-purple-600 rounded hover:bg-purple-600 hover:text-white cursor-pointer">
              <input type="file" multiple hidden @change="handleFileUpload" />
              Upload More Files
            </label>
          </div>
        </div>


        <div class="mt-4 text-xs text-gray-600 flex justify-between">
          <p>Accepted file types: All Types</p>
          <p>Max files: 20 | Max file size: 2GB</p>
        </div>
        <div class="mt-4 flex justify-end space-x-3">
          <button @click="goBackStep"
            class="px-4 py-1 border-2 border-red-600 text-red-600 rounded hover:bg-red-600 hover:text-white cursor-pointer">
            İptal Et
          </button>
        </div>
      </div>

    </div>
  </div>
</template>

<style scoped>
.scrollbar-hidden {
  scrollbar-width: none;
  /* Firefox için */
  -ms-overflow-style: none;
  /* IE için */
}

.scrollbar-hidden::-webkit-scrollbar {
  display: none;
  /* Chrome ve Safari için */
}

.fixed-size {
  width: 656px;
  height: 531px;
  /*w 492px h 398px 1080p*/
}

.fixed-size-border2 {
  width: 600px;
  height: 319px;
  /*w 450px h 239px 1080p*/
}
</style>
