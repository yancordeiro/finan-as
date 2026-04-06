<script setup>
import { computed } from 'vue'
import { ArrowUpRight, ArrowDownRight, ArrowRight } from 'lucide-vue-next'
import { useRouter } from 'vue-router'

const props = defineProps({
  transactions: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const router = useRouter()

const formatCurrency = (valueInCents) => {
  return (valueInCents / 100).toLocaleString('pt-BR', {
    style: 'currency',
    currency: 'BRL'
  })
}

const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('pt-BR', { day: '2-digit', month: 'short' })
}

const goToTransactions = () => {
  router.push('/transactions')
}
</script>

<template>
  <div class="card p-6">
    <div class="flex items-center justify-between mb-6">
      <h5 class="text-lg font-semibold">Transações Recentes</h5>
      <button
        @click="goToTransactions"
        class="text-sm text-blue-600 dark:text-blue-500 hover:underline flex items-center gap-1"
      >
        Ver Todas
        <ArrowRight :size="16" />
      </button>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="space-y-4">
      <div v-for="i in 5" :key="i" class="flex items-center gap-4">
        <div class="skeleton h-10 w-10 rounded-full"></div>
        <div class="flex-1 space-y-2">
          <div class="skeleton h-4 w-32"></div>
          <div class="skeleton h-3 w-20"></div>
        </div>
        <div class="skeleton h-6 w-20"></div>
      </div>
    </div>

    <!-- Empty state -->
    <div v-else-if="transactions.length === 0" class="text-center py-12">
      <p class="text-muted-foreground">Nenhuma transação encontrada</p>
      <button @click="goToTransactions" class="btn btn-outline mt-4">
        Adicionar Transação
      </button>
    </div>

    <!-- Transactions list -->
    <div v-else class="space-y-1">
      <div
        v-for="transaction in transactions"
        :key="transaction.id"
        class="flex items-center gap-4 py-3 px-2 rounded-lg hover:bg-muted/50 transition-colors cursor-pointer"
        @click="goToTransactions"
      >
        <div
          :class="[
            'w-10 h-10 rounded-full flex items-center justify-center flex-shrink-0',
            transaction.type === 'income'
              ? 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-500'
              : 'bg-red-500/10 text-red-600 dark:text-red-500'
          ]"
        >
          <ArrowUpRight v-if="transaction.type === 'income'" :size="20" />
          <ArrowDownRight v-else :size="20" />
        </div>

        <div class="flex-1 min-w-0">
          <p class="font-medium truncate">{{ transaction.description }}</p>
          <div class="flex items-center gap-2 text-sm text-muted-foreground">
            <span>{{ formatDate(transaction.date) }}</span>
            <span v-if="transaction.category">•</span>
            <span v-if="transaction.category" class="truncate">
              {{ transaction.category.icon }} {{ transaction.category.name }}
            </span>
          </div>
        </div>

        <div class="text-right">
          <span
            :class="[
              'font-semibold',
              transaction.type === 'income'
                ? 'text-emerald-600 dark:text-emerald-500'
                : 'text-red-600 dark:text-red-500'
            ]"
          >
            {{ transaction.type === 'income' ? '+' : '-' }}
            {{ formatCurrency(transaction.amount) }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.skeleton {
  @apply relative overflow-hidden bg-muted rounded;
}

.skeleton::after {
  content: '';
  @apply absolute inset-0;
  background: linear-gradient(
    90deg,
    transparent,
    rgba(255, 255, 255, 0.08),
    transparent
  );
  animation: shimmer 2s infinite;
}

@keyframes shimmer {
  0% {
    transform: translateX(-100%);
  }
  100% {
    transform: translateX(100%);
  }
}
</style>
