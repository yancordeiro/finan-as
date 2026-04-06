<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import {
  LayoutDashboard,
  ArrowRightLeft,
  FolderKanban,
  FileUp,
  Wallet,
  ChevronLeft,
  ChevronRight,
  X,
  BrainCircuit,
} from 'lucide-vue-next'

const props = defineProps({
  open: Boolean,
  mobile: Boolean,
})

const emit = defineEmits(['toggle'])

const route = useRoute()

const menuItems = [
  {
    name: 'Dashboard',
    icon: LayoutDashboard,
    path: '/dashboard',
    description: 'Visão geral'
  },
  {
    name: 'Transações',
    icon: ArrowRightLeft,
    path: '/transactions',
    description: 'Gerenciar'
  },
  {
    name: 'Categorias',
    icon: FolderKanban,
    path: '/categories',
    description: 'Organizar'
  },
  {
    name: 'Assistente IA',
    icon: BrainCircuit,
    path: '/ai-assistant',
    description: 'Análise inteligente'
  },
  {
    name: 'Importar OFX',
    icon: FileUp,
    path: '/ofx/import',
    description: 'Extratos'
  },
]

const isActive = (path) => {
  return route.path === path
}

const asideClasses = computed(() => {
  return [
    'fixed md:relative inset-y-0 left-0 z-50 md:z-auto flex flex-col bg-bg-sidebar border-r border-border-subtle',
    'transition-all duration-300 ease-in-out',
    'w-72',
    props.open ? 'translate-x-0' : '-translate-x-full',
    // Desktop behavior
    'md:translate-x-0',
    props.open ? 'md:w-64' : 'md:w-20',
  ]
})
</script>

<template>
  <aside
    :class="asideClasses"
  >
    <!-- Logo -->
    <div class="h-16 flex items-center px-4 border-b border-border-subtle">
      <div class="flex items-center gap-3 flex-1 min-w-0">
        <div class="w-10 h-10 rounded-xl bg-gradient-to-br from-primary to-accent flex items-center justify-center shadow-glow-sm">
          <Wallet class="w-5 h-5 text-white" />
        </div>
        <div v-if="open" class="overflow-hidden">
          <h1 class="text-lg font-bold text-gradient-blue whitespace-nowrap">Finanças</h1>
          <p class="text-[10px] text-text-muted -mt-0.5">Controle financeiro</p>
        </div>
      </div>

      <!-- Close button (mobile) -->
      <button
        v-if="open"
        type="button"
        class="md:hidden p-2 rounded-lg text-text-muted hover:text-text-primary hover:bg-bg-elevated transition-colors"
        @click="emit('toggle')"
        aria-label="Fechar menu"
      >
        <X :size="18" />
      </button>
    </div>

    <!-- Menu Items -->
    <nav class="flex-1 p-3 space-y-1 overflow-y-auto scrollbar-hide">
      <RouterLink
        v-for="item in menuItems"
        :key="item.path"
        :to="item.path"
        :class="[
          'group relative flex items-center gap-3 px-3 py-3 rounded-xl transition-all duration-200',
          isActive(item.path)
            ? 'bg-primary/10 text-primary'
            : 'text-text-secondary hover:bg-bg-elevated hover:text-text-primary',
        ]"
      >
        <!-- Active indicator -->
        <div
          v-if="isActive(item.path)"
          class="absolute left-0 top-1/2 -translate-y-1/2 w-1 h-8 bg-primary rounded-r-full"
        />

        <!-- Icon container -->
        <div
          :class="[
            'flex-shrink-0 w-10 h-10 rounded-lg flex items-center justify-center transition-all duration-200',
            isActive(item.path)
              ? 'bg-primary/20 text-primary'
              : 'bg-bg-elevated group-hover:bg-bg-tertiary text-text-muted group-hover:text-text-primary',
          ]"
        >
          <component :is="item.icon" :size="20" />
        </div>

        <!-- Text -->
        <div v-if="open" class="overflow-hidden">
          <span class="font-medium whitespace-nowrap">{{ item.name }}</span>
          <p class="text-[10px] text-text-muted">{{ item.description }}</p>
        </div>

        <!-- Tooltip when collapsed -->
        <div
          v-if="!open"
          class="absolute left-full ml-2 px-2 py-1 bg-bg-elevated border border-border-subtle
                 rounded-lg text-sm font-medium text-text-primary whitespace-nowrap
                 opacity-0 invisible group-hover:opacity-100 group-hover:visible
                 transition-all duration-200 z-50 shadow-elevated"
        >
          {{ item.name }}
        </div>
      </RouterLink>
    </nav>

    <!-- Toggle Button -->
    <button
      @click="emit('toggle')"
      class="hidden md:flex absolute -right-3 top-20 w-6 h-6 rounded-full bg-bg-elevated border border-border-subtle
             flex items-center justify-center text-text-muted hover:text-text-primary
             hover:bg-bg-tertiary hover:border-border-default transition-all duration-200 shadow-card"
    >
      <ChevronLeft v-if="open" :size="14" />
      <ChevronRight v-else :size="14" />
    </button>

    <!-- Bottom decoration -->
    <div class="p-4 border-t border-border-subtle">
      <div v-if="open" class="text-center">
        <p class="text-[10px] text-text-muted">Versão 1.0.0</p>
      </div>
      <div v-else class="flex justify-center">
        <div class="w-2 h-2 rounded-full bg-success animate-pulse" />
      </div>
    </div>
  </aside>
</template>
