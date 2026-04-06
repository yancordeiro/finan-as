<script setup>
import { PieChart, FileText, TrendingUp, Wallet } from 'lucide-vue-next'
import AnalyticsAward from '../components/AnalyticsAward.vue'
import AnalyticsTransactions from '../components/AnalyticsTransactions.vue'
import WeeklyOverview from '../components/WeeklyOverview.vue'
import TotalEarning from '../components/TotalEarning.vue'
import CardStatisticsVertical from '../components/CardStatisticsVertical.vue'
import RecentTransactions from '../components/RecentTransactions.vue'

const totalProfit = {
  title: 'Total Profit',
  stats: '$25.6k',
  change: 42,
  subtitle: 'Weekly Project',
  icon: PieChart,
  color: 'purple'
}

const newProject = {
  title: 'New Project',
  stats: '862',
  change: -18,
  subtitle: 'Yearly Project',
  icon: FileText,
  color: 'primary'
}

const salesGrowth = {
  title: 'Sales Growth',
  stats: '+48%',
  change: 22,
  subtitle: 'This Quarter',
  icon: TrendingUp,
  color: 'success'
}

const totalRevenue = {
  title: 'Total Revenue',
  stats: '$156k',
  change: 15,
  subtitle: 'This Month',
  icon: Wallet,
  color: 'warning'
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div>
      <h1 class="text-3xl font-bold tracking-tight mb-2">Example Dashboard</h1>
      <p class="text-muted-foreground">Dashboard inspirado no template Materio</p>
    </div>

    <!-- Main Grid Layout -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
      <!-- Award Card - 4 cols -->
      <div class="lg:col-span-4">
        <AnalyticsAward />
      </div>

      <!-- Transactions - 8 cols -->
      <div class="lg:col-span-8">
        <AnalyticsTransactions />
      </div>

      <!-- Weekly Overview - 4 cols -->
      <div class="lg:col-span-4">
        <WeeklyOverview />
      </div>

      <!-- Total Earning - 4 cols -->
      <div class="lg:col-span-4">
        <TotalEarning />
      </div>

      <!-- Stats Cards Grid - 4 cols -->
      <div class="lg:col-span-4">
        <div class="grid grid-cols-2 gap-4 h-full">
          <div class="col-span-2 sm:col-span-1">
            <CardStatisticsVertical v-bind="totalProfit" />
          </div>

          <div class="col-span-2 sm:col-span-1">
            <CardStatisticsVertical v-bind="newProject" />
          </div>

          <div class="col-span-2 sm:col-span-1">
            <CardStatisticsVertical v-bind="salesGrowth" />
          </div>

          <div class="col-span-2 sm:col-span-1">
            <CardStatisticsVertical v-bind="totalRevenue" />
          </div>
        </div>
      </div>

      <!-- Recent Transactions Table - Full width -->
      <div class="lg:col-span-12">
        <RecentTransactions />
      </div>
    </div>

    <!-- Info Box -->
    <div class="card p-6 bg-blue-500/10 border-blue-500/20">
      <div class="flex items-start gap-4">
        <div class="w-12 h-12 rounded-full bg-blue-500/20 flex items-center justify-center flex-shrink-0">
          <span class="text-2xl">ℹ️</span>
        </div>
        <div>
          <h6 class="font-semibold mb-1 text-blue-600 dark:text-blue-500">Dashboard de Exemplo</h6>
          <p class="text-sm text-muted-foreground">
            Este é um dashboard de exemplo criado na pasta <code class="px-2 py-1 rounded bg-muted text-foreground">exemplo-dash/</code>.
            Ele demonstra componentes reutilizáveis inspirados no template Materio com design moderno usando shadcn/ui.
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
