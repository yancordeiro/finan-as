<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/composables/useToast'
import { useGoogleAuth } from '@/composables/useGoogleAuth'
import { Wallet, TrendingUp, PieChart, Shield } from 'lucide-vue-next'

const router = useRouter()
const authStore = useAuthStore()
const toast = useToast()
const { initGoogleAuth, renderGoogleButton } = useGoogleAuth()

const loading = ref(false)
const googleReady = ref(false)

// Features list
const features = [
  { icon: TrendingUp, text: 'Acompanhe suas finanças em tempo real' },
  { icon: PieChart, text: 'Visualize gastos por categoria' },
  { icon: Shield, text: 'Seus dados protegidos e seguros' },
]

// Callback do Google OAuth
const handleGoogleCallback = async (response) => {
  loading.value = true

  try {
    const result = await authStore.login(response.credential)

    if (result.success) {
      toast.success('Login realizado com sucesso!')
      router.push({ name: 'dashboard' })
    } else {
      toast.error(result.error || 'Erro ao fazer login')
    }
  } catch (error) {
    console.error('Erro no login:', error)
    toast.error('Erro ao fazer login. Tente novamente.')
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  const loaded = await initGoogleAuth()

  if (loaded) {
    googleReady.value = true
    setTimeout(() => {
      renderGoogleButton('google-signin-button', handleGoogleCallback)
    }, 100)
  } else {
    toast.error('Erro ao carregar autenticação do Google')
  }
})
</script>

<template>
  <div class="min-h-screen flex items-center justify-center relative overflow-hidden bg-bg-primary">
    <!-- Animated background -->
    <div class="absolute inset-0 overflow-hidden">
      <!-- Gradient orbs -->
      <div class="absolute top-1/4 -left-20 w-96 h-96 bg-primary/20 rounded-full blur-[128px] animate-pulse-slow" />
      <div class="absolute bottom-1/4 -right-20 w-96 h-96 bg-accent/20 rounded-full blur-[128px] animate-pulse-slow" style="animation-delay: 1s;" />
      <div class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] bg-secondary/10 rounded-full blur-[128px] animate-pulse-slow" style="animation-delay: 2s;" />

      <!-- Grid pattern -->
      <div class="absolute inset-0 bg-[linear-gradient(rgba(59,130,246,0.03)_1px,transparent_1px),linear-gradient(90deg,rgba(59,130,246,0.03)_1px,transparent_1px)] bg-[size:64px_64px]" />

      <!-- Noise texture -->
      <div class="absolute inset-0 bg-noise pointer-events-none" />
    </div>

    <!-- Content -->
    <div class="relative z-10 w-full max-w-md mx-4">
      <!-- Card -->
      <div class="bg-bg-card/80 backdrop-blur-xl border border-border-subtle rounded-2xl shadow-elevated p-8">
        <!-- Logo -->
        <div class="text-center mb-8">
          <div class="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-gradient-to-br from-primary to-accent mb-4 shadow-glow">
            <Wallet class="w-8 h-8 text-white" />
          </div>
          <h1 class="text-3xl font-bold text-gradient mb-2">Finanças</h1>
          <p class="text-text-secondary">Controle financeiro inteligente</p>
        </div>

        <!-- Divider -->
        <div class="divider mb-8" />

        <!-- Google Button -->
        <div class="space-y-4">
          <!-- Loading state -->
          <div v-if="!googleReady" class="space-y-4">
            <div class="skeleton h-12 w-full rounded-xl"></div>
            <p class="text-text-muted text-sm text-center animate-pulse">
              Carregando autenticação...
            </p>
          </div>

          <!-- Google Sign In Button -->
          <div v-else>
            <div id="google-signin-button" class="w-full flex justify-center"></div>

            <!-- Loading overlay -->
            <div v-if="loading" class="mt-6 text-center">
              <div class="inline-flex items-center gap-3 px-4 py-2 rounded-xl bg-primary/10 border border-primary/20">
                <div class="w-5 h-5 border-2 border-primary border-t-transparent rounded-full animate-spin" />
                <span class="text-sm text-primary">Autenticando...</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Divider -->
        <div class="divider my-8" />

        <!-- Features -->
        <div class="space-y-3">
          <div
            v-for="(feature, index) in features"
            :key="index"
            class="flex items-center gap-3 text-sm text-text-secondary"
          >
            <div class="w-8 h-8 rounded-lg bg-bg-elevated flex items-center justify-center flex-shrink-0">
              <component :is="feature.icon" :size="16" class="text-primary" />
            </div>
            <span>{{ feature.text }}</span>
          </div>
        </div>
      </div>

      <!-- Footer -->
      <p class="text-text-muted text-xs text-center mt-6">
        Ao entrar, você concorda com nossos
        <a href="#" class="text-primary hover:underline">Termos de Uso</a>
        e
        <a href="#" class="text-primary hover:underline">Política de Privacidade</a>
      </p>
    </div>
  </div>
</template>
