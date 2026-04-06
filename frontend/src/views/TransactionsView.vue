<script setup>
import { ref, onMounted, computed } from 'vue'
import { Plus, Download, Edit2, Trash2, ChevronLeft, ChevronRight, Search } from 'lucide-vue-next'
import { useTransactionsStore } from '@/stores/transactions'
import { useCategoriesStore } from '@/stores/categories'
import { useToast } from '@/composables/useToast'
import TransactionForm from '@/components/transactions/TransactionForm.vue'
import DatePicker from '@/components/ui/DatePicker.vue'
import Dialog from '@/components/ui/Dialog.vue'

const transactionsStore = useTransactionsStore()
const categoriesStore = useCategoriesStore()
const toast = useToast()

const showForm = ref(false)
const selectedTransaction = ref(null)
const showDeleteDialog = ref(false)
const transactionToDelete = ref(null)

const localFilters = ref({
  startDate: '',
  endDate: '',
  type: '',
  categoryId: '',
  search: ''
})

const formatCurrency = (valueInCents) => {
  return (valueInCents / 100).toLocaleString('pt-BR', {
    style: 'currency',
    currency: 'BRL'
  })
}

const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('pt-BR')
}

const categoryOptions = computed(() => {
  return categoriesStore.categories.map(c => ({
    value: c.id,
    label: `${c.icon} ${c.name}`
  }))
})

const fetchTransactions = async () => {
  const filters = {
    startDate: localFilters.value.startDate || null,
    endDate: localFilters.value.endDate || null,
    type: localFilters.value.type || '',
    categoryId: localFilters.value.categoryId || null,
    search: localFilters.value.search || ''
  }
  await transactionsStore.fetchTransactions(filters)
}

const applyFilters = () => {
  transactionsStore.pagination.page = 1
  fetchTransactions()
}

const clearFilters = () => {
  localFilters.value = {
    startDate: '',
    endDate: '',
    type: '',
    categoryId: '',
    search: ''
  }
  transactionsStore.clearFilters()
  fetchTransactions()
}

const openNewTransactionForm = () => {
  selectedTransaction.value = null
  showForm.value = true
}

const openEditTransactionForm = (transaction) => {
  selectedTransaction.value = transaction
  showForm.value = true
}

const confirmDelete = (transaction) => {
  transactionToDelete.value = transaction
  showDeleteDialog.value = true
}

const deleteTransaction = async () => {
  if (!transactionToDelete.value) return

  const result = await transactionsStore.deleteTransaction(transactionToDelete.value.id)

  if (result.success) {
    toast.success('Transação excluída com sucesso!')
  } else {
    toast.error(result.error || 'Erro ao excluir transação')
  }

  showDeleteDialog.value = false
  transactionToDelete.value = null
}

const exportCSV = async () => {
  const result = await transactionsStore.exportCSV()
  if (result.success) {
    toast.success('Exportação realizada com sucesso!')
  } else {
    toast.error(result.error || 'Erro ao exportar transações')
  }
}

const changePage = (page) => {
  if (page < 1 || page > transactionsStore.pagination.totalPages) return
  transactionsStore.changePage(page)
}

const onFormSuccess = () => {
  fetchTransactions()
}

onMounted(async () => {
  await categoriesStore.fetchCategories()
  await fetchTransactions()
})
</script>

<template>
  <div class="space-y-4">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h1 class="text-3xl font-bold tracking-tight">Transações</h1>
        <p class="text-muted-foreground">Gerencie suas receitas e despesas</p>
      </div>
      <div class="flex items-center gap-3">
        <button
          @click="exportCSV"
          class="btn btn-outline"
          :disabled="transactionsStore.transactions.length === 0"
        >
          <Download :size="16" class="text-emerald-600 dark:text-emerald-500" />
          Exportar
        </button>
        <button
          @click="openNewTransactionForm"
          class="btn btn-default"
        >
          <Plus :size="16" />
          Nova Transação
        </button>
      </div>
    </div>

    <!-- Filtros -->
    <div class="card p-6">
      <h3 class="text-sm font-medium mb-4">Filtros</h3>
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-4">
        <!-- Data inicial -->
        <DatePicker
          v-model="localFilters.startDate"
          label="Data inicial"
          placeholder="Data inicial"
        />

        <!-- Data final -->
        <DatePicker
          v-model="localFilters.endDate"
          label="Data final"
          placeholder="Data final"
        />

        <!-- Tipo -->
        <div class="space-y-2">
          <label class="text-sm font-medium">Tipo</label>
          <select v-model="localFilters.type" class="input">
            <option value="">Todos</option>
            <option value="income">Entradas</option>
            <option value="expense">Saídas</option>
          </select>
        </div>

        <!-- Categoria -->
        <div class="space-y-2">
          <label class="text-sm font-medium">Categoria</label>
          <select v-model="localFilters.categoryId" class="input">
            <option value="">Todas</option>
            <option v-for="cat in categoryOptions" :key="cat.value" :value="cat.value">
              {{ cat.label }}
            </option>
          </select>
        </div>

        <!-- Busca -->
        <div class="space-y-2">
          <label class="text-sm font-medium">Buscar</label>
          <div class="relative">
            <Search class="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground" :size="16" />
            <input
              type="text"
              v-model="localFilters.search"
              placeholder="Descrição..."
              class="input pl-10"
              @keyup.enter="applyFilters"
            />
          </div>
        </div>
      </div>

      <div class="flex items-center gap-3 mt-4">
        <button @click="applyFilters" class="btn btn-default">
          Filtrar
        </button>
        <button @click="clearFilters" class="btn btn-outline">
          Limpar
        </button>
      </div>
    </div>

    <!-- Tabela -->
    <div class="card">
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead>
            <tr class="border-b">
              <th class="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Data</th>
              <th class="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Descrição</th>
              <th class="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Categoria</th>
              <th class="h-12 px-4 text-right align-middle font-medium text-muted-foreground">Valor</th>
              <th class="h-12 px-4 text-right align-middle font-medium text-muted-foreground">Ações</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="transaction in transactionsStore.transactions"
              :key="transaction.id"
              class="border-b transition-colors hover:bg-muted/50"
            >
              <td class="p-4 align-middle text-sm text-muted-foreground">
                {{ formatDate(transaction.date) }}
              </td>
              <td class="p-4 align-middle">
                <div>
                  <div class="font-medium">{{ transaction.description }}</div>
                  <div v-if="transaction.notes" class="text-sm text-muted-foreground truncate max-w-xs">
                    {{ transaction.notes }}
                  </div>
                </div>
              </td>
              <td class="p-4 align-middle">
                <div class="flex items-center gap-2">
                  <div
                    class="flex h-8 w-8 items-center justify-center rounded-full text-sm"
                    :style="{ backgroundColor: (transaction.category?.color || '#94a3b8') + '20' }"
                  >
                    {{ transaction.category?.icon || '📦' }}
                  </div>
                  <span class="text-sm">{{ transaction.category?.name || 'Sem categoria' }}</span>
                </div>
              </td>
              <td class="p-4 align-middle text-right">
                <span
                  class="font-semibold"
                  :class="transaction.type === 'income' ? 'text-green-600 dark:text-green-500' : 'text-red-600 dark:text-red-500'"
                >
                  {{ transaction.type === 'income' ? '+' : '-' }}
                  {{ formatCurrency(transaction.amount) }}
                </span>
              </td>
              <td class="p-4 align-middle text-right">
                <div class="flex items-center justify-end gap-2">
                  <button
                    @click="openEditTransactionForm(transaction)"
                    class="inline-flex items-center justify-center rounded-md text-sm font-medium hover:bg-blue-500/10 text-blue-600 dark:text-blue-500 h-8 w-8 transition-colors"
                  >
                    <Edit2 :size="16" />
                  </button>
                  <button
                    @click="confirmDelete(transaction)"
                    class="inline-flex items-center justify-center rounded-md text-sm font-medium hover:bg-red-500/10 text-red-600 dark:text-red-500 h-8 w-8 transition-colors"
                  >
                    <Trash2 :size="16" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Paginação -->
      <div
        v-if="transactionsStore.pagination.totalPages > 1"
        class="flex items-center justify-between px-4 py-4 border-t"
      >
        <p class="text-sm text-muted-foreground">
          Mostrando {{ transactionsStore.transactions.length }} de {{ transactionsStore.pagination.total }} transações
        </p>
        <div class="flex items-center gap-2">
          <button
            @click="changePage(transactionsStore.pagination.page - 1)"
            :disabled="transactionsStore.pagination.page === 1"
            class="btn btn-outline h-8 w-8 p-0"
          >
            <ChevronLeft :size="16" />
          </button>

          <span class="text-sm">
            Página {{ transactionsStore.pagination.page }} de {{ transactionsStore.pagination.totalPages }}
          </span>

          <button
            @click="changePage(transactionsStore.pagination.page + 1)"
            :disabled="transactionsStore.pagination.page === transactionsStore.pagination.totalPages"
            class="btn btn-outline h-8 w-8 p-0"
          >
            <ChevronRight :size="16" />
          </button>
        </div>
      </div>
    </div>

    <!-- Formulário de Transação -->
    <TransactionForm
      v-model:open="showForm"
      :transaction="selectedTransaction"
      @success="onFormSuccess"
    />

    <!-- Dialog de Confirmação de Exclusão -->
    <Dialog
      :open="showDeleteDialog"
      title="Excluir Transação"
      max-width="max-w-md"
      @update:open="showDeleteDialog = false"
    >
      <p class="text-sm text-muted-foreground">
        Tem certeza que deseja excluir a transação
        <strong class="text-foreground">"{{ transactionToDelete?.description }}"</strong>?
      </p>
      <p class="text-sm text-muted-foreground mt-2">
        Esta ação não pode ser desfeita.
      </p>

      <template #footer>
        <div class="flex items-center justify-end gap-3">
          <button
            @click="showDeleteDialog = false"
            class="btn btn-outline"
          >
            Cancelar
          </button>
          <button
            @click="deleteTransaction"
            class="btn btn-destructive"
          >
            Excluir
          </button>
        </div>
      </template>
    </Dialog>
  </div>
</template>
