<script setup>
import { X } from 'lucide-vue-next'
import { watch } from 'vue'

const props = defineProps({
  open: Boolean,
  title: String,
  maxWidth: {
    type: String,
    default: 'max-w-md'
  }
})

const emit = defineEmits(['update:open', 'close'])

const handleClose = () => {
  emit('update:open', false)
  emit('close')
}

// Fechar com ESC
watch(
  () => props.open,
  (isOpen, _, onCleanup) => {
    if (!isOpen) return
    const handleEsc = (e) => {
      if (e.key === 'Escape') {
        handleClose()
      }
    }
    document.addEventListener('keydown', handleEsc)
    onCleanup(() => document.removeEventListener('keydown', handleEsc))
  }
)
</script>

<template>
  <transition name="dialog">
    <div
      v-if="open"
      class="fixed inset-0 z-50 flex items-center justify-center p-4"
      @click.self="handleClose"
    >
      <!-- Backdrop com blur -->
      <div class="absolute inset-0 bg-black/75 backdrop-blur-md"></div>

      <!-- Dialog -->
      <div
        :class="[
          'relative rounded-lg shadow-2xl border border-border',
          'w-full animate-fade-in flex flex-col max-h-[calc(100vh-2rem)] overflow-hidden',
          maxWidth
        ]"
        style="background-color: hsl(var(--card));"
      >
        <!-- Header -->
        <div class="flex items-center justify-between p-4 sm:p-6 border-b border-border" style="background-color: hsl(var(--card));">
          <h2 class="text-xl font-semibold text-foreground">{{ title }}</h2>
          <button
            @click="handleClose"
            class="text-muted-foreground hover:text-foreground p-2 rounded-lg hover:bg-accent transition-colors"
          >
            <X :size="20" />
          </button>
        </div>

        <!-- Content -->
        <div class="p-4 sm:p-6 overflow-y-auto" style="background-color: hsl(var(--card));">
          <slot />
        </div>

        <!-- Footer (opcional) -->
        <div v-if="$slots.footer" class="p-4 sm:p-6 border-t border-border" style="background-color: hsl(var(--card));">
          <slot name="footer" />
        </div>
      </div>
    </div>
  </transition>
</template>

<style scoped>
.dialog-enter-active,
.dialog-leave-active {
  transition: opacity 0.2s ease;
}

.dialog-enter-from,
.dialog-leave-to {
  opacity: 0;
}

.dialog-enter-active .relative,
.dialog-leave-active .relative {
  transition: transform 0.2s ease, opacity 0.2s ease;
}

.dialog-enter-from .relative,
.dialog-leave-to .relative {
  transform: scale(0.95);
  opacity: 0;
}
</style>
