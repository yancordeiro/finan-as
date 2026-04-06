# 🚀 Setup Docker - Aplicação Financeira

## 📋 Pré-requisitos

- ✅ Docker Desktop instalado
- ✅ Conta Google (para login)

---

## ⚡ PASSO A PASSO RÁPIDO

### **1. Configure o Google OAuth** (5 minutos)

#### a) Criar credenciais:
1. Acesse: https://console.cloud.google.com/
2. Crie um **novo projeto** (ou selecione existente)
3. No menu lateral: **APIs e Serviços** → **Credenciais**
4. Clique em **Criar Credenciais** → **ID do cliente OAuth 2.0**
5. Se solicitado, configure a **Tela de consentimento OAuth** (básico)

#### b) Configurar o OAuth:
- **Tipo de aplicativo**: Aplicativo da Web
- **Nome**: Finanças App (ou qualquer nome)
- **Origens JavaScript autorizadas**:
  ```
  http://localhost:3000
  ```
- **URIs de redirecionamento autorizados**: (deixe vazio)

#### c) Copiar credenciais:
- Após criar, copie o **Client ID** (ex: `123456.apps.googleusercontent.com`)
- Copie o **Client Secret**

---

### **2. Criar arquivo `.env`**

Na **raiz do projeto** (`C:\repos\finanças`), crie um arquivo chamado `.env`:

```env
# JWT Secrets (MUDE PARA STRINGS ALEATÓRIAS)
JWT_SECRET=sua-string-aleatoria-aqui-minimo-32-caracteres
JWT_REFRESH_SECRET=outra-string-aleatoria-diferente-minimo-32-caracteres

# Google OAuth (COLE AQUI AS CREDENCIAIS DO PASSO 1)
GOOGLE_CLIENT_ID=seu-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=seu-client-secret-aqui

# URLs (não precisa alterar)
GOOGLE_REDIRECT_URI=http://localhost:3000/auth/callback
FRONTEND_URL=http://localhost:3000
VITE_API_URL=http://localhost:8080
```

#### 🔐 Para gerar JWT secrets fortes:

**Opção 1 - Online:**
Acesse: https://randomkeygen.com/ e copie uma "Fort Knox Password"

**Opção 2 - PowerShell (Windows):**
```powershell
-join ((48..57) + (65..90) + (97..122) | Get-Random -Count 32 | % {[char]$_})
```

**Opção 3 - Qualquer coisa:**
Use qualquer string aleatória de 32+ caracteres (ex: `minha-senha-super-secreta-12345678`)

---

### **3. Subir os containers**

No terminal, na pasta do projeto:

```bash
docker-compose up -d
```

**Aguarde 30-60 segundos** para tudo subir.

---

### **4. Acessar a aplicação**

Abra o navegador:
```
http://localhost:3000
```

Você verá a tela de login! 🎉

---

## ✅ Verificar se está funcionando

### Verificar containers:
```bash
docker-compose ps

# Deve mostrar 3 containers rodando:
# ✅ financas-postgres   (porta 5432)
# ✅ financas-backend    (porta 8080)
# ✅ financas-frontend   (porta 3000)
```

### Testar backend:
```bash
curl http://localhost:8080/health

# Resposta esperada:
# {"status":"ok","message":"API Financas rodando"}
```

### Ver logs:
```bash
# Todos os containers
docker-compose logs -f

# Apenas backend
docker-compose logs -f backend

# Apenas frontend
docker-compose logs -f frontend
```

---

## 🎯 Usar a aplicação

### 1. Fazer Login
- Clique em **"Entrar com Google"**
- Selecione sua conta Google
- Autorize o acesso
- Você será redirecionado para o **Dashboard**

### 2. Criar Categorias
- Menu lateral: **Categorias**
- Você verá 13 categorias padrão (🍔 Alimentação, 🚗 Transporte, etc.)
- Clique em **"Nova Categoria"**
- Escolha emoji, cor e tipo
- Clique em **"Criar"**

### 3. Funcionalidades disponíveis:
- ✅ **Dashboard** - Resumo financeiro
- ✅ **Categorias** - CRUD completo
- ✅ **Transações** - Em breve
- ✅ **Importar OFX** - Em breve

---

## 🛠️ Comandos Úteis

### Parar containers:
```bash
docker-compose down
```

### Resetar banco de dados (apaga tudo):
```bash
docker-compose down -v
docker-compose up -d
```

### Rebuild (após alterar código):
```bash
docker-compose up -d --build
```

### Reiniciar apenas um serviço:
```bash
docker-compose restart backend
# ou
docker-compose restart frontend
```

### Entrar no banco de dados:
```bash
docker-compose exec postgres psql -U financas -d financas_db

# Ver tabelas
\dt

# Ver categorias
SELECT name, icon, type FROM categories WHERE user_id IS NULL;

# Sair
\q
```

---

## 🐛 Problemas Comuns

### ❌ Erro: "JWT_SECRET é obrigatório"
**Causa**: Arquivo `.env` não existe ou está vazio
**Solução**: Crie o `.env` conforme passo 2

---

### ❌ Erro: "Token do Google inválido"
**Causa**: `GOOGLE_CLIENT_ID` incorreto no `.env`
**Solução**:
1. Verifique se copiou corretamente do Google Console
2. Confirme que adicionou `http://localhost:3000` nas origens autorizadas

---

### ❌ Porta 3000 ou 8080 já em uso
**Solução Windows:**
```bash
# Descobrir o que está usando a porta
netstat -ano | findstr :3000

# Matar o processo (substitua PID)
taskkill /PID numero_do_pid /F
```

---

### ❌ Botão do Google não aparece
**Causa**: Script do Google não carregou
**Solução**:
1. Abra DevTools (F12) → Console
2. Procure por erros
3. Aguarde 5 segundos e recarregue a página
4. Verifique se `GOOGLE_CLIENT_ID` está correto no `.env`

---

### ❌ Containers não sobem / Erro genérico
**Solução**:
```bash
# Ver logs detalhados
docker-compose logs

# Rebuild completo
docker-compose down
docker-compose up -d --build

# Ver se todos os containers subiram
docker-compose ps
```

---

### ❌ Categorias padrão não aparecem
**Causa**: Seed não executou
**Solução**:
```bash
# Verificar logs do backend
docker-compose logs backend | findstr "categorias padrão"

# Se não aparecer, resetar banco
docker-compose down -v
docker-compose up -d
```

---

## 📁 Estrutura do Projeto

```
finanças/
├── backend/              # API Go
│   ├── cmd/api/         # Entry point
│   ├── internal/        # Código interno
│   └── Dockerfile
│
├── frontend/            # Vue 3 SPA
│   ├── src/
│   └── Dockerfile
│
├── docker-compose.yml   # Orquestração
├── .env                 # Configurações (CRIAR!)
└── README.md
```

---

## 🎨 Funcionalidades Implementadas

| Funcionalidade | Status |
|----------------|--------|
| ✅ Login Google OAuth | 100% |
| ✅ CRUD Categorias | 100% |
| ✅ Dashboard básico | 100% |
| ⏳ CRUD Transações | 0% |
| ⏳ Gráficos | 0% |
| ⏳ Importação OFX | 0% |

---

## 📞 Ajuda

Se nada funcionar:

1. ✅ Confirme que o Docker Desktop está **rodando**
2. ✅ Confirme que o `.env` existe e tem todos os campos
3. ✅ Rode `docker-compose down -v` e `docker-compose up -d` (reset total)
4. ✅ Veja os logs: `docker-compose logs`

---

**Pronto! Sua aplicação está rodando!** 🚀

Acesse: http://localhost:3000
