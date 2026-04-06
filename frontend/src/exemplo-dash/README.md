# 📊 Example Dashboard

Dashboard de exemplo inspirado no template **Materio Vuetify** com design moderno usando **shadcn/ui**.

## 🎨 Componentes Criados

### 1. **CardStatisticsVertical**
Card de estatísticas vertical com:
- Ícone colorido em avatar
- Título e valor principal
- Indicador de variação percentual (↑/↓)
- Subtítulo descritivo
- Suporte a cores: primary, success, danger, warning, purple, indigo

**Props:**
```vue
{
  title: String,
  stats: String,
  change: Number,
  subtitle: String,
  icon: Component,
  color: String
}
```

### 2. **AnalyticsTransactions**
Grid responsivo com múltiplas estatísticas:
- 4 cards de métricas (Sales, Customers, Product, Revenue)
- Ícones coloridos (lucide-vue-next)
- Layout adaptativo: 4 cols (desktop) → 2 cols (tablet) → 1 col (mobile)

### 3. **WeeklyOverview**
Gráfico de barras semanal usando **Chart.js**:
- 7 barras (Sun-Sat)
- Barra destacada em azul (quarta-feira)
- Outras barras em cinza claro
- Eixo Y formatado (37k, 57k, etc.)
- Texto de performance
- Botão "Details" no rodapé

### 4. **AnalyticsAward**
Card de conquista com gradiente:
- Fundo gradiente azul
- Círculos decorativos
- Troféu ilustrativo
- Call-to-action button
- Texto de congratulações

### 5. **TotalEarning**
Card de ganhos totais:
- Valor principal centralizado
- Lista de earnings (Sales, Income, Expense)
- Ícones e badges de variação
- Indicador de % de mudança

### 6. **RecentTransactions**
Tabela de transações recentes:
- Colunas: Description, Date, Amount, Status
- Ícones coloridos (↗ verde / ↘ vermelho)
- Badges de status (Completed/Pending)
- Hover effect nas linhas
- Link "View All"

## 📐 Layout

O dashboard usa um sistema de grid responsivo (12 colunas):

```
┌──────────────────────────────────────┐
│  Award (4)    │  Transactions (8)    │
├──────────────────────────────────────┤
│ Weekly (4) │ Earning (4) │ Stats(4) │
├──────────────────────────────────────┤
│     Recent Transactions (12)         │
└──────────────────────────────────────┘
```

## 🎯 Como Usar

### 1. Acessar o Dashboard
Navegue para: `http://localhost:5173/exemplo-dash`

### 2. Reutilizar Componentes
```vue
<script setup>
import CardStatisticsVertical from '@/exemplo-dash/components/CardStatisticsVertical.vue'
import { PieChart } from 'lucide-vue-next'

const stats = {
  title: 'Total Profit',
  stats: '$25.6k',
  change: 42,
  subtitle: 'Weekly Project',
  icon: PieChart,
  color: 'purple'
}
</script>

<template>
  <CardStatisticsVertical v-bind="stats" />
</template>
```

## 🎨 Paleta de Cores

- **Primary (Azul)**: `text-blue-600 dark:text-blue-500`
- **Success (Verde)**: `text-emerald-600 dark:text-emerald-500`
- **Danger (Vermelho)**: `text-red-600 dark:text-red-500`
- **Warning (Âmbar)**: `text-amber-600 dark:text-amber-500`
- **Purple**: `text-purple-600 dark:text-purple-500`
- **Indigo**: `text-indigo-600 dark:text-indigo-500`
- **Cyan**: `text-cyan-600 dark:text-cyan-500`

## 📦 Dependências

- Vue 3
- Vue Router
- Chart.js + vue-chartjs
- lucide-vue-next (ícones)
- Tailwind CSS
- shadcn/ui (estilos)

## 🚀 Features

- ✅ Design moderno e limpo
- ✅ Totalmente responsivo
- ✅ Suporte a dark mode
- ✅ Ícones coloridos
- ✅ Gráficos interativos
- ✅ Componentes reutilizáveis
- ✅ Animações suaves
- ✅ Código bem documentado

## 📝 Estrutura de Arquivos

```
exemplo-dash/
├── components/
│   ├── AnalyticsAward.vue
│   ├── AnalyticsTransactions.vue
│   ├── CardStatisticsVertical.vue
│   ├── RecentTransactions.vue
│   ├── TotalEarning.vue
│   └── WeeklyOverview.vue
├── views/
│   └── ExampleDashboard.vue
└── README.md
```

## 💡 Próximos Passos

Você pode expandir este dashboard adicionando:
- [ ] Mais gráficos (pizza, linha, área)
- [ ] Filtros por data
- [ ] Exportação de dados
- [ ] Dashboards específicos por setor
- [ ] Integração com dados reais da API
- [ ] Modo de comparação (mês atual vs anterior)

---

Criado com 💙 inspirado no Materio Vuetify Template
