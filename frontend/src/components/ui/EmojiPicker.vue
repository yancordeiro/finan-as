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
    return props.modelValue || '📦'
  },
  set(val) {
    emit('update:modelValue', val)
  }
})

// Emojis categorizados
const emojiCategories = {
  'Alimentação': ['🍔', '🍕', '🍜', '🍱', '🥗', '🍝', '🥘', '🍲', '🍗', '🥙', '🌮', '🌯'],
  'Transporte': ['🚗', '🚕', '🚙', '🚌', '🚎', '🏎️', '🚓', '🚑', '🚒', '🚐', '🚛', '🚲', '🛵', '🏍️', '✈️', '🚁', '🚂', '🚇', '🚊'],
  'Moradia': ['🏠', '🏡', '🏢', '🏣', '🏤', '🏥', '🏦', '🏨', '🏩', '🏪', '🏫'],
  'Lazer': ['🎮', '🎯', '🎲', '🎰', '🎳', '🎭', '🎪', '🎨', '🎬', '🎤', '🎧', '🎵', '🎸', '🎹', '🎺', '🎻', '⚽', '🏀', '🏈', '⚾', '🎾'],
  'Saúde': ['🏥', '⚕️', '💊', '💉', '🩺', '🩹', '🩼', '🧬'],
  'Educação': ['📚', '📖', '📝', '✏️', '📓', '📕', '📗', '📘', '📙', '🎓', '🎒', '🏫'],
  'Trabalho': ['💼', '💻', '⌨️', '🖥️', '🖨️', '📱', '☎️', '📞', '📧', '📨'],
  'Compras': ['🛒', '🛍️', '💳', '💰', '💸', '🏪', '🏬'],
  'Outros': ['📦', '🎁', '🎈', '🎉', '🎊', '❤️', '⭐', '✨', '🔥', '💎', '🏆', '🎯']
}
</script>

<template>
  <div class="w-full">
    <label v-if="label" class="block text-sm font-medium text-text-primary mb-2">
      {{ label }}
      <span v-if="required" class="text-danger ml-1">*</span>
    </label>

    <!-- Emoji selecionado -->
    <div class="flex items-center gap-3 mb-3">
      <div class="w-16 h-16 bg-bg-primary border-2 border-border-subtle rounded-lg flex items-center justify-center text-4xl">
        {{ value }}
      </div>
      <input
        type="text"
        v-model="value"
        placeholder="📦"
        class="input flex-1 text-center text-2xl"
        maxlength="2"
      />
    </div>

    <!-- Grid de emojis por categoria -->
    <div class="max-h-64 overflow-y-auto border border-border-subtle rounded-lg p-3 bg-bg-primary">
      <div v-for="(emojis, category) in emojiCategories" :key="category" class="mb-4 last:mb-0">
        <p class="text-xs font-medium text-text-secondary mb-2">{{ category }}</p>
        <div class="grid grid-cols-8 gap-2">
          <button
            v-for="emoji in emojis"
            :key="emoji"
            type="button"
            @click="value = emoji"
            :class="[
              'w-10 h-10 rounded-lg text-2xl transition-all hover:scale-110 hover:bg-bg-secondary',
              value === emoji && 'bg-primary/20 ring-2 ring-primary'
            ]"
          >
            {{ emoji }}
          </button>
        </div>
      </div>
    </div>

    <p v-if="error" class="mt-1 text-sm text-danger">
      {{ error }}
    </p>
  </div>
</template>
