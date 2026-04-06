import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useApi } from '@/composables/useApi'

export const useDashboardStore = defineStore('dashboard', () => {
  const api = useApi()

  // Estado
  const summary = ref({
    totalIncome: 0,
    totalExpense: 0,
    balance: 0,
    incomeChange: 0,
    expenseChange: 0,
    transactionCount: 0,
    averageTicket: 0
  })

  const expensesByCategory = ref([])
  const incomesByCategory = ref([])
  const monthlyFlow = ref([])
  const recentTransactions = ref([])

  const loading = ref(false)
  const error = ref(null)

  // Mês/ano selecionado
  const selectedMonth = ref(new Date().getMonth() + 1)
  const selectedYear = ref(new Date().getFullYear())

  // Buscar resumo do mês
  async function fetchSummary(month = null, year = null) {
    const m = month || selectedMonth.value
    const y = year || selectedYear.value

    try {
      const response = await api.get(`/api/dashboard/summary?month=${m}&year=${y}`)
      summary.value = response.data.data
      return { success: true }
    } catch (err) {
      console.error('Erro ao buscar resumo:', err)
      error.value = err.response?.data?.error || 'Erro ao carregar resumo'
      return { success: false, error: error.value }
    }
  }

  // Buscar gastos por categoria
  async function fetchExpensesByCategory(month = null, year = null) {
    const m = month || selectedMonth.value
    const y = year || selectedYear.value

    try {
      const response = await api.get(`/api/dashboard/expenses-by-category?month=${m}&year=${y}`)
      expensesByCategory.value = response.data.data || []
      return { success: true }
    } catch (err) {
      console.error('Erro ao buscar gastos por categoria:', err)
      error.value = err.response?.data?.error || 'Erro ao carregar gastos por categoria'
      return { success: false, error: error.value }
    }
  }

  // Buscar entradas por categoria
  async function fetchIncomesByCategory(month = null, year = null) {
    const m = month || selectedMonth.value
    const y = year || selectedYear.value

    try {
      const response = await api.get(`/api/dashboard/incomes-by-category?month=${m}&year=${y}`)
      incomesByCategory.value = response.data.data || []
      return { success: true }
    } catch (err) {
      console.error('Erro ao buscar entradas por categoria:', err)
      error.value = err.response?.data?.error || 'Erro ao carregar entradas por categoria'
      return { success: false, error: error.value }
    }
  }

  // Buscar fluxo mensal
  async function fetchMonthlyFlow(months = 6) {
    try {
      const response = await api.get(`/api/dashboard/monthly-flow?months=${months}`)
      monthlyFlow.value = response.data.data || []
      return { success: true }
    } catch (err) {
      console.error('Erro ao buscar fluxo mensal:', err)
      error.value = err.response?.data?.error || 'Erro ao carregar fluxo mensal'
      return { success: false, error: error.value }
    }
  }

  // Buscar transações recentes
  async function fetchRecentTransactions(limit = 10) {
    try {
      const response = await api.get(`/api/dashboard/recent-transactions?limit=${limit}`)
      recentTransactions.value = response.data.data || []
      return { success: true }
    } catch (err) {
      console.error('Erro ao buscar transações recentes:', err)
      error.value = err.response?.data?.error || 'Erro ao carregar transações recentes'
      return { success: false, error: error.value }
    }
  }

  // Buscar todos os dados do dashboard
  async function fetchAll() {
    loading.value = true
    error.value = null

    try {
      await Promise.all([
        fetchSummary(),
        fetchExpensesByCategory(),
        fetchIncomesByCategory(),
        fetchMonthlyFlow(),
        fetchRecentTransactions()
      ])
      return { success: true }
    } catch (err) {
      console.error('Erro ao carregar dashboard:', err)
      return { success: false, error: 'Erro ao carregar dashboard' }
    } finally {
      loading.value = false
    }
  }

  // Mudar mês/ano selecionado
  async function changeMonth(month, year) {
    selectedMonth.value = month
    selectedYear.value = year

    loading.value = true
    try {
      await Promise.all([
        fetchSummary(month, year),
        fetchExpensesByCategory(month, year),
        fetchIncomesByCategory(month, year)
      ])
    } finally {
      loading.value = false
    }
  }

  return {
    summary,
    expensesByCategory,
    incomesByCategory,
    monthlyFlow,
    recentTransactions,
    loading,
    error,
    selectedMonth,
    selectedYear,
    fetchSummary,
    fetchExpensesByCategory,
    fetchIncomesByCategory,
    fetchMonthlyFlow,
    fetchRecentTransactions,
    fetchAll,
    changeMonth
  }
})
