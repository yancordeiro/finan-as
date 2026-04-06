<script setup>
import { ref } from 'vue'
import { Plus, CreditCard } from 'lucide-vue-next'

const accounts = ref([])
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-3xl font-bold">Contas Bancárias</h1>
      <button class="btn btn-primary flex items-center gap-2">
        <Plus :size="20" />
        Nova Conta
      </button>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <!-- Conta exemplo -->
      <div class="card p-6 card-hover">
        <div class="flex items-start justify-between mb-4">
          <div class="flex items-center gap-3">
            <div class="w-12 h-12 rounded-full bg-primary/20 flex items-center justify-center">
              <CreditCard :size="24" class="text-primary" />
            </div>
            <div>
              <p class="font-semibold">Nubank</p>
              <p class="text-sm text-text-secondary">Conta Corrente</p>
            </div>
          </div>
        </div>
        <div class="border-t border-border-subtle pt-4">
          <p class="text-sm text-text-secondary mb-1">Saldo</p>
          <p class="text-2xl font-bold text-primary">R$ 1.234,56</p>
        </div>
      </div>
    </div>
  </div>
</template>
