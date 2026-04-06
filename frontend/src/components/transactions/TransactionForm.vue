<script setup>
import { ref, reactive, watch, computed } from 'vue'
import Dialog from '@/components/ui/Dialog.vue'
import Input from '@/components/ui/Input.vue'
import Select from '@/components/ui/Select.vue'
import DatePicker from '@/components/ui/DatePicker.vue'
import { useTransactionsStore } from '@/stores/transactions'
import { useCategoriesStore } from '@/stores/categories'
import { useToast } from '@/composables/useToast'

const props = defineProps({
  open: Boolean,
  transaction: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['update:open', 'close', 'success'])

const transactionsStore = useTransactionsStore()
const categoriesStore = useCategoriesStore()
const toast = useToast()

const form = reactive({
  type: 'expense',
  description: '',
  amount: '',
  date: new Date().toISOString().split('T')[0],
  categoryId: '',
  notes: ''
})

const errors = reactive({
  type: '',
  description: '',
  amount: '',
  date: '',
  categoryId: ''
})

const loading = ref(false)

// Opções de tipo
const typeOptions = [
  { value: 'expense', label: 'Saída' },
  { value: 'income', label: 'Entrada' }
]

// Categorias filtradas por tipo
const filteredCategories = computed(() => {
  const categories = form.type === 'income'
    ? categoriesStore.incomeCategories
    : categoriesStore.expenseCategories

  return categories.map(c => ({
    value: c.id,
    label: `${c.icon} ${c.name}`
  }))
})

// Formatar valor monetário
const formatCurrency = (value) => {
  if (!value) return ''
  const number = parseFloat(value.toString().replace(/[^\d]/g, '')) / 100
  return number.toLocaleString('pt-BR', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  })
}

// Converter valor formatado para centavos
const parseCurrency = (value) => {
  if (!value) return 0
  const cleanValue = value.toString().replace(/[^\d]/g, '')
  return parseInt(cleanValue) || 0
}

// Handler para input de valor
const handleAmountInput = (event) => {
  const value = event.target.value.replace(/[^\d]/g, '')
  form.amount = formatCurrency(value)
}

// Resetar formulário
const resetForm = () => {
  form.type = 'expense'
  form.description = ''
  form.amount = ''
  form.date = new Date().toISOString().split('T')[0]
  form.categoryId = ''
  form.notes = ''

  errors.type = ''
  errors.description = ''
  errors.amount = ''
  errors.date = ''
  errors.categoryId = ''
}

// Preencher formulário com dados da transação (modo edição)
watch(() => props.transaction, (transaction) => {
  if (transaction) {
    form.type = transaction.type
    form.description = transaction.description
    form.amount = formatCurrency(transaction.amount)
    form.date = transaction.date.split('T')[0]
    form.categoryId = transaction.categoryId
    form.notes = transaction.notes || ''
  } else {
    resetForm()
  }
}, { immediate: true })

// Limpar categoria quando mudar o tipo (se não for compatível)
watch(() => form.type, () => {
  const currentCategory = categoriesStore.getCategoryById(form.categoryId)
  if (currentCategory && currentCategory.type !== 'both' && currentCategory.type !== form.type) {
    form.categoryId = ''
  }
})

// Validar formulário
const validate = () => {
  let isValid = true

  errors.type = ''
  errors.description = ''
  errors.amount = ''
  errors.date = ''
  errors.categoryId = ''

  if (!form.type) {
    errors.type = 'Tipo é obrigatório'
    isValid = false
  }

  if (!form.description.trim()) {
    errors.description = 'Descrição é obrigatória'
    isValid = false
  }

  const amountInCents = parseCurrency(form.amount)
  if (amountInCents <= 0) {
    errors.amount = 'Valor deve ser maior que zero'
    isValid = false
  }

  if (!form.date) {
    errors.date = 'Data é obrigatória'
    isValid = false
  }

  if (!form.categoryId) {
    errors.categoryId = 'Categoria é obrigatória'
    isValid = false
  }

  return isValid
}

// Salvar transação
const handleSubmit = async () => {
  if (!validate()) {
    toast.error('Por favor, preencha todos os campos obrigatórios')
    return
  }

  loading.value = true

  try {
    const transactionData = {
      type: form.type,
      description: form.description.trim(),
      amount: parseCurrency(form.amount),
      date: form.date,
      categoryId: parseInt(form.categoryId),
      notes: form.notes.trim()
    }

    let result

    if (props.transaction) {
      result = await transactionsStore.updateTransaction(props.transaction.id, transactionData)
    } else {
      result = await transactionsStore.createTransaction(transactionData)
    }

    if (result.success) {
      toast.success(
        props.transaction
          ? 'Transação atualizada com sucesso!'
          : 'Transação criada com sucesso!'
      )

      emit('success')
      handleClose()
    } else {
      toast.error(result.error || 'Erro ao salvar transação')
    }
  } catch (error) {
    console.error('Erro ao salvar transação:', error)
    toast.error('Erro ao salvar transação')
  } finally {
    loading.value = false
  }
}

// Fechar dialog
const handleClose = () => {
  resetForm()
  emit('update:open', false)
  emit('close')
}
</script>

<template>
  <Dialog
    :open="open"
    :title="transaction ? 'Editar Transação' : 'Nova Transação'"
    max-width="max-w-2xl"
    @update:open="handleClose"
  >
    <form @submit.prevent="handleSubmit" class="space-y-6">
      <!-- Tipo -->
      <div>
        <label class="block text-sm font-medium text-text-primary mb-2">Tipo</label>
        <div class="flex gap-4">
          <button
            type="button"
            @click="form.type = 'expense'"
            :class="[
              'flex-1 py-3 px-4 rounded-lg border-2 font-medium transition-all',
              form.type === 'expense'
                ? 'border-danger bg-danger/10 text-danger'
                : 'border-border-subtle bg-bg-secondary text-text-secondary hover:border-danger/50'
            ]"
            :disabled="loading"
          >
            <span class="mr-2">↓</span> Saída
          </button>
          <button
            type="button"
            @click="form.type = 'income'"
            :class="[
              'flex-1 py-3 px-4 rounded-lg border-2 font-medium transition-all',
              form.type === 'income'
                ? 'border-success bg-success/10 text-success'
                : 'border-border-subtle bg-bg-secondary text-text-secondary hover:border-success/50'
            ]"
            :disabled="loading"
          >
            <span class="mr-2">↑</span> Entrada
          </button>
        </div>
        <p v-if="errors.type" class="mt-1 text-sm text-danger">{{ errors.type }}</p>
      </div>

      <!-- Descrição -->
      <Input
        v-model="form.description"
        label="Descrição"
        placeholder="Ex: Almoço, Salário, Compras..."
        required
        :error="errors.description"
        :disabled="loading"
      />

      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- Valor -->
        <div>
          <label class="block text-sm font-medium text-text-primary mb-2">
            Valor <span class="text-danger">*</span>
          </label>
          <div class="relative">
            <span class="absolute left-3 top-1/2 -translate-y-1/2 text-text-secondary">R$</span>
            <input
              type="text"
              :value="form.amount"
              @input="handleAmountInput"
              placeholder="0,00"
              class="input pl-10"
              :class="{ 'border-danger': errors.amount }"
              :disabled="loading"
            />
          </div>
          <p v-if="errors.amount" class="mt-1 text-sm text-danger">{{ errors.amount }}</p>
        </div>

        <!-- Data -->
        <DatePicker
          v-model="form.date"
          label="Data"
          placeholder="Selecione uma data"
          required
          :error="errors.date"
          :disabled="loading"
        />
      </div>

      <!-- Categoria -->
      <Select
        v-model="form.categoryId"
        label="Categoria"
        :options="filteredCategories"
        placeholder="Selecione uma categoria"
        required
        :error="errors.categoryId"
        :disabled="loading"
      />

      <!-- Observações -->
      <div>
        <label class="block text-sm font-medium text-text-primary mb-2">
          Observações <span class="text-text-secondary">(opcional)</span>
        </label>
        <textarea
          v-model="form.notes"
          placeholder="Adicione notas ou detalhes..."
          rows="3"
          class="input resize-none"
          :disabled="loading"
        ></textarea>
      </div>
    </form>

    <template #footer>
      <div class="flex items-center justify-end gap-3">
        <button
          type="button"
          @click="handleClose"
          class="btn btn-ghost"
          :disabled="loading"
        >
          Cancelar
        </button>
        <button
          type="button"
          @click="handleSubmit"
          class="btn btn-primary"
          :disabled="loading"
        >
          <span v-if="loading">Salvando...</span>
          <span v-else>{{ transaction ? 'Atualizar' : 'Criar' }}</span>
        </button>
      </div>
    </template>
  </Dialog>
</template>
