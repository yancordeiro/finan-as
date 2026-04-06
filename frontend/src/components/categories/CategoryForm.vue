<script setup>
import { ref, reactive, watch } from 'vue'
import Dialog from '@/components/ui/Dialog.vue'
import Input from '@/components/ui/Input.vue'
import Select from '@/components/ui/Select.vue'
import ColorPicker from '@/components/ui/ColorPicker.vue'
import EmojiPicker from '@/components/ui/EmojiPicker.vue'
import { useCategoriesStore } from '@/stores/categories'
import { useToast } from '@/composables/useToast'

const props = defineProps({
  open: Boolean,
  category: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['update:open', 'close', 'success'])

const categoriesStore = useCategoriesStore()
const toast = useToast()

const form = reactive({
  name: '',
  icon: '📦',
  color: '#4f8ef7',
  type: 'expense'
})

const errors = reactive({
  name: '',
  icon: '',
  color: '',
  type: ''
})

const loading = ref(false)

// Opções de tipo
const typeOptions = [
  { value: 'income', label: 'Receita' },
  { value: 'expense', label: 'Despesa' },
  { value: 'both', label: 'Ambos' }
]

// Resetar formulário
const resetForm = () => {
  form.name = ''
  form.icon = '📦'
  form.color = '#4f8ef7'
  form.type = 'expense'

  errors.name = ''
  errors.icon = ''
  errors.color = ''
  errors.type = ''
}

// Preencher formulário com dados da categoria (modo edição)
watch(() => props.category, (category) => {
  if (category) {
    form.name = category.name
    form.icon = category.icon
    form.color = category.color
    form.type = category.type
  } else {
    resetForm()
  }
}, { immediate: true })

// Validar formulário
const validate = () => {
  let isValid = true

  errors.name = ''
  errors.icon = ''
  errors.color = ''
  errors.type = ''

  if (!form.name.trim()) {
    errors.name = 'Nome é obrigatório'
    isValid = false
  }

  if (!form.icon.trim()) {
    errors.icon = 'Ícone é obrigatório'
    isValid = false
  }

  if (!form.color.trim()) {
    errors.color = 'Cor é obrigatória'
    isValid = false
  }

  if (!form.type) {
    errors.type = 'Tipo é obrigatório'
    isValid = false
  }

  return isValid
}

// Salvar categoria
const handleSubmit = async () => {
  if (!validate()) {
    toast.error('Por favor, preencha todos os campos obrigatórios')
    return
  }

  loading.value = true

  try {
    let result

    if (props.category) {
      // Editar categoria existente
      result = await categoriesStore.updateCategory(props.category.id, form)
    } else {
      // Criar nova categoria
      result = await categoriesStore.createCategory(form)
    }

    if (result.success) {
      toast.success(
        props.category
          ? 'Categoria atualizada com sucesso!'
          : 'Categoria criada com sucesso!'
      )

      emit('success')
      handleClose()
    } else {
      toast.error(result.error || 'Erro ao salvar categoria')
    }
  } catch (error) {
    console.error('Erro ao salvar categoria:', error)
    toast.error('Erro ao salvar categoria')
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
    :title="category ? 'Editar Categoria' : 'Nova Categoria'"
    max-width="max-w-2xl"
    @update:open="handleClose"
  >
    <form @submit.prevent="handleSubmit" class="space-y-6">
      <!-- Nome -->
      <Input
        v-model="form.name"
        label="Nome"
        placeholder="Ex: Alimentação, Transporte..."
        required
        :error="errors.name"
        :disabled="loading"
      />

      <!-- Tipo -->
      <Select
        v-model="form.type"
        label="Tipo"
        :options="typeOptions"
        required
        :error="errors.type"
        :disabled="loading"
      />

      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- Ícone -->
        <EmojiPicker
          v-model="form.icon"
          label="Ícone"
          required
          :error="errors.icon"
        />

        <!-- Cor -->
        <ColorPicker
          v-model="form.color"
          label="Cor"
          required
          :error="errors.color"
        />
      </div>

      <!-- Preview -->
      <div class="p-4 bg-bg-primary border border-border-subtle rounded-lg">
        <p class="text-sm text-text-secondary mb-3">Prévia:</p>
        <div class="flex items-center gap-3">
          <div
            class="w-12 h-12 rounded-full flex items-center justify-center text-2xl"
            :style="{ backgroundColor: form.color + '33' }"
          >
            {{ form.icon }}
          </div>
          <div>
            <p class="font-semibold text-text-primary">{{ form.name || 'Nome da categoria' }}</p>
            <p class="text-sm text-text-secondary">
              {{ typeOptions.find(t => t.value === form.type)?.label }}
            </p>
          </div>
        </div>
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
          <span v-else>{{ category ? 'Atualizar' : 'Criar' }}</span>
        </button>
      </div>
    </template>
  </Dialog>
</template>
