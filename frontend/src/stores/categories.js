import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useApi } from '@/composables/useApi'

export const useCategoriesStore = defineStore('categories', () => {
  const api = useApi()
  const categories = ref([])
  const loading = ref(false)
  const error = ref(null)

  // Computed: Categorias de receita
  const incomeCategories = computed(() => {
    return categories.value.filter(c => c.type === 'income' || c.type === 'both')
  })

  // Computed: Categorias de despesa
  const expenseCategories = computed(() => {
    return categories.value.filter(c => c.type === 'expense' || c.type === 'both')
  })

  // Computed: Categorias do usuário (não globais)
  const userCategories = computed(() => {
    return categories.value.filter(c => c.userId !== null)
  })

  // Computed: Categorias globais (padrão)
  const globalCategories = computed(() => {
    return categories.value.filter(c => c.userId === null)
  })

  // Buscar todas as categorias
  async function fetchCategories() {
    loading.value = true
    error.value = null

    try {
      const response = await api.get('/api/categories')
      categories.value = response.data.data
      return { success: true }
    } catch (err) {
      console.error('Erro ao buscar categorias:', err)
      error.value = err.response?.data?.error || 'Erro ao carregar categorias'
      return { success: false, error: error.value }
    } finally {
      loading.value = false
    }
  }

  // Criar nova categoria
  async function createCategory(categoryData) {
    loading.value = true
    error.value = null

    try {
      const response = await api.post('/api/categories', categoryData)
      const newCategory = response.data.data

      categories.value.push(newCategory)

      return { success: true, data: newCategory }
    } catch (err) {
      console.error('Erro ao criar categoria:', err)
      error.value = err.response?.data?.error || 'Erro ao criar categoria'
      return { success: false, error: error.value }
    } finally {
      loading.value = false
    }
  }

  // Atualizar categoria
  async function updateCategory(categoryId, categoryData) {
    loading.value = true
    error.value = null

    try {
      const response = await api.put(`/api/categories/${categoryId}`, categoryData)
      const updatedCategory = response.data.data

      const index = categories.value.findIndex(c => c.id === categoryId)
      if (index !== -1) {
        categories.value[index] = updatedCategory
      }

      return { success: true, data: updatedCategory }
    } catch (err) {
      console.error('Erro ao atualizar categoria:', err)
      error.value = err.response?.data?.error || 'Erro ao atualizar categoria'
      return { success: false, error: error.value }
    } finally {
      loading.value = false
    }
  }

  // Deletar categoria
  async function deleteCategory(categoryId) {
    loading.value = true
    error.value = null

    try {
      await api.delete(`/api/categories/${categoryId}`)

      const index = categories.value.findIndex(c => c.id === categoryId)
      if (index !== -1) {
        categories.value.splice(index, 1)
      }

      return { success: true }
    } catch (err) {
      console.error('Erro ao deletar categoria:', err)
      error.value = err.response?.data?.error || 'Erro ao deletar categoria'
      return { success: false, error: error.value }
    } finally {
      loading.value = false
    }
  }

  // Buscar categoria por ID
  function getCategoryById(categoryId) {
    return categories.value.find(c => c.id === categoryId)
  }

  return {
    categories,
    loading,
    error,
    incomeCategories,
    expenseCategories,
    userCategories,
    globalCategories,
    fetchCategories,
    createCategory,
    updateCategory,
    deleteCategory,
    getCategoryById,
  }
})
