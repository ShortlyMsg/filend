<template>
  <div id="app" class="min-h-screen [background:radial-gradient(125%_125%_at_50%_10%,#EDFEF4_40%,#673084_100%)]">
    <FilendNavbar />
    <main class="w-full">
      <router-view></router-view>
    </main>
  </div>
</template>

<script setup>
import FilendNavbar from "@/components/FilendNavbar.vue";

</script>

<style>

body {
  overflow: hidden;
}
</style>