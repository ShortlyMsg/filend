<template>
  <div id="app" class="min-h-screen bg-gray-200">
    <FilendNavbar />
    <main class="w-full">
      <router-view></router-view>
    </main>
  </div>
</template>

<script setup>
import FilendNavbar from "@/components/FilendNavbar.vue";
//[background:radial-gradient(125%_125%_at_50%_10%,#EDFEF4_40%,#673084_100%)]
</script>

<style>

body {
  overflow: hidden;
}
</style>