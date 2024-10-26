<template>
  <div class="home min-h-screen flex">
    <!-- Sol taraf (Upload) -->
    <div class="flex-1 flex flex-col items-center justify-center bg-blue-100">
      <h2 class="text-3xl font-semibold mb-4 text-blue-800">Dosya Gönder</h2>
      <button 
        @click="goToUpload" 
        class="bg-blue-500 text-white py-3 px-6 rounded-lg shadow-lg hover:bg-blue-700 transition duration-300">
        Gönder
      </button>
    </div>
    <!-- Sağ taraf (Download) -->
    <div class="flex-1 flex flex-col items-center justify-center bg-green-100">
      <h2 class="text-3xl font-semibold mb-4 text-green-800">Dosya Al</h2>
      <button 
        @click="goToDownload" 
        class="bg-green-500 text-white py-3 px-6 rounded-lg shadow-lg hover:bg-green-700 transition duration-300">
        Al
      </button>
    </div>
  </div>
</template>

<script>
export default {
  name: "HomePage",
  methods: {
    goToUpload() {
      this.$router.push({ name: 'FileUpload' }); 
    },
    goToDownload() {
      this.$router.push({ name: 'FileDownload' }); 
    }
  }
}
</script>

<style scoped>

</style>