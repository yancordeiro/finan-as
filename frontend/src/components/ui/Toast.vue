<script setup>
import { CheckCircle, XCircle, AlertCircle, Info, X } from 'lucide-vue-next'
import { useToast } from '@/composables/useToast'

const { toasts, removeToast } = useToast()

const getIcon = (type) => {
  switch (type) {
    case 'success':
      return CheckCircle
    case 'error':
      return XCircle
    case 'warning':
      return AlertCircle
    default:
      return Info
  }
}

const getColorClass = (type) => {
  switch (type) {
    case 'success':
      return 'bg-success/20 border-success text-success'
    case 'error':
      return 'bg-danger/20 border-danger text-danger'
    case 'warning':
      return 'bg-yellow-500/20 border-yellow-500 text-yellow-500'
    default:
      return 'bg-primary/20 border-primary text-primary'
  }
}
</script>

<template>
  <div class="fixed top-4 right-4 z-50 space-y-2 max-w-md w-full pointer-events-none">
    <transition-group name="toast">
      <div
        v-for="toast in toasts"
        :key="toast.id"
        :class="[
          'flex items-start gap-3 p-4 rounded-lg border shadow-lg pointer-events-auto',
          'bg-bg-secondary backdrop-blur-sm',
          getColorClass(toast.type),
          'animate-slide-in'
        ]"
      >
        <component :is="getIcon(toast.type)" :size="20" class="flex-shrink-0 mt-0.5" />

        <p class="flex-1 text-sm text-text-primary">{{ toast.message }}</p>

        <button
          @click="removeToast(toast.id)"
          class="flex-shrink-0 opacity-70 hover:opacity-100 transition-opacity"
        >
          <X :size="16" />
        </button>
      </div>
    </transition-group>
  </div>
</template>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from {
  opacity: 0;
  transform: translateX(100%);
}

.toast-leave-to {
  opacity: 0;
  transform: translateX(100%);
}
</style>
