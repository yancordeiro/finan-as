<script setup>
import { computed } from 'vue'
import { Bar } from 'vue-chartjs'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend
} from 'chart.js'

ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend)

const props = defineProps({
  data: {
    type: Array,
    default: () => []
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
      labels: [],
      datasets: []
    }
  }

  return {
    labels: props.data.map(item => item.label),
    datasets: [
      {
        label: 'Entradas',
        data: props.data.map(item => item.income / 100), // Converter para reais
        backgroundColor: '#00e676',
        borderRadius: 4,
        barPercentage: 0.7,
        categoryPercentage: 0.8
      },
      {
        label: 'Saídas',
        data: props.data.map(item => item.expense / 100), // Converter para reais
        backgroundColor: '#ff4d4d',
        borderRadius: 4,
        barPercentage: 0.7,
        categoryPercentage: 0.8
      }
    ]
  }
})

// Opções do gráfico
const chartOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  interaction: {
    intersect: false,
    mode: 'index'
  },
  scales: {
    x: {
      grid: {
        display: false
      },
      ticks: {
        color: '#7a7a9a'
      }
    },
    y: {
      beginAtZero: true,
      grid: {
        color: '#2a2a3a'
      },
      ticks: {
        color: '#7a7a9a',
        callback: (value) => {
          if (value >= 1000) {
            return `R$ ${(value / 1000).toFixed(0)}k`
          }
          return `R$ ${value}`
        }
      }
    }
  },
  plugins: {
    legend: {
      position: 'top',
      align: 'end',
      labels: {
        color: '#e8e8f0',
        usePointStyle: true,
        pointStyle: 'circle',
        padding: 20
      }
    },
    tooltip: {
      backgroundColor: '#1a1a24',
      titleColor: '#e8e8f0',
      bodyColor: '#e8e8f0',
      borderColor: '#2a2a3a',
      borderWidth: 1,
      padding: 12,
      callbacks: {
        label: (context) => {
          return ` ${context.dataset.label}: ${formatCurrency(context.raw * 100)}`
        }
      }
    }
  }
}))
</script>

<template>
  <div class="card card-glow p-6">
    <div class="flex items-center justify-between mb-6">
      <h3 class="text-lg font-semibold">Fluxo Mensal</h3>
      <div class="flex items-center gap-4 text-xs">
        <div class="flex items-center gap-1.5">
          <span class="w-2.5 h-2.5 rounded-full bg-success"></span>
          <span class="text-text-muted">Entradas</span>
        </div>
        <div class="flex items-center gap-1.5">
          <span class="w-2.5 h-2.5 rounded-full bg-danger"></span>
          <span class="text-text-muted">Saídas</span>
        </div>
      </div>
    </div>

    <div v-if="data.length === 0" class="h-64 flex items-center justify-center text-text-muted">
      Nenhum dado disponível
    </div>

    <div v-else class="h-64">
      <Bar :data="chartData" :options="chartOptions" />
    </div>
  </div>
</template>
