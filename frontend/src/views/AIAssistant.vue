<template>
  <div class="ai-assistant-view">
    <AIChat />
  </div>
</template>

<script setup>
import AIChat from '@/components/ai/AIChat.vue'
</script>

<style scoped>
.ai-assistant-view {
  height: calc(100vh - 56px - 2rem);
  max-width: 1200px;
  margin: 0 auto;
}

@media (min-width: 1024px) {
  .ai-assistant-view {
    height: calc(100vh - 56px - 3rem);
  }
}
</style>
