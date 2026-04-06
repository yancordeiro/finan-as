<script setup>
import { ArrowUpRight, ArrowDownRight } from 'lucide-vue-next'

const transactions = [
  {
    id: 1,
    type: 'income',
    description: 'Paypal Transfer',
    date: '24 Jan 2024',
    amount: '+$2,500',
    status: 'completed'
  },
  {
    id: 2,
    type: 'expense',
    description: 'Office Supplies',
    date: '23 Jan 2024',
    amount: '-$450',
    status: 'completed'
  },
  {
    id: 3,
    type: 'income',
    description: 'Client Payment',
    date: '22 Jan 2024',
    amount: '+$5,200',
    status: 'pending'
  },
  {
    id: 4,
    type: 'expense',
    description: 'Marketing Campaign',
    date: '21 Jan 2024',
    amount: '-$1,800',
    status: 'completed'
  },
  {
    id: 5,
    type: 'income',
    description: 'Product Sales',
    date: '20 Jan 2024',
    amount: '+$3,400',
    status: 'completed'
  }
]
</script>

<template>
  <div class="card p-6">
    <div class="flex items-center justify-between mb-6">
      <h5 class="text-lg font-semibold">Recent Transactions</h5>
      <button class="text-sm text-blue-600 dark:text-blue-500 hover:underline">
        View All
      </button>
    </div>

    <div class="overflow-x-auto">
      <table class="w-full">
        <thead>
          <tr class="border-b">
            <th class="text-left py-3 px-2 text-sm font-medium text-muted-foreground">Description</th>
            <th class="text-left py-3 px-2 text-sm font-medium text-muted-foreground">Date</th>
            <th class="text-right py-3 px-2 text-sm font-medium text-muted-foreground">Amount</th>
            <th class="text-right py-3 px-2 text-sm font-medium text-muted-foreground">Status</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="transaction in transactions"
            :key="transaction.id"
            class="border-b last:border-b-0 hover:bg-muted/50 transition-colors"
          >
            <td class="py-4 px-2">
              <div class="flex items-center gap-3">
                <div
                  :class="[
                    'w-10 h-10 rounded-full flex items-center justify-center',
                    transaction.type === 'income'
                      ? 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-500'
                      : 'bg-red-500/10 text-red-600 dark:text-red-500'
                  ]"
                >
                  <ArrowUpRight v-if="transaction.type === 'income'" :size="20" />
                  <ArrowDownRight v-else :size="20" />
                </div>
                <span class="font-medium">{{ transaction.description }}</span>
              </div>
            </td>
            <td class="py-4 px-2 text-sm text-muted-foreground">
              {{ transaction.date }}
            </td>
            <td class="py-4 px-2 text-right">
              <span
                :class="[
                  'font-semibold',
                  transaction.type === 'income'
                    ? 'text-emerald-600 dark:text-emerald-500'
                    : 'text-red-600 dark:text-red-500'
                ]"
              >
                {{ transaction.amount }}
              </span>
            </td>
            <td class="py-4 px-2 text-right">
              <span
                :class="[
                  'inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium',
                  transaction.status === 'completed'
                    ? 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-500 border border-emerald-500/20'
                    : 'bg-amber-500/10 text-amber-600 dark:text-amber-500 border border-amber-500/20'
                ]"
              >
                {{ transaction.status === 'completed' ? 'Completed' : 'Pending' }}
              </span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
