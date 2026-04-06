<script setup>
import { ref } from 'vue'
import { LogOut, User, ChevronDown, Settings, Bell, Menu } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

const emit = defineEmits(['toggleSidebar'])

const authStore = useAuthStore()
const router = useRouter()

const showUserMenu = ref(false)

const handleLogout = async () => {
  await authStore.logout()
  router.push({ name: 'login' })
}

const toggleUserMenu = () => {
  showUserMenu.value = !showUserMenu.value
}

// Close menu when clicking outside
const closeMenu = () => {
  showUserMenu.value = false
}
</script>

<template>
  <header class="h-16 bg-bg-secondary/50 backdrop-blur-sm border-b border-border-subtle flex items-center justify-between px-4 sm:px-6">
    <!-- Left side - Breadcrumb/Title -->
    <div class="flex items-center gap-4">
      <button
        type="button"
        class="md:hidden p-2 rounded-xl bg-bg-elevated border border-border-subtle
               text-text-muted hover:text-text-primary hover:border-border-default
               transition-all duration-200"
        @click="emit('toggleSidebar')"
        aria-label="Abrir menu"
      >
        <Menu :size="18" />
      </button>
      <div>
        <h2 class="text-lg font-semibold text-text-primary leading-tight">Bem-vindo de volta</h2>
        <p class="hidden sm:block text-xs text-text-muted">
          {{ new Date().toLocaleDateString('pt-BR', { weekday: 'long', day: 'numeric', month: 'long' }) }}
        </p>
      </div>
    </div>

    <!-- Right side - User Menu -->
    <div class="flex items-center gap-3">
      <!-- Notifications (placeholder) -->
      <button class="relative p-2.5 rounded-xl bg-bg-elevated border border-border-subtle
                     text-text-muted hover:text-text-primary hover:border-border-default
                     transition-all duration-200">
        <Bell :size="18" />
        <span class="absolute top-1.5 right-1.5 w-2 h-2 bg-danger rounded-full"></span>
      </button>

      <!-- User dropdown -->
      <div class="relative" v-if="authStore.user">
        <button
          @click="toggleUserMenu"
          class="flex items-center gap-3 p-2 pr-3 rounded-xl bg-bg-elevated border border-border-subtle
                 hover:border-border-default transition-all duration-200"
        >
          <!-- Avatar -->
          <div class="w-8 h-8 rounded-lg bg-gradient-to-br from-primary to-accent flex items-center justify-center overflow-hidden">
            <img
              v-if="authStore.user.picture"
              :src="authStore.user.picture"
              :alt="authStore.user.name"
              class="w-full h-full object-cover"
            />
            <User v-else :size="16" class="text-white" />
          </div>

          <!-- Name -->
          <div class="hidden sm:block text-left">
            <p class="text-sm font-medium text-text-primary leading-tight">{{ authStore.user.name?.split(' ')[0] }}</p>
            <p class="text-[10px] text-text-muted leading-tight">{{ authStore.user.email?.split('@')[0] }}</p>
          </div>

          <ChevronDown
            :size="16"
            :class="[
              'text-text-muted transition-transform duration-200',
              showUserMenu ? 'rotate-180' : ''
            ]"
          />
        </button>

        <!-- Dropdown menu -->
        <Transition
          enter-active-class="transition duration-200 ease-out"
          enter-from-class="opacity-0 scale-95 -translate-y-2"
          enter-to-class="opacity-100 scale-100 translate-y-0"
          leave-active-class="transition duration-150 ease-in"
          leave-from-class="opacity-100 scale-100 translate-y-0"
          leave-to-class="opacity-0 scale-95 -translate-y-2"
        >
          <div
            v-if="showUserMenu"
            class="absolute right-0 top-full mt-2 w-56 bg-bg-card border border-border-subtle
                   rounded-xl shadow-elevated overflow-hidden z-50"
            @click="closeMenu"
          >
            <!-- User info header -->
            <div class="p-4 border-b border-border-subtle bg-bg-elevated/50">
              <p class="text-sm font-medium text-text-primary">{{ authStore.user.name }}</p>
              <p class="text-xs text-text-muted truncate">{{ authStore.user.email }}</p>
            </div>

            <!-- Menu items -->
            <div class="p-2">
              <button
                class="w-full flex items-center gap-3 px-3 py-2.5 rounded-lg
                       text-text-secondary hover:text-text-primary hover:bg-bg-elevated
                       transition-all duration-200 text-left"
              >
                <Settings :size="16" />
                <span class="text-sm">Configurações</span>
              </button>

              <div class="my-2 h-px bg-border-subtle" />

              <button
                @click="handleLogout"
                class="w-full flex items-center gap-3 px-3 py-2.5 rounded-lg
                       text-danger hover:bg-danger/10
                       transition-all duration-200 text-left"
              >
                <LogOut :size="16" />
                <span class="text-sm">Sair</span>
              </button>
            </div>
          </div>
        </Transition>
      </div>
    </div>
  </header>

  <!-- Overlay to close menu -->
  <div
    v-if="showUserMenu"
    class="fixed inset-0 z-40"
    @click="closeMenu"
  />
</template>
