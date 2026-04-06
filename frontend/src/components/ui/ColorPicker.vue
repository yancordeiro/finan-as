<script setup>
import { computed } from 'vue'

const props = defineProps({
  modelValue: String,
  label: String,
  error: String,
  required: Boolean
})

const emit = defineEmits(['update:modelValue'])

const value = computed({
  get() {
    return props.modelValue || '#4f8ef7'
  },
  set(val) {
    emit('update:modelValue', val)
  }
})

// Cores pré-definidas para quick select
const presetColors = [
  '#FF6B6B', '#4ECDC4', '#45B7D1', '#A259FF',
  '#FFA07A', '#95E1D3', '#FFB6C1', '#FFD700',
  '#00E676', '#00D4FF', '#FFD93D', '#6BCF7F',
  '#4f8ef7', '#a259ff', '#ff4d4d', '#7A7A9A'
]
</script>

<template>
  <div class="w-full">
    <label v-if="label" class="block text-sm font-medium text-text-primary mb-2">
      {{ label }}
      <span v-if="required" class="text-danger ml-1">*</span>
    </label>

    <div class="flex items-center gap-3">
      <!-- Color input nativo -->
      <div class="relative">
        <input
          type="color"
          v-model="value"
          class="w-16 h-10 rounded-lg border-2 border-border-subtle cursor-pointer overflow-hidden"
        />
        <div
          class="absolute inset-0 rounded-lg pointer-events-none"
          :style="{ backgroundColor: value }"
        ></div>
      </div>

      <!-- Input de texto com o hex -->
      <input
        type="text"
        v-model="value"
        placeholder="#000000"
        class="input flex-1 font-mono"
        maxlength="7"
      />
    </div>

    <!-- Preset colors -->
    <div class="mt-3 flex flex-wrap gap-2">
      <button
        v-for="color in presetColors"
        :key="color"
        type="button"
        @click="value = color"
        :class="[
          'w-8 h-8 rounded-lg border-2 transition-all',
          value === color ? 'border-white scale-110' : 'border-border-subtle hover:scale-105'
        ]"
        :style="{ backgroundColor: color }"
        :title="color"
      />
    </div>

    <p v-if="error" class="mt-1 text-sm text-danger">
      {{ error }}
    </p>
  </div>
</template>
