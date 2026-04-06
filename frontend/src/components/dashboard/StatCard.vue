<script setup>
import { computed } from 'vue'
import { TrendingUp, TrendingDown } from 'lucide-vue-next'

const props = defineProps({
  title: {
    type: String,
    required: true
  },
  value: {
    type: String,
    required: true
  },
  change: {
    type: Number,
    default: 0
  },
  changeLabel: {
    type: String,
    default: 'vs mês anterior'
  },
  icon: {
    type: Object,
    required: true
  },
  color: {
    type: String,
    default: 'primary',
    validator: (value) => ['primary', 'success', 'danger', 'warning', 'purple', 'indigo'].includes(value)
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const isPositive = computed(() => Math.sign(props.change) === 1)

const displayChange = computed(() => {
  if (!isFinite(props.change) || isNaN(props.change)) {
    return null
  }
  return Math.abs(props.change)
})

const shouldShowChange = computed(() => {
  return displayChange.value !== null && displayChange.value !== 0
})

const colorClasses = computed(() => {
  const colors = {
    primary: 'bg-blue-500/10 text-blue-600 dark:text-blue-500 border-blue-500/20',
    success: 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-500 border-emerald-500/20',
    danger: 'bg-red-500/10 text-red-600 dark:text-red-500 border-red-500/20',
    warning: 'bg-amber-500/10 text-amber-600 dark:text-amber-500 border-amber-500/20',
    purple: 'bg-purple-500/10 text-purple-600 dark:text-purple-500 border-purple-500/20',
    indigo: 'bg-indigo-500/10 text-indigo-600 dark:text-indigo-500 border-indigo-500/20'
  }
  return colors[props.color]
})

const textColorClasses = computed(() => {
  const colors = {
    primary: 'text-blue-600 dark:text-blue-500',
    success: 'text-emerald-600 dark:text-emerald-500',
    danger: 'text-red-600 dark:text-red-500',
    warning: 'text-amber-600 dark:text-amber-500',
    purple: 'text-purple-600 dark:text-purple-500',
    indigo: 'text-indigo-600 dark:text-indigo-500'
  }
  return colors[props.color]
})
</script>

<template>
  <div class="card p-6 h-full transition-all hover:shadow-md">
    <!-- Loading state -->
    <div v-if="loading" class="space-y-4">
      <div class="skeleton h-12 w-12 rounded-xl"></div>
      <div class="skeleton h-4 w-24"></div>
      <div class="skeleton h-8 w-32"></div>
      <div class="skeleton h-3 w-20"></div>
    </div>

    <!-- Content -->
    <div v-else>
      <div class="flex items-center justify-between mb-4">
        <div :class="['w-12 h-12 rounded-xl flex items-center justify-center border transition-transform group-hover:scale-110', colorClasses]">
          <component :is="icon" :size="24" />
        </div>
        <div v-if="shouldShowChange" :class="['flex items-center gap-1 text-sm font-medium', isPositive ? 'text-emerald-600 dark:text-emerald-500' : 'text-red-600 dark:text-red-500']">
          <TrendingUp v-if="isPositive" :size="16" />
          <TrendingDown v-else :size="16" />
          <span>{{ displayChange.toFixed(1) }}%</span>
        </div>
      </div>

      <h6 class="text-sm font-medium text-muted-foreground mb-1">{{ title }}</h6>
      <h4 :class="['text-3xl font-bold tracking-tight mb-2', textColorClasses]">{{ value }}</h4>

      <p class="text-sm text-muted-foreground">{{ changeLabel }}</p>
    </div>
  </div>
</template>

<style scoped>
.skeleton {
  @apply relative overflow-hidden bg-muted rounded;
}

.skeleton::after {
  content: '';
  @apply absolute inset-0;
  background: linear-gradient(
    90deg,
    transparent,
    rgba(255, 255, 255, 0.08),
    transparent
  );
  animation: shimmer 2s infinite;
}

@keyframes shimmer {
  0% {
    transform: translateX(-100%);
  }
  100% {
    transform: translateX(100%);
  }
}
</style>
