<script setup>
import { ref, computed } from 'vue'
import { Bar } from 'vue-chartjs'
import { Chart as ChartJS, Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale } from 'chart.js'

ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale)

const isDark = ref(document.documentElement.classList.contains('dark'))

const chartData = {
  labels: ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'],
  datasets: [
    {
      data: [37, 57, 45, 75, 57, 40, 65],
      backgroundColor: (context) => {
        const index = context.dataIndex
        return index === 3 ? 'rgb(59, 130, 246)' : isDark.value ? 'rgba(100, 116, 139, 0.3)' : 'rgba(226, 232, 240, 1)'
      },
      borderRadius: 8,
      barThickness: 24,
    }
  ]
}

const chartOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      display: false
    },
    tooltip: {
      backgroundColor: isDark.value ? 'rgba(15, 23, 42, 0.9)' : 'rgba(255, 255, 255, 0.9)',
      titleColor: isDark.value ? '#e2e8f0' : '#1e293b',
      bodyColor: isDark.value ? '#cbd5e1' : '#475569',
      borderColor: isDark.value ? 'rgba(71, 85, 105, 0.2)' : 'rgba(226, 232, 240, 0.8)',
      borderWidth: 1,
      padding: 12,
      displayColors: false,
      callbacks: {
        label: (context) => `${context.parsed.y}k`
      }
    }
  },
  scales: {
    x: {
      grid: {
        display: false
      },
      ticks: {
        display: false
      },
      border: {
        display: false
      }
    },
    y: {
      beginAtZero: true,
      max: 80,
      ticks: {
        stepSize: 20,
        callback: (value) => `${value}k`,
        color: isDark.value ? '#64748b' : '#94a3b8',
        font: {
          size: 12
        }
      },
      grid: {
        color: isDark.value ? 'rgba(71, 85, 105, 0.2)' : 'rgba(226, 232, 240, 0.8)',
        drawTicks: false
      },
      border: {
        display: false,
        dash: [5, 5]
      }
    }
  }
}))
</script>

<template>
  <div class="card p-6 h-full flex flex-col">
    <div class="flex items-center justify-between mb-4">
      <h5 class="text-lg font-semibold">Weekly Overview</h5>
    </div>

    <div class="flex-1 min-h-[200px] mb-6">
      <Bar :data="chartData" :options="chartOptions" />
    </div>

    <div class="space-y-4">
      <div class="flex items-start gap-3">
        <h4 class="text-3xl font-bold text-blue-600 dark:text-blue-500">45%</h4>
        <p class="text-sm text-muted-foreground flex-1 pt-1">
          Your sales performance is 45% 😎 better compared to last month
        </p>
      </div>

      <button class="btn btn-default w-full">
        Details
      </button>
    </div>
  </div>
</template>
