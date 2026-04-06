<script setup>
import { computed } from 'vue'

const props = defineProps({
  variant: {
    type: String,
    default: 'primary',
    validator: (value) => ['primary', 'success', 'danger', 'warning', 'purple', 'indigo', 'pink', 'cyan'].includes(value)
  },
  size: {
    type: String,
    default: 'md',
    validator: (value) => ['sm', 'md', 'lg', 'xl'].includes(value)
  },
  pulse: {
    type: Boolean,
    default: false
  }
})

const sizeClasses = computed(() => {
  const sizes = {
    sm: 'w-8 h-8',
    md: 'w-12 h-12',
    lg: 'w-16 h-16',
    xl: 'w-20 h-20'
  }
  return sizes[props.size]
})

const variantClasses = computed(() => {
  return `icon-${props.variant}`
})
</script>

<template>
  <div
    :class="[
      'icon-container',
      sizeClasses,
      variantClasses,
      pulse ? 'animate-pulse' : 'group-hover:scale-110'
    ]"
  >
    <slot />
  </div>
</template>
