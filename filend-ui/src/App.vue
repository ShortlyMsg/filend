<template>
  <div class="container mx-auto">
    <router-view></router-view>
  </div>
</template>

<script>
export default {
  name: 'App',
};
</script>

<style scoped>
.container {
  @apply px-4 py-6;
}
</style>
