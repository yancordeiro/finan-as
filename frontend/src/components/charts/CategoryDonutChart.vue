<script setup>
import { computed } from 'vue'
import { Doughnut } from 'vue-chartjs'
import {
  Chart as ChartJS,
  ArcElement,
  Tooltip,
  Legend
} from 'chart.js'

ChartJS.register(ArcElement, Tooltip, Legend)

const props = defineProps({
  data: {
    type: Array,
    default: () => []
  },
  title: {
    type: String,
    default: ''
  }
})

// Formatar valor em reais
const formatCurrency = (valueInCents) => {
  return (valueInCents / 100).toLocaleString('pt-BR', {
    style: 'currency',
    currency: 'BRL'
  })
}

// Dados do gráfico
const chartData = computed(() => {
  if (!props.data || props.data.length === 0) {
    return {
      labels: ['Sem dados'],
      datasets: [{
        data: [1],
        backgroundColor: ['#2a2a3a'],
        borderWidth: 0
      }]
    }
  }

  return {
    labels: props.data.map(item => item.categoryName),
    datasets: [{
      data: props.data.map(item => item.total),
      backgroundColor: props.data.map(item => item.categoryColor || '#4f8ef7'),
      borderWidth: 0,
      hoverOffset: 4
    }]
  }
})

// Opções do gráfico
const chartOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  cutout: '65%',
  plugins: {
    legend: {
      display: false
    },
    tooltip: {
      backgroundColor: '#1a1a24',
      titleColor: '#e8e8f0',
      bodyColor: '#e8e8f0',
      borderColor: '#2a2a3a',
      borderWidth: 1,
      padding: 12,
      displayColors: true,
      callbacks: {
        label: (context) => {
          const value = context.raw
          const total = context.dataset.data.reduce((a, b) => a + b, 0)
          const percentage = ((value / total) * 100).toFixed(1)
          return ` ${formatCurrency(value)} (${percentage}%)`
        }
      }
    }
  }
}))

// Total
const total = computed(() => {
  return props.data.reduce((sum, item) => sum + item.total, 0)
})
</script>

<template>
  <div class="card card-glow p-6">
    <h3 class="text-lg font-semibold mb-6">{{ title }}</h3>

    <div class="flex flex-col lg:flex-row items-center gap-6">
      <!-- Gráfico -->
      <div class="relative w-44 h-44 flex-shrink-0">
        <Doughnut :data="chartData" :options="chartOptions" />
        <!-- Total no centro -->
        <div class="absolute inset-0 flex items-center justify-center flex-col">
          <span class="text-xs text-text-muted uppercase tracking-wide">Total</span>
          <span class="text-lg font-bold text-text-primary">{{ formatCurrency(total) }}</span>
        </div>
      </div>

      <!-- Legenda -->
      <div class="flex-1 w-full space-y-1.5 max-h-44 overflow-y-auto scrollbar-hide">
        <div
          v-if="data.length === 0"
          class="text-text-muted text-sm text-center py-4"
        >
          Nenhum dado disponível
        </div>
        <div
          v-for="item in data"
          :key="item.categoryId"
          class="flex items-center justify-between gap-3 p-2 rounded-lg hover:bg-bg-elevated/50 transition-colors"
        >
          <div class="flex items-center gap-2.5 min-w-0">
            <span
              class="w-2.5 h-2.5 rounded-full flex-shrink-0 ring-2 ring-offset-1 ring-offset-bg-card"
              :style="{ backgroundColor: item.categoryColor, ringColor: item.categoryColor + '40' }"
            ></span>
            <span class="text-sm text-text-secondary truncate">
              {{ item.categoryIcon }} {{ item.categoryName }}
            </span>
          </div>
          <div class="flex items-center gap-2 flex-shrink-0">
            <span class="text-sm font-medium text-text-primary">{{ formatCurrency(item.total) }}</span>
            <span class="text-[10px] text-text-muted bg-bg-elevated px-1.5 py-0.5 rounded">{{ item.percentage.toFixed(0) }}%</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
