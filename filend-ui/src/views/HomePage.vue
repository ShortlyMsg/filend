<template>
  <div class="home">
    <h1 class="text-4xl font-bold text-center mt-12">Hoş Geldiniz</h1>
    <div class="flex justify-center mt-8 space-x-6">
      <button 
        @click="goToUpload" 
        class="bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-700 transition duration-300">
        Gönder
      </button>
      <button 
        @click="goToDownload" 
        class="bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-700 transition duration-300">
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