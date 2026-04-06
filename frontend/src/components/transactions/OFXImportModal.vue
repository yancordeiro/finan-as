<script setup>
import { ref } from 'vue'
import Dialog from '@/components/ui/Dialog.vue'
import { useApi } from '@/composables/useApi'
import { useToast } from '@/composables/useToast'
import { Upload, FileText, CheckCircle, AlertCircle, Loader2 } from 'lucide-vue-next'

const props = defineProps({
  open: Boolean
})

const emit = defineEmits(['update:open', 'close', 'success'])

const api = useApi()
const toast = useToast()

// Estados
const step = ref(1) // 1: Upload e Preview, 2: Resultado
const loading = ref(false)
const selectedFile = ref(null)
const previewTransactions = ref([])
const importResult = ref(null)
const dragOver = ref(false)
const fileInputRef = ref(null)

// Abrir seletor de arquivo
const openFileSelector = () => {
  fileInputRef.value?.click()
}

// Handlers de drag and drop
const handleDragOver = (e) => {
  e.preventDefault()
  dragOver.value = true
}

const handleDragLeave = () => {
  dragOver.value = false
}

const handleDrop = (e) => {
  e.preventDefault()
  dragOver.value = false

  const files = e.dataTransfer.files
  if (files.length > 0) {
    handleFileSelect(files[0])
  }
}

// Selecionar arquivo
const handleFileInput = (e) => {
  const file = e.target.files[0]
  if (file) {
    handleFileSelect(file)
  }
}

const handleFileSelect = (file) => {
  if (!file.name.toLowerCase().endsWith('.ofx')) {
    toast.error('Por favor, selecione um arquivo .ofx')
    return
  }

  selectedFile.value = file
}

// Importar transações (diretamente sem categoria)
const handleImport = async () => {
  if (!selectedFile.value) {
    toast.error('Selecione um arquivo OFX')
    return
  }

  loading.value = true

  try {
    const formData = new FormData()
    formData.append('file', selectedFile.value)
    // Não enviamos categoryId - permitindo importação sem categoria

    console.log('Importando arquivo:', selectedFile.value.name)

    const response = await api.post('/api/ofx/import', formData)

    console.log('Resultado:', response.data)

    importResult.value = response.data.data
    step.value = 2

    if (importResult.value.imported > 0) {
      emit('success')
      toast.success(`${importResult.value.imported} transações importadas com sucesso!`)
    }
  } catch (err) {
    console.error('Erro ao importar:', err)
    console.error('Response:', err.response)
    toast.error(err.response?.data?.error || 'Erro ao importar transacoes')
  } finally {
    loading.value = false
  }
}

// Resetar e fechar
const handleClose = () => {
  step.value = 1
  selectedFile.value = null
  importResult.value = null
  loading.value = false
  dragOver.value = false

  emit('update:open', false)
  emit('close')
}
</script>

<template>
  <Dialog
    :open="open"
    title="Importar OFX"
    max-width="max-w-2xl"
    @update:open="handleClose"
  >
    <!-- Step 1: Upload do arquivo -->
    <div v-if="step === 1" class="space-y-6">
      <p class="text-text-secondary text-sm">
        Importe transacoes de um arquivo OFX do seu banco.
      </p>

      <!-- Dropzone -->
      <div
        @dragover="handleDragOver"
        @dragleave="handleDragLeave"
        @drop="handleDrop"
        :class="[
          'border-2 border-dashed rounded-xl p-8 text-center transition-all cursor-pointer',
          dragOver ? 'border-primary bg-primary/5' : 'border-border-subtle hover:border-primary/50',
          selectedFile ? 'bg-success/5 border-success' : ''
        ]"
        @click="openFileSelector"
      >
        <input
          ref="fileInputRef"
          type="file"
          accept=".ofx,.OFX"
          class="hidden"
          @change="handleFileInput"
        />

        <div v-if="selectedFile" class="space-y-3">
          <div class="w-16 h-16 mx-auto rounded-full bg-success/10 flex items-center justify-center">
            <FileText class="w-8 h-8 text-success" />
          </div>
          <div>
            <p class="font-medium text-text-primary">{{ selectedFile.name }}</p>
            <p class="text-sm text-text-secondary">
              {{ (selectedFile.size / 1024).toFixed(1) }} KB
            </p>
          </div>
          <button
            type="button"
            class="text-sm text-primary hover:underline"
            @click.stop="selectedFile = null"
          >
            Trocar arquivo
          </button>
        </div>

        <div v-else class="space-y-3">
          <div class="w-16 h-16 mx-auto rounded-full bg-bg-elevated flex items-center justify-center">
            <Upload class="w-8 h-8 text-text-secondary" />
          </div>
          <div>
            <p class="font-medium text-text-primary">
              Arraste um arquivo ou clique para selecionar
            </p>
            <p class="text-sm text-text-secondary">
              Apenas arquivos .ofx
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- Step 2: Resultado -->
    <div v-else-if="step === 2" class="space-y-6">
      <div class="text-center py-6">
        <div
          :class="[
            'w-20 h-20 mx-auto rounded-full flex items-center justify-center mb-4',
            importResult?.imported > 0 ? 'bg-success/10' : 'bg-warning/10'
          ]"
        >
          <CheckCircle
            v-if="importResult?.imported > 0"
            class="w-10 h-10 text-success"
          />
          <AlertCircle
            v-else
            class="w-10 h-10 text-warning"
          />
        </div>

        <h3 class="text-xl font-semibold text-text-primary mb-2">
          Importacao Concluida
        </h3>

        <div class="space-y-1 text-sm">
          <p class="text-text-secondary">
            <span class="font-medium text-text-primary">{{ importResult?.totalRead }}</span>
            transacoes no arquivo
          </p>
          <p class="text-success">
            <span class="font-medium">{{ importResult?.imported }}</span>
            importadas com sucesso
          </p>
          <p v-if="importResult?.duplicates > 0" class="text-warning">
            <span class="font-medium">{{ importResult?.duplicates }}</span>
            duplicadas (ignoradas)
          </p>
          <p v-if="importResult?.skipped > 0" class="text-text-muted">
            <span class="font-medium">{{ importResult?.skipped }}</span>
            ignoradas por erro
          </p>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="flex items-center justify-end gap-3">
        <button
          type="button"
          @click="handleClose"
          class="btn btn-ghost"
          :disabled="loading"
        >
          {{ step === 2 ? 'Fechar' : 'Cancelar' }}
        </button>

        <button
          v-if="step === 1"
          type="button"
          @click="handleImport"
          class="btn btn-primary"
          :disabled="!selectedFile || loading"
        >
          <Loader2 v-if="loading" class="w-4 h-4 animate-spin mr-2" />
          <span>{{ loading ? 'Importando...' : 'Importar Arquivo' }}</span>
        </button>
      </div>
    </template>
  </Dialog>
</template>
