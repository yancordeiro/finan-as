<script setup>
import { computed } from 'vue'
import { ChevronLeft, ChevronRight } from 'lucide-vue-next'

const props = defineProps({
  selectedMonth: {
    type: Number,
    required: true
  },
  selectedYear: {
    type: Number,
    required: true
  }
})

const emit = defineEmits(['previous', 'next'])

const monthNames = [
  'Janeiro', 'Fevereiro', 'Março', 'Abril', 'Maio', 'Junho',
  'Julho', 'Agosto', 'Setembro', 'Outubro', 'Novembro', 'Dezembro'
]

const selectedPeriod = computed(() => {
  return `${monthNames[props.selectedMonth - 1]} ${props.selectedYear}`
})

const canGoNext = computed(() => {
  const now = new Date()
  const currentMonth = now.getMonth() + 1
  const currentYear = now.getFullYear()

  return props.selectedYear < currentYear ||
    (props.selectedYear === currentYear && props.selectedMonth < currentMonth)
})
</script>

<template>
  <div class="inline-flex items-center gap-1 card p-1">
    <button
      @click="emit('previous')"
      class="p-2 rounded-lg hover:bg-accent transition-colors"
    >
      <ChevronLeft :size="18" />
    </button>

    <span class="px-3 py-1 font-medium min-w-[140px] text-center text-sm">
      {{ selectedPeriod }}
    </span>

    <button
      @click="emit('next')"
      :disabled="!canGoNext"
      class="p-2 rounded-lg hover:bg-accent transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
    >
      <ChevronRight :size="18" />
    </button>
  </div>
</template>
