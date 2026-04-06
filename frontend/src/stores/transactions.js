import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useApi } from '@/composables/useApi'

export const useTransactionsStore = defineStore('transactions', () => {
  const api = useApi()

  const transactions = ref([])
  const loading = ref(false)
  const error = ref(null)

  // Paginação
  const pagination = ref({
    page: 1,
    limit: 20,
    total: 0,
    totalPages: 0
  })

  // Filtros ativos
  const filters = ref({
    startDate: null,
    endDate: null,
    type: '',
    categoryId: null,
    search: ''
  })

  // Computed: Total de entradas do período
  const totalIncome = computed(() => {
    return transactions.value
      .filter(t => t.type === 'income')
      .reduce((sum, t) => sum + t.amount, 0)
  })

  // Computed: Total de saídas do período
  const totalExpense = computed(() => {
    return transactions.value
      .filter(t => t.type === 'expense')
      .reduce((sum, t) => sum + t.amount, 0)
  })

  // Buscar transações com filtros e paginação
  async function fetchTransactions(newFilters = {}) {
    loading.value = true
    error.value = null

    // Atualizar filtros
    filters.value = { ...filters.value, ...newFilters }

    try {
      // Montar query params
      const params = new URLSearchParams()

      if (filters.value.startDate) {
        params.append('startDate', filters.value.startDate)
      }
      if (filters.value.endDate) {
        params.append('endDate', filters.value.endDate)
      }
      if (filters.value.type) {
        params.append('type', filters.value.type)
      }
      if (filters.value.categoryId) {
        params.append('categoryId', filters.value.categoryId)
      }
      if (filters.value.search) {
        params.append('search', filters.value.search)
      }

      params.append('page', pagination.value.page)
      params.append('limit', pagination.value.limit)

      const response = await api.get(`/api/transactions?${params.toString()}`)
      const data = response.data.data

      transactions.value = data.transactions || []
      pagination.value = {
        page: data.page,
        limit: data.limit,
        total: data.total,
        totalPages: data.totalPages
      }

      return { success: true }
    } catch (err) {
      console.error('Erro ao buscar transações:', err)
      error.value = err.response?.data?.error || 'Erro ao carregar transações'
      return { success: false, error: error.value }
    } finally {
      loading.value = false
    }
  }

  // Criar nova transação
  async function createTransaction(transactionData) {
    loading.value = true
    error.value = null

    try {
      const response = await api.post('/api/transactions', transactionData)
      const newTransaction = response.data.data

      // Adicionar ao início da lista
      transactions.value.unshift(newTransaction)
      pagination.value.total++

      return { success: true, data: newTransaction }
    } catch (err) {
      console.error('Erro ao criar transação:', err)
      error.value = err.response?.data?.error || 'Erro ao criar transação'
      return { success: false, error: error.value }
    } finally {
      loading.value = false
    }
  }

  // Atualizar transação
  async function updateTransaction(transactionId, transactionData) {
    loading.value = true
    error.value = null

    try {
      const response = await api.put(`/api/transactions/${transactionId}`, transactionData)
      const updatedTransaction = response.data.data

      const index = transactions.value.findIndex(t => t.id === transactionId)
      if (index !== -1) {
        transactions.value[index] = updatedTransaction
      }

      return { success: true, data: updatedTransaction }
    } catch (err) {
      console.error('Erro ao atualizar transação:', err)
      error.value = err.response?.data?.error || 'Erro ao atualizar transação'
      return { success: false, error: error.value }
    } finally {
      loading.value = false
    }
  }

  // Excluir transação
  async function deleteTransaction(transactionId) {
    loading.value = true
    error.value = null

    try {
      await api.delete(`/api/transactions/${transactionId}`)

      const index = transactions.value.findIndex(t => t.id === transactionId)
      if (index !== -1) {
        transactions.value.splice(index, 1)
        pagination.value.total--
      }

      return { success: true }
    } catch (err) {
      console.error('Erro ao excluir transação:', err)
      error.value = err.response?.data?.error || 'Erro ao excluir transação'
      return { success: false, error: error.value }
    } finally {
      loading.value = false
    }
  }

  // Mudar página
  async function changePage(page) {
    pagination.value.page = page
    return fetchTransactions()
  }

  // Limpar filtros
  function clearFilters() {
    filters.value = {
      startDate: null,
      endDate: null,
      type: '',
      categoryId: null,
      search: ''
    }
    pagination.value.page = 1
  }

  // Exportar para CSV
  async function exportCSV() {
    try {
      const params = new URLSearchParams()

      if (filters.value.startDate) {
        params.append('startDate', filters.value.startDate)
      }
      if (filters.value.endDate) {
        params.append('endDate', filters.value.endDate)
      }
      if (filters.value.type) {
        params.append('type', filters.value.type)
      }
      if (filters.value.categoryId) {
        params.append('categoryId', filters.value.categoryId)
      }
      if (filters.value.search) {
        params.append('search', filters.value.search)
      }

      const response = await api.get(`/api/transactions/export?${params.toString()}`, {
        responseType: 'blob'
      })

      // Criar link de download
      const url = window.URL.createObjectURL(new Blob([response.data]))
      const link = document.createElement('a')
      link.href = url
      link.setAttribute('download', 'transacoes.csv')
      document.body.appendChild(link)
      link.click()
      link.remove()
      window.URL.revokeObjectURL(url)

      return { success: true }
    } catch (err) {
      console.error('Erro ao exportar transações:', err)
      return { success: false, error: 'Erro ao exportar transações' }
    }
  }

  return {
    transactions,
    loading,
    error,
    pagination,
    filters,
    totalIncome,
    totalExpense,
    fetchTransactions,
    createTransaction,
    updateTransaction,
    deleteTransaction,
    changePage,
    clearFilters,
    exportCSV
  }
})
