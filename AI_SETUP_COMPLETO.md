# 🤖 Setup Completo - IA Local com Ollama

## ✅ O que foi implementado

Implementação completa de um assistente financeiro com IA rodando localmente via Docker!

### Backend (Go):
- ✅ `internal/models/ai_models.go` - Modelos de dados para IA
- ✅ `internal/services/ai_service.go` - Serviço que conecta com Ollama e analisa dados financeiros
- ✅ `internal/handlers/ai_handler.go` - Endpoints HTTP
- ✅ Rotas adicionadas em `internal/routes/routes.go`:
  - `POST /api/ai/chat` - Conversar com IA
  - `GET /api/ai/summary` - Resumo financeiro
  - `GET /api/ai/quick-analysis` - Análises rápidas

### Frontend (Vue 3):
- ✅ `stores/ai.js` - Store Pinia para gerenciar estado do chat
- ✅ `components/ai/AIChat.vue` - Componente de chat completo
- ✅ `views/AIAssistant.vue` - Página principal
- ✅ Rota adicionada: `/ai-assistant`
- ✅ Item de menu na sidebar: "Assistente IA"

### Docker:
- ✅ Ollama adicionado ao `docker-compose.yml`
- ✅ Scripts de setup criados
- ✅ Variáveis de ambiente configuradas

---

## 🚀 Como Usar (Passo a Passo)

### 1. Instalar Dependência do Frontend

Antes de iniciar, instale a biblioteca markdown:

```bash
cd frontend
npm install marked
cd ..
```

### 2. Atualizar arquivo .env

Copie as variáveis de Ollama do `.env.example` para seu `.env`:

```env
# Adicione ao seu .env
OLLAMA_BASE_URL=http://ollama:11434
OLLAMA_MODEL=llama3.1:8b
```

### 3. Iniciar Docker Compose

```bash
# Subir todos os serviços (postgres, ollama, backend, frontend)
docker-compose up -d

# Ver logs
docker-compose logs -f
```

**⏱️ Aguarde ~30 segundos** para o Ollama iniciar completamente.

### 4. Baixar Modelo de IA

**Opção A - Automático (Recomendado):**
```bash
bash setup-ollama.sh
```

**Opção B - Manual:**
```bash
# Baixar modelo Llama 3.1 8B (4.7GB)
docker-compose exec ollama ollama pull llama3.1:8b

# Verificar se baixou
docker-compose exec ollama ollama list
```

**Modelos disponíveis:**
- `llama3.1:8b` - Recomendado (4.7GB) - Melhor custo/benefício
- `phi3:medium` - Mais leve (7.9GB) - Mais rápido
- `mistral:7b` - Alternativa (4.1GB)
- `llama3.1:70b` - Avançado (40GB) - Requer GPU/muita RAM

### 5. Testar Ollama

```bash
# Teste simples
docker-compose exec ollama ollama run llama3.1:8b "Olá!"

# Teste com análise
docker-compose exec ollama ollama run llama3.1:8b "Gastei R$500 em delivery. Como economizar?"
```

### 6. Acessar o App

1. Abra o navegador: **http://localhost:3000**
2. Faça login com sua conta Google
3. Clique em **"Assistente IA"** no menu lateral

---

## 💬 Como Usar o Assistente IA

### Análises Rápidas (Clique nos botões):
- "Analise meus gastos do mês e identifique onde posso economizar"
- "Quais são minhas maiores despesas e como reduzi-las?"
- "Estou gastando muito? Sugira um plano de economia"
- "Por que meu saldo está negativo este mês?"

### Perguntas Personalizadas:
```
"Quanto gastei em delivery no último mês?"
"Estou gastando muito em alimentação?"
"Como posso economizar R$1000 por mês?"
"Crie uma meta de economia para viagem de R$5000"
"Analise meus gastos e me dê 3 dicas práticas"
```

### O que a IA sabe sobre você:
- ✅ Todas suas transações dos últimos 30 dias
- ✅ Receitas e despesas totais
- ✅ Gastos por categoria
- ✅ Saldo atual
- ✅ Padrões de consumo

---

## 🎨 Interface do Chat

```
┌─────────────────────────────────────────────┐
│ 🤖 Assistente Financeiro IA                 │
│ Análise inteligente com Ollama              │
├─────────────────────────────────────────────┤
│                                             │
│ [Análises Rápidas] ← Botões pré-definidos  │
│                                             │
│ ┌─────────────────────────────────────┐   │
│ │ Você: "Analise meus gastos"         │   │
│ └─────────────────────────────────────┘   │
│                                             │
│ ┌─────────────────────────────────────┐   │
│ │ 🤖 IA: "Analisando seus dados...   │   │
│ │                                     │   │
│ │ • Receitas: R$ 5.000,00            │   │
│ │ • Despesas: R$ 3.800,00            │   │
│ │ • Saldo: R$ 1.200,00               │   │
│ │                                     │   │
│ │ Principais gastos:                  │   │
│ │ • Alimentação: R$ 1.200 (32%)      │   │
│ │ • Delivery: R$ 600 (16%)           │   │
│ │                                     │   │
│ │ 💡 Sugestão:                       │   │
│ │ Reduzir delivery de R$600 para...  │   │
│ └─────────────────────────────────────┘   │
│                                             │
│ ┌─────────────────────────────────────┐   │
│ │ Digite sua mensagem...        [Enviar]│  │
│ └─────────────────────────────────────┘   │
└─────────────────────────────────────────────┘
```

---

## 🔧 Troubleshooting

### Erro: "Failed to process AI request"

**Causa:** Ollama não está rodando ou modelo não foi baixado

**Solução:**
```bash
# 1. Verificar se Ollama está rodando
docker-compose ps ollama

# 2. Ver logs do Ollama
docker-compose logs ollama

# 3. Verificar modelos
docker-compose exec ollama ollama list

# 4. Se vazio, baixar modelo
docker-compose exec ollama ollama pull llama3.1:8b

# 5. Reiniciar backend
docker-compose restart backend
```

### Erro: "marked is not defined"

**Causa:** Biblioteca `marked` não instalada

**Solução:**
```bash
cd frontend
npm install marked
docker-compose restart frontend
```

### Chat muito lento

**Opções:**

1. **Usar modelo menor:**
   ```bash
   docker-compose exec ollama ollama pull phi3:medium
   ```
   Atualizar `.env`:
   ```env
   OLLAMA_MODEL=phi3:medium
   ```

2. **Alocar mais RAM ao Docker:**
   - Docker Desktop → Settings → Resources → Memory: 8GB+

3. **Usar GPU** (se disponível):
   - Descomentar linhas de GPU no `docker-compose.yml`
   - Ver `OLLAMA_DOCKER.md` para instruções completas

### Erro 404 ao acessar /ai-assistant

**Causa:** Frontend precisa rebuild

**Solução:**
```bash
docker-compose restart frontend
# OU
docker-compose down
docker-compose up -d --build
```

---

## 📊 Arquitetura do Sistema

```
┌─────────────────────────────────────────────────────────┐
│ USUÁRIO                                                  │
│ http://localhost:3000/ai-assistant                      │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│ FRONTEND (Vue 3)                                         │
│ ┌─────────────────────────────────────────────────────┐ │
│ │ AIChat.vue                                          │ │
│ │ - Interface de chat                                 │ │
│ │ - Formatação markdown                               │ │
│ │ - Botões de análise rápida                          │ │
│ └──────────────────┬──────────────────────────────────┘ │
│                    │ POST /api/ai/chat                   │
│ ┌──────────────────▼──────────────────────────────────┐ │
│ │ ai.js (Pinia Store)                                 │ │
│ │ - Gerencia estado do chat                           │ │
│ │ - Envia mensagens via axios                         │ │
│ └─────────────────────────────────────────────────────┘ │
└────────────────────┬────────────────────────────────────┘
                     │ HTTP Request
                     ▼
┌─────────────────────────────────────────────────────────┐
│ BACKEND (Go/Fiber)                                       │
│ ┌─────────────────────────────────────────────────────┐ │
│ │ ai_handler.go                                       │ │
│ │ POST /api/ai/chat                                   │ │
│ └──────────────────┬──────────────────────────────────┘ │
│                    │                                     │
│ ┌──────────────────▼──────────────────────────────────┐ │
│ │ ai_service.go                                       │ │
│ │                                                      │ │
│ │ 1. buildFinancialContext(userID)                    │ │
│ │    └─> Query PostgreSQL                             │ │
│ │        ├─ Transações (últimos 30 dias)              │ │
│ │        ├─ Categorias do usuário                     │ │
│ │        └─ Calcula totais e percentuais              │ │
│ │                                                      │ │
│ │ 2. buildSystemPrompt(context)                       │ │
│ │    └─> Monta prompt com dados reais:                │ │
│ │        "Receitas: R$5000, Despesas: R$3800..."      │ │
│ │                                                      │ │
│ │ 3. callOllama(prompt, userMessage)                  │ │
│ │    └─> HTTP POST http://ollama:11434/api/chat       │ │
│ └─────────────────────────────────────────────────────┘ │
└────────────────────┬────────────────────────────────────┘
                     │ HTTP Request
                     ▼
┌─────────────────────────────────────────────────────────┐
│ OLLAMA (Container Docker)                                │
│ ┌─────────────────────────────────────────────────────┐ │
│ │ Modelo: Llama 3.1 8B                                │ │
│ │                                                      │ │
│ │ Input:                                              │ │
│ │ {                                                    │ │
│ │   "model": "llama3.1:8b",                          │ │
│ │   "messages": [                                     │ │
│ │     {                                               │ │
│ │       "role": "system",                             │ │
│ │       "content": "Você é assistente financeiro...   │ │
│ │                   Dados: Receitas R$5000..."        │ │
│ │     },                                              │ │
│ │     {                                               │ │
│ │       "role": "user",                               │ │
│ │       "content": "Analise meus gastos"              │ │
│ │     }                                               │ │
│ │   ]                                                 │ │
│ │ }                                                    │ │
│ │                                                      │ │
│ │ ⚙️ Processamento (LLM)...                          │ │
│ │                                                      │ │
│ │ Output:                                             │ │
│ │ "Analisando seus R$3800 em despesas, identifiquei  │ │
│ │  que você está gastando 32% em alimentação..."      │ │
│ └─────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────┘
                     │
                     ▼ (Retorna resposta)
              Backend → Frontend → Usuário
```

---

## 🎯 Próximos Passos (Melhorias Futuras)

### Já Implementado:
- ✅ Chat funcional com contexto financeiro
- ✅ Análise de gastos por categoria
- ✅ Sugestões básicas de economia
- ✅ Interface moderna e responsiva

### Possíveis Melhorias:
- [ ] **Tool Calling**: IA executar ações (criar metas, ajustar orçamento)
- [ ] **Histórico de conversas**: Salvar chats no BD
- [ ] **Análise de tendências**: Comparar meses
- [ ] **Alertas automáticos**: IA notifica gastos anormais
- [ ] **Export de relatórios**: IA gera PDFs com análises
- [ ] **Multi-modelos**: Permitir trocar modelo no frontend
- [ ] **Voice input**: Falar com IA via microfone

---

## 📚 Endpoints da API

### POST /api/ai/chat
Envia mensagem para IA e recebe análise contextualizada.

**Request:**
```json
{
  "message": "Analise meus gastos do mês"
}
```

**Response:**
```json
{
  "response": "Analisando seus dados dos últimos 30 dias:\n\n**Receitas**: R$ 5.000,00...",
  "suggestions": [
    {
      "type": "reduce_category",
      "title": "Reduzir gastos em Delivery",
      "description": "Você gastou R$ 600. Reduzir 20% economizaria R$ 120/mês",
      "action": "/transactions?category=Delivery"
    }
  ],
  "financialContext": {
    "totalIncome": 5000.00,
    "totalExpense": 3800.00,
    "balance": 1200.00,
    "transactionCount": 45,
    "period": "last_30_days",
    "categoryBreakdown": {
      "Alimentação": {
        "amount": 1200.00,
        "percentage": 31.58,
        "count": 12
      },
      "Delivery": {
        "amount": 600.00,
        "percentage": 15.79,
        "count": 8
      }
    }
  }
}
```

### GET /api/ai/summary
Retorna resumo financeiro sem análise de IA (mais rápido).

**Response:**
```json
{
  "totalIncome": 5000.00,
  "totalExpense": 3800.00,
  "balance": 1200.00,
  "transactionCount": 45,
  "period": "last_30_days",
  "categoryBreakdown": { ... }
}
```

### GET /api/ai/quick-analysis
Retorna lista de análises pré-definidas.

**Response:**
```json
{
  "userId": 123,
  "analyses": [
    "Analise meus gastos do mês e identifique onde posso economizar",
    "Quais são minhas maiores despesas e como reduzi-las?",
    "Estou gastando muito? Sugira um plano de economia"
  ]
}
```

---

## 🔒 Privacidade e Segurança

✅ **100% Local**
- Ollama roda no Docker na sua máquina
- Dados financeiros NUNCA saem do servidor
- Nenhuma comunicação com serviços externos

✅ **Autenticação**
- Endpoints protegidos com JWT
- Cada usuário vê apenas seus próprios dados

✅ **Isolamento**
- Docker network isolado
- Ollama não tem acesso à internet (opcional)

---

## 💰 Custo

**Total: R$ 0,00**

- Ollama: Grátis e open source
- Llama 3.1: Grátis e open source
- PostgreSQL: Grátis e open source
- Sem cobranças por token/uso
- Sem limites de requisições

**Único custo:** Eletricidade do PC/servidor rodando 😄

---

## 🎉 Pronto!

Seu assistente financeiro com IA está completo e funcionando!

**Comandos úteis:**

```bash
# Ver todos os serviços
docker-compose ps

# Logs em tempo real
docker-compose logs -f

# Reiniciar tudo
docker-compose restart

# Parar tudo
docker-compose down

# Ver modelos Ollama
docker-compose exec ollama ollama list

# Testar Ollama
docker-compose exec ollama ollama run llama3.1:8b "Hello!"
```

**Acesso rápido:**
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- Ollama API: http://localhost:11434

---

Para dúvidas, consulte:
- `OLLAMA_DOCKER.md` - Detalhes sobre Ollama no Docker
- `README.md` - Documentação geral do projeto
- Logs: `docker-compose logs ollama` ou `docker-compose logs backend`

**Desenvolvido com ❤️ usando Go, Vue.js, Ollama e Llama 3.1**
