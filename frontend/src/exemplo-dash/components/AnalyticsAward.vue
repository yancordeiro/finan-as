<script setup>
import { Trophy } from 'lucide-vue-next'
</script>

<template>
  <div class="card p-6 h-full bg-gradient-to-br from-blue-500 to-blue-600 text-white relative overflow-hidden">
    <!-- Decorative circles -->
    <div class="absolute -top-10 -right-10 w-40 h-40 rounded-full bg-white/10"></div>
    <div class="absolute -bottom-10 -left-10 w-40 h-40 rounded-full bg-white/10"></div>

    <div class="relative z-10">
      <div class="flex items-start justify-between mb-6">
        <div>
          <h5 class="text-xl font-bold mb-1">Congratulations John! 🎉</h5>
          <p class="text-blue-100 text-sm">You have done 57.6% more sales today.</p>
        </div>
      </div>

      <div class="mb-6">
        <h2 class="text-4xl font-bold mb-1">$89k</h2>
        <p class="text-blue-100 text-sm">Total Sales</p>
      </div>

      <button class="bg-white text-blue-600 hover:bg-blue-50 font-medium px-4 py-2 rounded-lg transition-colors">
        View Sales
      </button>
    </div>

    <!-- Trophy illustration -->
    <div class="absolute bottom-0 right-0 p-6">
      <div class="w-24 h-24 rounded-full bg-white/20 flex items-center justify-center">
        <Trophy :size="48" class="text-white" />
      </div>
    </div>
  </div>
</template>
