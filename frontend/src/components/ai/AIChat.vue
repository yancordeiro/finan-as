<template>
  <div class="ai-chat-container">
    <!-- Header -->
    <div class="chat-header">
      <div class="header-content">
        <div class="header-icon">
          <svg xmlns="http://www.w3.org/2000/svg" class="icon" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
          </svg>
        </div>
        <div>
          <h2 class="header-title">Assistente Financeiro IA</h2>
          <p class="header-subtitle">Análise inteligente dos seus gastos com Ollama</p>
        </div>
      </div>
      <button v-if="aiStore.hasMessages" @click="clearChat" class="clear-btn">
        <svg xmlns="http://www.w3.org/2000/svg" class="icon-sm" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
        </svg>
        Limpar
      </button>
    </div>

    <!-- Quick Analysis Buttons -->
    <div v-if="!aiStore.hasMessages && quickAnalyses.length > 0" class="quick-analysis">
      <p class="quick-title">Análises Rápidas:</p>
      <div class="quick-buttons">
        <button
          v-for="(analysis, index) in quickAnalyses"
          :key="index"
          @click="sendQuickAnalysis(analysis)"
          class="quick-btn"
        >
          {{ analysis }}
        </button>
      </div>
    </div>

    <!-- Messages Container -->
    <div ref="messagesContainer" class="messages-container">
      <div v-if="aiStore.messages.length === 0" class="empty-state">
        <svg xmlns="http://www.w3.org/2000/svg" class="empty-icon" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z" />
        </svg>
        <p class="empty-title">Comece uma conversa</p>
        <p class="empty-subtitle">Pergunte sobre seus gastos, metas ou peça sugestões de economia</p>
      </div>

      <div
        v-for="(message, index) in aiStore.messages"
        :key="index"
        :class="['message', `message-${message.role}`]"
      >
        <div class="message-avatar">
          <span v-if="message.role === 'user'" class="avatar-user">Você</span>
          <span v-else-if="message.role === 'assistant'" class="avatar-ai">🤖</span>
          <span v-else class="avatar-error">⚠️</span>
        </div>
        <div class="message-content">
          <div class="message-text" v-html="formatMessage(message.content)"></div>

          <!-- Suggestions -->
          <div v-if="message.suggestions && message.suggestions.length > 0" class="suggestions">
            <p class="suggestions-title">Sugestões:</p>
            <div
              v-for="(suggestion, sIndex) in message.suggestions"
              :key="sIndex"
              class="suggestion-card"
            >
              <h4 class="suggestion-title">{{ suggestion.title }}</h4>
              <p class="suggestion-desc">{{ suggestion.description }}</p>
            </div>
          </div>

          <span class="message-time">{{ formatTime(message.timestamp) }}</span>
        </div>
      </div>

      <!-- Loading indicator -->
      <div v-if="aiStore.isLoading" class="message message-assistant">
        <div class="message-avatar">
          <span class="avatar-ai">🤖</span>
        </div>
        <div class="message-content">
          <div class="typing-indicator">
            <span></span>
            <span></span>
            <span></span>
          </div>
        </div>
      </div>
    </div>

    <!-- Input Area -->
    <div class="input-area">
      <form @submit.prevent="sendMessage" class="input-form">
        <input
          v-model="userInput"
          type="text"
          placeholder="Pergunte sobre seus gastos, metas ou peça conselhos..."
          class="input-field"
          :disabled="aiStore.isLoading"
          maxlength="1000"
        />
        <button
          type="submit"
          :disabled="!userInput.trim() || aiStore.isLoading"
          class="send-btn"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="icon-sm" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8" />
          </svg>
          Enviar
        </button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick, watch } from 'vue'
import { useAIStore } from '@/stores/ai'
import { marked } from 'marked'

const aiStore = useAIStore()
const userInput = ref('')
const messagesContainer = ref(null)
const quickAnalyses = ref([])

onMounted(async () => {
  await aiStore.loadQuickAnalyses()
  quickAnalyses.value = aiStore.quickAnalyses
})

// Auto-scroll quando novas mensagens chegam
watch(() => aiStore.messages.length, () => {
  nextTick(() => {
    scrollToBottom()
  })
})

const sendMessage = async () => {
  if (!userInput.value.trim()) return

  const message = userInput.value
  userInput.value = ''

  await aiStore.sendMessage(message)
}

const sendQuickAnalysis = async (analysis) => {
  await aiStore.sendQuickAnalysis(analysis)
}

const clearChat = () => {
  if (confirm('Deseja limpar toda a conversa?')) {
    aiStore.clearMessages()
  }
}

const formatMessage = (content) => {
  // Converte markdown para HTML
  return marked.parse(content)
}

const formatTime = (timestamp) => {
  if (!timestamp) return ''
  const date = new Date(timestamp)
  return date.toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit' })
}

const scrollToBottom = () => {
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}
</script>

<style scoped>
.ai-chat-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 0;
  background: #1a1a24;
  border-radius: 12px;
  overflow: hidden;
}

.chat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  background: linear-gradient(135deg, #4f8ef7 0%, #a259ff 100%);
  border-bottom: 1px solid #2a2a3a;
}

.header-content {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.header-icon {
  width: 48px;
  height: 48px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.icon {
  width: 24px;
  height: 24px;
  color: white;
}

.icon-sm {
  width: 20px;
  height: 20px;
}

.header-title {
  font-size: 1.25rem;
  font-weight: 700;
  color: white;
  margin: 0;
}

.header-subtitle {
  font-size: 0.875rem;
  color: rgba(255, 255, 255, 0.8);
  margin: 0.25rem 0 0 0;
}

.clear-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  color: white;
  font-size: 0.875rem;
  cursor: pointer;
  transition: all 0.2s;
}

.clear-btn:hover {
  background: rgba(255, 255, 255, 0.2);
}

.quick-analysis {
  padding: 1rem;
  background: #13131c;
  border-bottom: 1px solid #2a2a3a;
}

.quick-title {
  font-size: 0.875rem;
  color: #7a7a9a;
  margin-bottom: 0.75rem;
}

.quick-buttons {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.quick-btn {
  padding: 0.5rem 1rem;
  background: #1a1a24;
  border: 1px solid #2a2a3a;
  border-radius: 8px;
  color: #e8e8f0;
  font-size: 0.875rem;
  cursor: pointer;
  transition: all 0.2s;
}

.quick-btn:hover {
  background: #4f8ef7;
  border-color: #4f8ef7;
  color: white;
}

.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 1.5rem;
  background: #0f0f13;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #7a7a9a;
  text-align: center;
}

.empty-icon {
  width: 64px;
  height: 64px;
  margin-bottom: 1rem;
  opacity: 0.3;
}

.empty-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: #e8e8f0;
  margin-bottom: 0.5rem;
}

.empty-subtitle {
  font-size: 0.875rem;
  max-width: 400px;
}

.message {
  display: flex;
  gap: 1rem;
  margin-bottom: 1.5rem;
  animation: slideIn 0.3s ease;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.message-avatar {
  flex-shrink: 0;
}

.avatar-user,
.avatar-ai,
.avatar-error {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  font-size: 0.75rem;
  font-weight: 600;
}

.avatar-user {
  background: #4f8ef7;
  color: white;
}

.avatar-ai {
  background: #a259ff;
  font-size: 1.25rem;
}

.avatar-error {
  background: #ff4d4d;
  font-size: 1.25rem;
}

.message-content {
  flex: 1;
  min-width: 0;
  max-width: calc(100% - 56px);
}

.message-text {
  background: #1a1a24;
  border: 1px solid #2a2a3a;
  border-radius: 12px;
  padding: 1rem;
  color: #e8e8f0;
  line-height: 1.6;
  overflow-wrap: break-word;
  word-break: break-word;
}

.message-text :deep(strong) {
  color: #4f8ef7;
}

.message-text :deep(ul),
.message-text :deep(ol) {
  margin: 0.5rem 0;
  padding-left: 1.5rem;
}

.message-text :deep(p) {
  margin: 0.5rem 0;
}

.message-text :deep(p:first-child) {
  margin-top: 0;
}

.message-text :deep(p:last-child) {
  margin-bottom: 0;
}

.message-text :deep(code) {
  background: #0f0f13;
  padding: 0.125rem 0.375rem;
  border-radius: 4px;
  font-size: 0.875em;
}

.message-text :deep(pre) {
  background: #0f0f13;
  padding: 1rem;
  border-radius: 8px;
  overflow-x: auto;
  margin: 0.75rem 0;
}

.message-text :deep(pre code) {
  background: none;
  padding: 0;
}

.message-error .message-text {
  background: rgba(255, 77, 77, 0.1);
  border-color: #ff4d4d;
  color: #ff4d4d;
}

.message-time {
  display: inline-block;
  margin-top: 0.5rem;
  font-size: 0.75rem;
  color: #7a7a9a;
}

.suggestions {
  margin-top: 1rem;
}

.suggestions-title {
  font-size: 0.875rem;
  color: #7a7a9a;
  margin-bottom: 0.5rem;
}

.suggestion-card {
  background: rgba(79, 142, 247, 0.1);
  border: 1px solid rgba(79, 142, 247, 0.3);
  border-radius: 8px;
  padding: 0.75rem;
  margin-bottom: 0.5rem;
}

.suggestion-title {
  font-size: 0.875rem;
  font-weight: 600;
  color: #4f8ef7;
  margin: 0 0 0.25rem 0;
}

.suggestion-desc {
  font-size: 0.875rem;
  color: #e8e8f0;
  margin: 0;
}

.typing-indicator {
  display: flex;
  gap: 0.5rem;
  padding: 1rem;
}

.typing-indicator span {
  width: 8px;
  height: 8px;
  background: #7a7a9a;
  border-radius: 50%;
  animation: typing 1.4s infinite;
}

.typing-indicator span:nth-child(2) {
  animation-delay: 0.2s;
}

.typing-indicator span:nth-child(3) {
  animation-delay: 0.4s;
}

@keyframes typing {
  0%, 60%, 100% {
    opacity: 0.3;
    transform: scale(0.8);
  }
  30% {
    opacity: 1;
    transform: scale(1);
  }
}

.input-area {
  padding: 1rem;
  background: #1a1a24;
  border-top: 1px solid #2a2a3a;
}

.input-form {
  display: flex;
  gap: 0.75rem;
}

.input-field {
  flex: 1;
  padding: 0.875rem 1rem;
  background: #0f0f13;
  border: 1px solid #2a2a3a;
  border-radius: 8px;
  color: #e8e8f0;
  font-size: 0.875rem;
  outline: none;
  transition: all 0.2s;
}

.input-field:focus {
  border-color: #4f8ef7;
  box-shadow: 0 0 0 3px rgba(79, 142, 247, 0.1);
}

.input-field:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.send-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.875rem 1.5rem;
  background: linear-gradient(135deg, #4f8ef7 0%, #a259ff 100%);
  border: none;
  border-radius: 8px;
  color: white;
  font-weight: 600;
  font-size: 0.875rem;
  cursor: pointer;
  transition: all 0.2s;
}

.send-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(79, 142, 247, 0.4);
}

.send-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
