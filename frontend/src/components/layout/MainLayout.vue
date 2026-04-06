<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import {
  LayoutDashboard,
  ArrowRightLeft,
  FolderKanban,
  FileUp,
  Menu,
  Sun,
  Moon,
  User,
  LogOut,
  ChevronDown,
  BrainCircuit,
} from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const sidebarOpen = ref(true)
const isDark = ref(false)
const userMenuOpen = ref(false)

const menuItems = [
  {
    title: 'Dashboard',
    icon: LayoutDashboard,
    to: '/dashboard',
    color: 'text-blue-600 dark:text-blue-500'
  },
  {
    title: 'Transações',
    icon: ArrowRightLeft,
    to: '/transactions',
    color: 'text-emerald-600 dark:text-emerald-500'
  },
  {
    title: 'Categorias',
    icon: FolderKanban,
    to: '/categories',
    color: 'text-purple-600 dark:text-purple-500'
  },
  {
    title: 'Assistente IA',
    icon: BrainCircuit,
    to: '/ai-assistant',
    color: 'text-pink-600 dark:text-pink-500'
  },
  {
    title: 'Importar OFX',
    icon: FileUp,
    to: '/ofx/import',
    color: 'text-amber-600 dark:text-amber-500'
  },
]

const isActive = (path) => route.path === path

const toggleTheme = () => {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark')
}

const handleLogout = async () => {
  await authStore.logout()
  router.push({ name: 'login' })
}

const toggleUserMenu = () => {
  userMenuOpen.value = !userMenuOpen.value
}

// Fechar o menu quando clicar fora
const handleClickOutside = (e) => {
  const userMenu = document.querySelector('[data-user-menu]')
  if (userMenu && !userMenu.contains(e.target) && userMenuOpen.value) {
    userMenuOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<template>
  <div class="flex h-screen overflow-hidden bg-background">
    <!-- Sidebar -->
    <aside
      :class="[
        'hidden md:flex flex-col border-r bg-card transition-all duration-300',
        sidebarOpen ? 'w-64' : 'w-16'
      ]"
    >
      <!-- Logo -->
      <div class="flex h-14 items-center border-b px-4">
        <h1 v-if="sidebarOpen" class="text-lg font-semibold">Finanças</h1>
      </div>

      <!-- Menu -->
      <nav class="flex-1 space-y-1 p-2">
        <router-link
          v-for="item in menuItems"
          :key="item.to"
          :to="item.to"
          :class="[
            'flex items-center gap-3 rounded-lg px-3 py-2 text-sm transition-colors',
            isActive(item.to)
              ? 'bg-secondary text-secondary-foreground'
              : 'text-muted-foreground hover:bg-secondary/50 hover:text-foreground'
          ]"
        >
          <component
            :is="item.icon"
            :size="20"
            class="flex-shrink-0"
            :class="isActive(item.to) ? item.color : ''"
          />
          <span v-if="sidebarOpen">{{ item.title }}</span>
        </router-link>
      </nav>

      <!-- Footer com logo ou informação -->
      <div class="border-t p-4">
        <p v-if="sidebarOpen" class="text-xs text-muted-foreground text-center">
          © 2024 Finanças
        </p>
      </div>
    </aside>

    <!-- Main Content -->
    <div class="flex flex-1 flex-col overflow-hidden">
      <!-- Header -->
      <header class="flex h-14 items-center gap-4 border-b bg-card px-4 lg:px-6">
        <button
          @click="sidebarOpen = !sidebarOpen"
          class="hidden md:inline-flex items-center justify-center rounded-md text-sm font-medium hover:bg-accent hover:text-accent-foreground h-9 w-9"
        >
          <Menu :size="20" />
        </button>

        <h2 class="text-lg font-semibold">{{ route.meta.title || 'Dashboard' }}</h2>

        <div class="ml-auto flex items-center gap-2">
          <!-- Theme Toggle Button -->
          <button
            @click="toggleTheme"
            class="inline-flex items-center justify-center rounded-md text-sm font-medium hover:bg-accent hover:text-accent-foreground h-9 w-9 transition-colors"
            :title="isDark ? 'Modo Claro' : 'Modo Escuro'"
          >
            <Sun v-if="!isDark" :size="20" class="text-amber-600" />
            <Moon v-else :size="20" class="text-indigo-400" />
          </button>

          <!-- User Menu -->
          <div class="relative" data-user-menu>
            <button
              @click.stop="toggleUserMenu"
              class="flex items-center gap-2 rounded-md px-3 h-9 hover:bg-accent hover:text-accent-foreground transition-colors"
            >
              <div class="w-7 h-7 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center">
                <User :size="16" class="text-white" />
              </div>
              <span class="hidden md:inline text-sm font-medium">
                {{ authStore.user?.name?.split(' ')[0] || 'Usuário' }}
              </span>
              <ChevronDown :size="16" class="hidden md:inline" />
            </button>

            <!-- Dropdown Menu -->
            <Transition
              enter-active-class="transition ease-out duration-100"
              enter-from-class="transform opacity-0 scale-95"
              enter-to-class="transform opacity-100 scale-100"
              leave-active-class="transition ease-in duration-75"
              leave-from-class="transform opacity-100 scale-100"
              leave-to-class="transform opacity-0 scale-95"
            >
              <div
                v-if="userMenuOpen"
                @click.stop
                class="absolute right-0 mt-2 w-56 rounded-lg bg-card border shadow-lg z-50"
              >
                <!-- User Info -->
                <div class="px-4 py-3 border-b">
                  <p class="text-sm font-medium">{{ authStore.user?.name || 'Usuário' }}</p>
                  <p class="text-xs text-muted-foreground truncate">{{ authStore.user?.email || '' }}</p>
                </div>

                <!-- Menu Items -->
                <div class="py-1">
                  <button
                    @click="handleLogout"
                    class="flex w-full items-center gap-3 px-4 py-2 text-sm text-red-600 dark:text-red-500 hover:bg-secondary/50 transition-colors"
                  >
                    <LogOut :size="16" />
                    <span>Sair</span>
                  </button>
                </div>
              </div>
            </Transition>
          </div>
        </div>
      </header>

      <!-- Page Content -->
      <main class="flex-1 overflow-y-auto p-4 lg:p-6">
        <RouterView />
      </main>
    </div>
  </div>
</template>
