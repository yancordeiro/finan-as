import { defineStore } from 'pinia'
import axios from 'axios'

export const useAIStore = defineStore('ai', {
  state: () => ({
    messages: [],
    isLoading: false,
    error: null,
    financialSummary: null,
    quickAnalyses: [],
  }),

  actions: {
    // Envia mensagem para a IA
    async sendMessage(message) {
      if (!message || message.trim() === '') {
        return
      }

      // Adiciona mensagem do usuário
      this.messages.push({
        role: 'user',
        content: message.trim(),
        timestamp: new Date(),
      })

      this.isLoading = true
      this.error = null

      try {
        const token = localStorage.getItem('accessToken')
        const response = await axios.post(
          '/api/ai/chat',
          { message: message.trim() },
          {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          }
        )

        // Adiciona resposta da IA
        this.messages.push({
          role: 'assistant',
          content: response.data.response,
          suggestions: response.data.suggestions || [],
          financialContext: response.data.financialContext,
          timestamp: new Date(),
        })

        // Atualiza contexto financeiro se disponível
        if (response.data.financialContext) {
          this.financialSummary = response.data.financialContext
        }
      } catch (error) {
        console.error('Erro ao enviar mensagem:', error)
        this.error = error.response?.data?.error || 'Erro ao processar mensagem'

        // Adiciona mensagem de erro
        this.messages.push({
          role: 'error',
          content: this.error,
          timestamp: new Date(),
        })
      } finally {
        this.isLoading = false
      }
    },

    // Carrega resumo financeiro
    async loadFinancialSummary() {
      try {
        const token = localStorage.getItem('accessToken')
        const response = await axios.get('/api/ai/summary', {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        })

        this.financialSummary = response.data
      } catch (error) {
        console.error('Erro ao carregar resumo financeiro:', error)
      }
    },

    // Carrega análises rápidas
    async loadQuickAnalyses() {
      try {
        const token = localStorage.getItem('accessToken')
        const response = await axios.get('/api/ai/quick-analysis', {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        })

        this.quickAnalyses = response.data.analyses || []
      } catch (error) {
        console.error('Erro ao carregar análises rápidas:', error)
      }
    },

    // Envia análise rápida
    async sendQuickAnalysis(analysis) {
      await this.sendMessage(analysis)
    },

    // Limpa conversa
    clearMessages() {
      this.messages = []
      this.error = null
    },

    // Remove mensagem específica
    removeMessage(index) {
      this.messages.splice(index, 1)
    },
  },

  getters: {
    // Retorna apenas mensagens válidas (sem erros)
    validMessages: (state) => {
      return state.messages.filter((m) => m.role !== 'error')
    },

    // Conta mensagens
    messageCount: (state) => state.messages.length,

    // Verifica se tem mensagens
    hasMessages: (state) => state.messages.length > 0,

    // Última mensagem
    lastMessage: (state) => {
      return state.messages.length > 0
        ? state.messages[state.messages.length - 1]
        : null
    },

    // Contexto financeiro formatado
    formattedSummary: (state) => {
      if (!state.financialSummary) return null

      return {
        income: state.financialSummary.totalIncome.toFixed(2),
        expense: state.financialSummary.totalExpense.toFixed(2),
        balance: state.financialSummary.balance.toFixed(2),
        isPositive: state.financialSummary.balance >= 0,
        transactionCount: state.financialSummary.transactionCount,
        topCategories: Object.entries(
          state.financialSummary.categoryBreakdown || {}
        )
          .sort((a, b) => b[1].amount - a[1].amount)
          .slice(0, 5)
          .map(([name, stats]) => ({
            name,
            amount: stats.amount.toFixed(2),
            percentage: stats.percentage.toFixed(1),
            count: stats.count,
          })),
      }
    },
  },
})
