<script setup>
import { computed } from 'vue'

const props = defineProps({
  modelValue: [String, Number],
  type: {
    type: String,
    default: 'text'
  },
  placeholder: String,
  label: String,
  error: String,
  disabled: Boolean,
  required: Boolean
})

const emit = defineEmits(['update:modelValue'])

const value = computed({
  get() {
    return props.modelValue
  },
  set(val) {
    emit('update:modelValue', val)
  }
})
</script>

<template>
  <div class="w-full">
    <label v-if="label" class="block text-sm font-medium text-text-primary mb-2">
      {{ label }}
      <span v-if="required" class="text-danger ml-1">*</span>
    </label>

    <input
      v-model="value"
      :type="type"
      :placeholder="placeholder"
      :disabled="disabled"
      :class="[
        'input w-full',
        error && 'border-danger focus:ring-danger',
        disabled && 'opacity-50 cursor-not-allowed'
      ]"
    />

    <p v-if="error" class="mt-1 text-sm text-danger">
      {{ error }}
    </p>
  </div>
</template>
