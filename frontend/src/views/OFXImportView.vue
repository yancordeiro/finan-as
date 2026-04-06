<script setup>
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import OFXImportModal from '@/components/transactions/OFXImportModal.vue'
import { useCategoriesStore } from '@/stores/categories'

const router = useRouter()
const categoriesStore = useCategoriesStore()

const open = ref(true)

onMounted(async () => {
  if (!categoriesStore.categories?.length) {
    await categoriesStore.fetchCategories()
  }
})

const handleClose = () => {
  router.push('/dashboard')
}

const handleSuccess = () => {
  router.push('/transactions')
}
</script>

<template>
  <div class="min-h-[60vh]">
    <OFXImportModal
      v-model:open="open"
      @close="handleClose"
      @success="handleSuccess"
    />
  </div>
</template>
