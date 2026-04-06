<script setup>
import { ShoppingCart, Users, Package, DollarSign } from 'lucide-vue-next'

const statistics = [
  {
    title: 'Sales',
    stats: '245k',
    icon: ShoppingCart,
    color: 'primary',
  },
  {
    title: 'Customers',
    stats: '12.5k',
    icon: Users,
    color: 'success',
  },
  {
    title: 'Product',
    stats: '1.54k',
    icon: Package,
    color: 'warning',
  },
  {
    title: 'Revenue',
    stats: '$88k',
    icon: DollarSign,
    color: 'info',
  },
]

const colorClasses = {
  primary: 'bg-blue-500/10 text-blue-600 dark:text-blue-500 border-blue-500/20',
  success: 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-500 border-emerald-500/20',
  warning: 'bg-amber-500/10 text-amber-600 dark:text-amber-500 border-amber-500/20',
  info: 'bg-cyan-500/10 text-cyan-600 dark:text-cyan-500 border-cyan-500/20'
}
</script>

<template>
  <div class="card p-6 h-full">
    <div class="flex items-center justify-between mb-6">
      <h5 class="text-lg font-semibold">Transactions</h5>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
      <div
        v-for="stat in statistics"
        :key="stat.title"
        class="flex items-center gap-4"
      >
        <div :class="['w-12 h-12 rounded-xl flex items-center justify-center border flex-shrink-0', colorClasses[stat.color]]">
          <component :is="stat.icon" :size="24" />
        </div>
        <div>
          <p class="text-sm text-muted-foreground">{{ stat.title }}</p>
          <h4 class="text-2xl font-bold tracking-tight">{{ stat.stats }}</h4>
        </div>
      </div>
    </div>
  </div>
</template>
