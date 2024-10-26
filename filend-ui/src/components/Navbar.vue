<template>
<nav class="fixed w-full h-14 bg-white shadow-inner flex items-center justify-between px-4">
    <router-link to="/">
        <img src="@/assets/logo.png" alt="Logo" class="h-10 ml-20" />
    </router-link>
</nav>
</template>

<script>
export default {
    name: "FNavbar",
};
</script>

<style scoped>
.shadow-inner {
    box-shadow: inset 0 -1px 0 rgba(0, 0, 0, 0.1);
}
</style>