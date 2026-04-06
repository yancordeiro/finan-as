# 📊 Dashboard Components

Componentes reutilizáveis para o dashboard principal da aplicação.

## Componentes

### 1. StatCard
Card de estatística com ícone, valor e indicador de variação.

**Props:**
```vue
{
  title: String,        // Título do card
  value: String,        // Valor principal a exibir
  change: Number,       // Variação percentual (0 = sem badge)
  changeLabel: String,  // Label da variação
  icon: Component,      // Componente de ícone
  color: String,        // Cor do card (primary, success, danger, warning, purple, indigo)
  loading: Boolean      // Estado de carregamento
}
```

**Exemplo:**
```vue
<StatCard
  title="Total de Entradas"
  value="R$ 5.000,00"
  :change="15"
  change-label="vs mês anterior"
  :icon="ArrowTrendingUpIcon"
  color="success"
/>
```

### 2. RecentTransactionsCard
Lista de transações recentes com navegação para a página completa.

**Props:**
```vue
{
  transactions: Array,  // Array de transações
  loading: Boolean      // Estado de carregamento
}
```

**Exemplo:**
```vue
<RecentTransactionsCard
  :transactions="recentTransactions"
  :loading="dashboardStore.loading"
/>
```

### 3. MonthSelector
Seletor de mês com navegação anterior/próxima.

**Props:**
```vue
{
  selectedMonth: Number,  // Mês selecionado (1-12)
  selectedYear: Number    // Ano selecionado
}
```

**Events:**
- `@previous` - Navegação para mês anterior
- `@next` - Navegação para próximo mês

**Exemplo:**
```vue
<MonthSelector
  :selected-month="dashboardStore.selectedMonth"
  :selected-year="dashboardStore.selectedYear"
  @previous="previousMonth"
  @next="nextMonth"
/>
```

## Features

- ✅ Loading states (skeleton)
- ✅ Responsive design
- ✅ Dark mode support
- ✅ Ícones coloridos
- ✅ Animações suaves
- ✅ Tipografia consistente

## Cores Disponíveis

### StatCard Colors
- **primary** - Azul (`text-blue-600 dark:text-blue-500`)
- **success** - Verde (`text-emerald-600 dark:text-emerald-500`)
- **danger** - Vermelho (`text-red-600 dark:text-red-500`)
- **warning** - Âmbar (`text-amber-600 dark:text-amber-500`)
- **purple** - Roxo (`text-purple-600 dark:text-purple-500`)
- **indigo** - Índigo (`text-indigo-600 dark:text-indigo-500`)

## Estrutura de Arquivos

```
components/dashboard/
├── StatCard.vue                 // Card de estatística
├── RecentTransactionsCard.vue   // Lista de transações recentes
├── MonthSelector.vue            // Seletor de mês
└── README.md                    // Esta documentação
```

## Uso no Dashboard

```vue
<script setup>
import StatCard from '@/components/dashboard/StatCard.vue'
import RecentTransactionsCard from '@/components/dashboard/RecentTransactionsCard.vue'
import MonthSelector from '@/components/dashboard/MonthSelector.vue'
</script>

<template>
  <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
    <div class="lg:col-span-8">
      <div class="grid grid-cols-4 gap-4">
        <StatCard v-bind="stats1" />
        <StatCard v-bind="stats2" />
        <StatCard v-bind="stats3" />
        <StatCard v-bind="stats4" />
      </div>
    </div>

    <div class="lg:col-span-12">
      <RecentTransactionsCard :transactions="recent" />
    </div>
  </div>
</template>
```

## Integração com Stores

Os componentes são projetados para trabalhar com Pinia stores:

```javascript
// No DashboardView.vue
const incomeStats = computed(() => ({
  title: 'Total de Entradas',
  value: formatCurrency(dashboardStore.summary.totalIncome),
  change: dashboardStore.summary.incomeChange,
  changeLabel: 'vs mês anterior',
  icon: ArrowTrendingUpIcon,
  color: 'success'
}))
```

---

Criado seguindo o padrão do exemplo-dash ✨
