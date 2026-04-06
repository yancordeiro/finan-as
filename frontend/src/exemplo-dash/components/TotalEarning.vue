<script setup>
import { DollarSign, TrendingUp } from 'lucide-vue-next'

const earnings = [
  { label: 'Total Sales', amount: '$21,845', percentage: '+18%', color: 'success' },
  { label: 'Total Income', amount: '$38,658', percentage: '+42%', color: 'success' },
  { label: 'Total Expense', amount: '$16,758', percentage: '-8%', color: 'danger' },
]

const colorClasses = {
  success: 'text-emerald-600 dark:text-emerald-500',
  danger: 'text-red-600 dark:text-red-500'
}
</script>

<template>
  <div class="card p-6 h-full">
    <div class="flex items-center justify-between mb-6">
      <h5 class="text-lg font-semibold">Total Earning</h5>
    </div>

    <div class="text-center mb-8">
      <h2 class="text-4xl font-bold text-blue-600 dark:text-blue-500 mb-2">$24,895</h2>
      <div class="inline-flex items-center gap-1 text-sm font-medium text-emerald-600 dark:text-emerald-500">
        <TrendingUp :size="16" />
        <span>+10% from last month</span>
      </div>
    </div>

    <div class="space-y-4">
      <div
        v-for="(earning, index) in earnings"
        :key="index"
        class="flex items-center justify-between py-3 border-b last:border-b-0"
      >
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-full bg-muted/50 flex items-center justify-center">
            <DollarSign :size="20" class="text-muted-foreground" />
          </div>
          <span class="text-sm font-medium">{{ earning.label }}</span>
        </div>
        <div class="text-right">
          <p class="font-semibold">{{ earning.amount }}</p>
          <p :class="['text-sm font-medium', colorClasses[earning.color]]">
            {{ earning.percentage }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
