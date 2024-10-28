<template>
  <div id="app">
    <Navbar />
    <main class="pt-14">
      <router-view></router-view>
    </main>
  </div>
</template>

<script setup>
import Navbar from "@/components/Navbar.vue";

</script>

<style>

</style>