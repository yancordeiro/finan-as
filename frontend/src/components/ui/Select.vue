<script setup>
import { computed } from 'vue'

const props = defineProps({
  modelValue: [String, Number],
  options: {
    type: Array,
    required: true
  },
  label: String,
  placeholder: String,
  error: String,
  disabled: Boolean,
  required: Boolean,
  valueKey: {
    type: String,
    default: 'value'
  },
  labelKey: {
    type: String,
    default: 'label'
  }
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

    <select
      v-model="value"
      :disabled="disabled"
      :class="[
        'input w-full',
        error && 'border-danger focus:ring-danger',
        disabled && 'opacity-50 cursor-not-allowed'
      ]"
    >
      <option v-if="placeholder" value="" disabled>{{ placeholder }}</option>
      <option
        v-for="option in options"
        :key="option[valueKey]"
        :value="option[valueKey]"
      >
        {{ option[labelKey] }}
      </option>
    </select>

    <p v-if="error" class="mt-1 text-sm text-danger">
      {{ error }}
    </p>
  </div>
</template>
