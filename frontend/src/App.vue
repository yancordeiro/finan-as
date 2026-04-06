<script setup>
import { RouterView } from 'vue-router'
import Toast from '@/components/ui/Toast.vue'
</script>

<template>
  <div class="min-h-screen bg-background">
    <RouterView />
    <Toast />
  </div>
</template>

<style scoped>
</style>
