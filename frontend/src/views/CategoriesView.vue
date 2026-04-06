<script setup>
import { ref, computed, onMounted } from 'vue'
import { Plus, Trash2, Lock, Edit2, FolderKanban, TrendingUp, TrendingDown } from 'lucide-vue-next'
import { useCategoriesStore } from '@/stores/categories'
import { useToast } from '@/composables/useToast'
import CategoryForm from '@/components/categories/CategoryForm.vue'

const categoriesStore = useCategoriesStore()
const toast = useToast()

const showForm = ref(false)
const editingCategory = ref(null)
const selectedFilter = ref('all')

// Opções de filtro
const filterOptions = [
  { value: 'all', label: 'Todas' },
  { value: 'income', label: 'Receitas' },
  { value: 'expense', label: 'Despesas' },
  { value: 'user', label: 'Minhas Categorias' },
  { value: 'global', label: 'Padrão do Sistema' }
]

// Categorias filtradas
const filteredCategories = computed(() => {
  let cats = categoriesStore.categories

  if (selectedFilter.value === 'income') {
    cats = cats.filter(c => c.type === 'income' || c.type === 'both')
  } else if (selectedFilter.value === 'expense') {
    cats = cats.filter(c => c.type === 'expense' || c.type === 'both')
  } else if (selectedFilter.value === 'user') {
    cats = cats.filter(c => c.userId !== null)
  } else if (selectedFilter.value === 'global') {
    cats = cats.filter(c => c.userId === null)
  }

  return cats
})

// Carregar categorias
onMounted(async () => {
  const result = await categoriesStore.fetchCategories()
  if (!result.success) {
    toast.error('Erro ao carregar categorias')
  }
})

// Abrir formulário para criar nova categoria
const handleCreate = () => {
  editingCategory.value = null
  showForm.value = true
}

// Abrir formulário para editar categoria
const handleEdit = (category) => {
  // Não permite editar categorias globais
  if (category.userId === null) {
    toast.warning('Categorias padrão do sistema não podem ser editadas')
    return
  }

  editingCategory.value = category
  showForm.value = true
}

// Deletar categoria
const handleDelete = async (category) => {
  // Não permite deletar categorias globais
  if (category.userId === null) {
    toast.warning('Categorias padrão do sistema não podem ser deletadas')
    return
  }

  if (!confirm(`Tem certeza que deseja deletar a categoria "${category.name}"?`)) {
    return
  }

  const result = await categoriesStore.deleteCategory(category.id)

  if (result.success) {
    toast.success('Categoria deletada com sucesso!')
  } else {
    toast.error(result.error || 'Erro ao deletar categoria')
  }
}

// Obter label do tipo
const getTypeLabel = (type) => {
  switch (type) {
    case 'income':
      return 'Receita'
    case 'expense':
      return 'Despesa'
    case 'both':
      return 'Ambos'
    default:
      return type
  }
}

// Callback de sucesso do form
const handleFormSuccess = () => {
  showForm.value = false
  editingCategory.value = null
}
</script>

<template>
  <div>
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4 mb-6">
      <h1 class="text-3xl font-bold">Categorias</h1>

      <div class="flex items-center gap-3">
        <!-- Filtros -->
        <select
          v-model="selectedFilter"
          class="input"
        >
          <option
            v-for="option in filterOptions"
            :key="option.value"
            :value="option.value"
          >
            {{ option.label }}
          </option>
        </select>

        <!-- Botão Nova Categoria -->
        <button
          @click="handleCreate"
          class="btn btn-primary flex items-center gap-2 whitespace-nowrap"
        >
          <Plus :size="20" />
          Nova Categoria
        </button>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="categoriesStore.loading && categoriesStore.categories.length === 0" class="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-4 gap-4">
      <div v-for="i in 8" :key="i" class="card p-4">
        <div class="flex items-center gap-3">
          <div class="skeleton w-12 h-12 rounded-full"></div>
          <div class="flex-1">
            <div class="skeleton h-4 w-24 mb-2"></div>
            <div class="skeleton h-3 w-16"></div>
          </div>
        </div>
      </div>
    </div>

    <!-- Lista de Categorias -->
    <div v-else-if="filteredCategories.length > 0" class="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-4 gap-4">
      <div
        v-for="category in filteredCategories"
        :key="category.id"
        class="card p-4 card-hover group relative"
      >
        <div
          @click="handleEdit(category)"
          class="flex items-center gap-3 cursor-pointer"
        >
          <div
            class="w-12 h-12 rounded-full flex items-center justify-center text-2xl flex-shrink-0"
            :style="{ backgroundColor: category.color + '33' }"
          >
            {{ category.icon }}
          </div>
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2">
              <p class="font-semibold truncate">{{ category.name }}</p>
              <Lock v-if="category.userId === null" :size="14" class="text-text-secondary flex-shrink-0" />
            </div>
            <p class="text-sm text-text-secondary">{{ getTypeLabel(category.type) }}</p>
          </div>
        </div>

        <!-- Botão Delete (apenas para categorias do usuário) -->
        <button
          v-if="category.userId !== null"
          @click.stop="handleDelete(category)"
          class="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity btn-ghost p-2 rounded-lg hover:bg-danger/10 hover:text-danger"
          title="Deletar categoria"
        >
          <Trash2 :size="16" />
        </button>
      </div>
    </div>

    <!-- Empty state -->
    <div v-else class="card p-12 text-center">
      <p class="text-text-secondary mb-4">Nenhuma categoria encontrada.</p>
      <button
        @click="selectedFilter = 'all'"
        class="btn btn-ghost"
      >
        Limpar filtros
      </button>
    </div>

    <!-- Formulário de Categoria -->
    <CategoryForm
      v-model:open="showForm"
      :category="editingCategory"
      @success="handleFormSuccess"
    />
  </div>
</template>
