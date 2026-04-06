# 🧪 Guia de Teste - Autenticação Google OAuth

## ✅ O que foi implementado

### Backend (100% completo)
- ✅ Service de autenticação com validação do Google
- ✅ Geração e validação de JWT (access + refresh tokens)
- ✅ Endpoints: `/api/auth/google`, `/api/auth/refresh`, `/api/auth/logout`
- ✅ Endpoint `/api/user/me` para obter dados do usuário autenticado
- ✅ Middleware de autenticação JWT
- ✅ Refresh tokens salvos no banco com expiração

### Frontend (100% completo)
- ✅ Integração com Google Identity Services
- ✅ Botão de login do Google (renderizado oficialmente)
- ✅ Pinia store de autenticação
- ✅ Auto-refresh de token em caso de 401
- ✅ Toast notifications para feedback
- ✅ Router guards para rotas protegidas
- ✅ Persistência de sessão (localStorage)

---

## 🚀 Como Testar

### Passo 1: Configurar Google OAuth 2.0

1. Acesse https://console.cloud.google.com/
2. Crie um novo projeto ou selecione um existente
3. Vá em "APIs e serviços" > "Credenciais"
4. Clique em "Criar credenciais" > "ID do cliente OAuth 2.0"
5. Configure a tela de consentimento OAuth se necessário
6. Em "Tipo de aplicativo", selecione "Aplicativo da Web"
7. Adicione as seguintes **Origens JavaScript autorizadas**:
   ```
   http://localhost:3000
   http://localhost:8080
   ```
8. Adicione os seguintes **URIs de redirecionamento autorizados**:
   ```
   http://localhost:3000/auth/callback
   ```
9. Copie o **Client ID** gerado

### Passo 2: Criar arquivo .env

Crie um arquivo `.env` na raiz do projeto:

```bash
# Na pasta raiz
cp .env.example .env
```

Edite o `.env` e adicione suas credenciais:

```env
# JWT Secrets (GERE SECRETS FORTES!)
JWT_SECRET=cole-aqui-um-secret-forte-minimo-32-caracteres
JWT_REFRESH_SECRET=cole-aqui-outro-secret-forte-diferente

# Google OAuth
GOOGLE_CLIENT_ID=seu-client-id-aqui.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=seu-client-secret-aqui

# URLs
FRONTEND_URL=http://localhost:3000
VITE_API_URL=http://localhost:8080
```

**Para gerar secrets fortes:**

```bash
# Linux/Mac
openssl rand -hex 32

# Windows PowerShell
[Convert]::ToBase64String((1..32 | ForEach-Object { Get-Random -Maximum 256 }))

# Ou use um site como: https://randomkeygen.com/
```

### Passo 3: Iniciar a aplicação com Docker

```bash
# Na raiz do projeto
docker-compose up -d

# Aguardar containers subirem (30-60 segundos)

# Verificar status
docker-compose ps

# Ver logs em tempo real
docker-compose logs -f
```

### Passo 4: Testar a autenticação

1. **Abra o navegador**: http://localhost:3000

2. **Tela de Login**: Você deve ver:
   - Logo "Finanças" com gradiente
   - Botão oficial do Google "Entrar com Google"
   - Mensagem sobre termos de uso

3. **Clique em "Entrar com Google"**:
   - Popup do Google deve abrir
   - Selecione sua conta Google
   - Autorize o acesso

4. **Após login bem-sucedido**:
   - Toast de sucesso aparece no canto superior direito
   - Você é redirecionado para `/dashboard`
   - Topbar exibe seu avatar e nome
   - Sidebar mostra o menu lateral

5. **Testar persistência**:
   - Recarregue a página (F5)
   - Você deve continuar logado
   - Feche a aba e reabra: ainda logado

6. **Testar logout**:
   - Clique no ícone de logout (🚪) no topbar
   - Você volta para a tela de login
   - Toast de confirmação aparece

---

## 🔍 Verificações Técnicas

### 1. Verificar Backend (API)

```bash
# Health check
curl http://localhost:8080/health

# Deve retornar:
# {"status":"ok","message":"API Financas rodando"}
```

### 2. Verificar Banco de Dados

```bash
# Entrar no container do PostgreSQL
docker-compose exec postgres psql -U financas -d financas_db

# Listar tabelas
\dt

# Verificar categorias padrão
SELECT id, name, icon, color, type FROM categories WHERE user_id IS NULL;

# Verificar usuários criados
SELECT id, email, name FROM users;

# Sair
\q
```

### 3. Verificar Logs

```bash
# Logs do backend
docker-compose logs backend

# Procurar por erros
docker-compose logs backend | grep -i error

# Logs do frontend
docker-compose logs frontend
```

### 4. Testar Endpoints Manualmente

**Obter um token de teste** (após fazer login na aplicação):

1. Abra DevTools (F12) > Application > Local Storage
2. Copie o valor de `accessToken`

**Testar endpoint protegido:**

```bash
# Substitua SEU_TOKEN_AQUI pelo token copiado
curl http://localhost:8080/api/user/me \
  -H "Authorization: Bearer SEU_TOKEN_AQUI"

# Deve retornar seus dados:
# {"success":true,"message":"Usuário encontrado","data":{"id":1,"email":"...","name":"...","picture":"..."}}
```

**Testar endpoint de refresh:**

```bash
curl http://localhost:8080/api/auth/refresh \
  -X POST \
  -H "Content-Type: application/json" \
  --cookie "refreshToken=SEU_REFRESH_TOKEN" \
  -c cookies.txt

# Deve retornar novos tokens
```

---

## 🐛 Troubleshooting

### Erro: "Token do Google inválido"

**Causa**: Client ID não configurado ou incorreto

**Solução**:
1. Verifique se `GOOGLE_CLIENT_ID` está no `.env`
2. Verifique se é o mesmo Client ID do Google Cloud Console
3. Verifique se `http://localhost:3000` está nas origens autorizadas

### Erro: "CORS error"

**Causa**: FRONTEND_URL incorreto no backend

**Solução**:
```env
# No .env, verifique:
FRONTEND_URL=http://localhost:3000
```

Reinicie o backend:
```bash
docker-compose restart backend
```

### Erro: "JWT_SECRET é obrigatório"

**Causa**: Secrets não configurados

**Solução**:
Gere secrets fortes e adicione ao `.env`:
```bash
openssl rand -hex 32
```

### Erro: "Erro ao conectar ao banco de dados"

**Causa**: PostgreSQL não iniciou

**Solução**:
```bash
# Verificar status
docker-compose ps

# Reiniciar PostgreSQL
docker-compose restart postgres

# Aguardar 10 segundos e testar
```

### Botão do Google não aparece

**Causa**: Google Identity Services não carregou

**Solução**:
1. Abra DevTools (F12) > Console
2. Procure por erros de rede
3. Verifique se o script do Google carregou:
   ```javascript
   console.log(window.google)
   // Deve retornar um objeto
   ```

### Token expira muito rápido

**Comportamento esperado**:
- Access token expira em **1 hora**
- Refresh token expira em **30 dias**
- Auto-refresh acontece automaticamente em caso de 401

Para testar auto-refresh, aguarde 1 hora e faça uma requisição à API.

---

## ✨ Próximos Passos

Após confirmar que a autenticação está funcionando perfeitamente:

1. ✅ **Sprint 2 completo!** Autenticação 100% funcional
2. 🚀 **Próximo**: Implementar CRUD de Categorias
3. 🚀 **Depois**: CRUD de Contas Bancárias
4. 🚀 **Depois**: CRUD de Transações
5. 🚀 **Depois**: Dashboard com gráficos
6. 🚀 **Depois**: Importação OFX

---

## 📞 Dúvidas?

Se encontrar algum problema:
1. Verifique os logs: `docker-compose logs -f`
2. Verifique o console do navegador (F12)
3. Confirme que todas as variáveis estão no `.env`
4. Tente reiniciar: `docker-compose down && docker-compose up -d`
