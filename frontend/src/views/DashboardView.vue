<script setup>
import { ref, computed, onMounted } from 'vue'
import { Plus, Upload, Wallet, TrendingUp, TrendingDown, CreditCard } from 'lucide-vue-next'
import { ArrowTrendingUpIcon, ArrowTrendingDownIcon, BanknotesIcon, CreditCardIcon } from '@heroicons/vue/24/solid'
import { useDashboardStore } from '@/stores/dashboard'
import { useCategoriesStore } from '@/stores/categories'
import { useToast } from '@/composables/useToast'
import CategoryDonutChart from '@/components/charts/CategoryDonutChart.vue'
import MonthlyFlowChart from '@/components/charts/MonthlyFlowChart.vue'
import TransactionForm from '@/components/transactions/TransactionForm.vue'
import OFXImportModal from '@/components/transactions/OFXImportModal.vue'
import StatCard from '@/components/dashboard/StatCard.vue'
import RecentTransactionsCard from '@/components/dashboard/RecentTransactionsCard.vue'
import MonthSelector from '@/components/dashboard/MonthSelector.vue'

const dashboardStore = useDashboardStore()
const categoriesStore = useCategoriesStore()
const toast = useToast()

const showTransactionForm = ref(false)
const showOFXImport = ref(false)

const onTransactionCreated = () => {
  toast.success('Transação criada com sucesso!')
  dashboardStore.fetchAll()
}

const onOFXImported = () => {
  toast.success('Transações importadas com sucesso!')
  dashboardStore.fetchAll()
}

const formatCurrency = (valueInCents) => {
  return (valueInCents / 100).toLocaleString('pt-BR', {
    style: 'currency',
    currency: 'BRL'
  })
}

const previousMonth = () => {
  let month = dashboardStore.selectedMonth - 1
  let year = dashboardStore.selectedYear

  if (month < 1) {
    month = 12
    year--
  }

  dashboardStore.changeMonth(month, year)
}

const nextMonth = () => {
  let month = dashboardStore.selectedMonth + 1
  let year = dashboardStore.selectedYear

  if (month > 12) {
    month = 1
    year++
  }

  dashboardStore.changeMonth(month, year)
}

// Computed para os cards de estatísticas
const incomeStats = computed(() => ({
  title: 'Total de Entradas',
  value: formatCurrency(dashboardStore.summary.totalIncome),
  change: dashboardStore.summary.incomeChange,
  changeLabel: 'vs mês anterior',
  icon: ArrowTrendingUpIcon,
  color: 'success'
}))

const expenseStats = computed(() => ({
  title: 'Total de Gastos',
  value: formatCurrency(dashboardStore.summary.totalExpense),
  change: dashboardStore.summary.expenseChange,
  changeLabel: 'vs mês anterior',
  icon: ArrowTrendingDownIcon,
  color: 'danger'
}))

const balanceStats = computed(() => ({
  title: 'Saldo do Mês',
  value: formatCurrency(dashboardStore.summary.balance),
  change: 0,
  changeLabel: `${dashboardStore.summary.transactionCount} transações`,
  icon: Wallet,
  color: dashboardStore.summary.balance >= 0 ? 'primary' : 'danger'
}))

const averageStats = computed(() => ({
  title: 'Ticket Médio',
  value: formatCurrency(dashboardStore.summary.averageTicket),
  change: 0,
  changeLabel: 'por transação',
  icon: CreditCardIcon,
  color: 'purple'
}))

// Limitar transações recentes a 8
const recentTransactions = computed(() => {
  return dashboardStore.recentTransactions.slice(0, 8)
})

onMounted(async () => {
  await categoriesStore.fetchCategories()
  await dashboardStore.fetchAll()
})
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4">
      <div>
        <h1 class="text-3xl font-bold tracking-tight mb-2">Dashboard</h1>
        <p class="text-muted-foreground">Visão geral das suas finanças</p>
      </div>

      <div class="flex flex-wrap items-center gap-3">
        <!-- Botão Importar OFX -->
        <button
          @click="showOFXImport = true"
          class="btn btn-outline"
          title="Importar OFX"
        >
          <Upload :size="18" class="text-indigo-600 dark:text-indigo-500" />
          <span class="hidden sm:inline">Importar OFX</span>
        </button>

        <!-- Botão Nova Transação -->
        <button
          @click="showTransactionForm = true"
          class="btn btn-default"
        >
          <Plus :size="18" />
          <span class="hidden sm:inline">Nova Transação</span>
        </button>

        <!-- Seletor de mês -->
        <MonthSelector
          :selected-month="dashboardStore.selectedMonth"
          :selected-year="dashboardStore.selectedYear"
          @previous="previousMonth"
          @next="nextMonth"
        />
      </div>
    </div>

    <!-- Main Grid Layout -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
      <!-- Stats Cards - 8 cols -->
      <div class="lg:col-span-8">
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          <StatCard v-bind="incomeStats" :loading="dashboardStore.loading" />
          <StatCard v-bind="expenseStats" :loading="dashboardStore.loading" />
          <StatCard v-bind="balanceStats" :loading="dashboardStore.loading" />
          <StatCard v-bind="averageStats" :loading="dashboardStore.loading" />
        </div>
      </div>

      <!-- Month Selector Visual - 4 cols -->
      <div class="lg:col-span-4">
        <div class="card p-6 h-full bg-gradient-to-br from-blue-500 to-blue-600 text-white relative overflow-hidden">
          <div class="absolute -top-10 -right-10 w-40 h-40 rounded-full bg-white/10"></div>
          <div class="absolute -bottom-10 -left-10 w-40 h-40 rounded-full bg-white/10"></div>

          <div class="relative z-10">
            <h5 class="text-xl font-bold mb-1">Período Selecionado</h5>
            <p class="text-blue-100 text-sm mb-6">Análise mensal das finanças</p>

            <div class="mb-4">
              <h2 class="text-4xl font-bold mb-1">
                {{ formatCurrency(dashboardStore.summary.balance) }}
              </h2>
              <p class="text-blue-100 text-sm">Saldo Total</p>
            </div>

            <div class="flex items-center gap-4 text-sm">
              <div>
                <p class="text-blue-100">Entradas</p>
                <p class="font-semibold">{{ dashboardStore.summary.transactionCount }}</p>
              </div>
              <div class="h-8 w-px bg-white/20"></div>
              <div>
                <p class="text-blue-100">Categorias</p>
                <p class="font-semibold">{{ dashboardStore.expensesByCategory.length }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Monthly Flow Chart - 8 cols -->
      <div class="lg:col-span-8">
        <MonthlyFlowChart :data="dashboardStore.monthlyFlow" />
      </div>

      <!-- Category Donut Chart - 4 cols -->
      <div class="lg:col-span-4">
        <CategoryDonutChart
          :data="dashboardStore.expensesByCategory"
          title="Gastos por Categoria"
        />
      </div>

      <!-- Recent Transactions - Full width -->
      <div class="lg:col-span-12">
        <RecentTransactionsCard
          :transactions="recentTransactions"
          :loading="dashboardStore.loading"
        />
      </div>
    </div>

    <!-- Modals -->
    <TransactionForm
      v-model:open="showTransactionForm"
      @success="onTransactionCreated"
    />

    <OFXImportModal
      v-model:open="showOFXImport"
      @success="onOFXImported"
    />
  </div>
</template>
