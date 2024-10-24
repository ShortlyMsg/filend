<template>
  <div class="flex justify-center items-center h-screen bg-gray-100">
    <div class="w-full max-w-md bg-white rounded-lg shadow-lg p-6">
      <!-- Başlık -->
      <div class="mb-4">
        <h2 class="text-2xl font-bold">Filend - File Send</h2>
        <p class="text-sm text-gray-600">Tek tıkla gönder, tek kodla al!</p>
        <p class="text-sm text-red-500 mt-2">
          Tek Seferlik Kodu Kaydetmeyi Unutmayınız!
        </p>
      </div>
      <!-- Gönderen -->
      <div class="flex items-center mb-4">
        <div class="w-10 h-10 bg-green-500 rounded-full flex justify-center items-center text-white">
          Fe
        </div>
        <p class="ml-3 text-gray-800">Gönderen</p>
      </div>
      <!-- Dosya Yükleme Alanı -->
      <div class="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center"
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
        <input id="file-upload" type="file" class="hidden" @change="handleFileUpload" />
      </div>
      <!-- Dosya Bilgileri -->
      <div class="mt-4 text-sm text-gray-600">
        <p>Accepted file types: All Types</p>
        <p>Max files: 20 | Max file size: 2GB</p>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "FileUpload",
  data() {
    return {
      selectedFiles: [],
    };
  },
  methods: {
    dragEnter() {
      console.log('Dosya alanına girildi.');
    },
    dragOver() {
      console.log('Dosya alanının üstünde.');
    },
    dragLeave() {
      console.log('Dosya alanından çıkıldı.');
    },
    handleDrop(event) {
      const files = event.dataTransfer.files;
      console.log('Bırakılan dosyalar:', files);
      this.handleFileUpload({ target: { files } });
    },
    handleFileUpload(event) {
      const files = event.target.files;
      this.selectedFiles = Array.from(files); // Seçilen dosyaları listeye ekle
      console.log('Yüklenen dosyalar:', this.selectedFiles);
    },
  },
};
</script>

<style scoped>

</style>