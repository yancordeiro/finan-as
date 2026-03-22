# 💰 Finanças - Controle Financeiro Pessoal

Aplicação web full-stack de controle financeiro doméstico com visual dark mode premium, autenticação Google OAuth, gestão de transações e importação de extratos OFX.

## 🚀 Stack Tecnológica

### Backend
- **Go** (Golang) 1.21+
- **Fiber** - Framework web rápido e expressivo
- **PostgreSQL** - Banco de dados relacional
- **GORM** - ORM para Go
- **JWT** - Autenticação e autorização
- **OFXgo** - Parser de arquivos OFX

### Frontend
- **Vue 3** - Framework JavaScript progressivo
- **Vite** - Build tool ultrarrápida
- **Tailwind CSS** - Framework CSS utilitário
- **Pinia** - State management
- **Vue Router** - Roteamento
- **Chart.js** - Gráficos e visualizações
- **Axios** - Cliente HTTP

### Infraestrutura
- **Docker** & **Docker Compose** - Containerização
- **PostgreSQL 15** - Banco de dados
- **Nginx** - Servidor web (produção)

---

## 📋 Pré-requisitos

- [Docker](https://www.docker.com/get-started) (versão 20.10+)
- [Docker Compose](https://docs.docker.com/compose/install/) (versão 2.0+)
- [Node.js](https://nodejs.org/) (versão 20+) - apenas para desenvolvimento local
- [Go](https://golang.org/dl/) (versão 1.21+) - apenas para desenvolvimento local
- Conta Google Cloud Platform com OAuth 2.0 configurado

---

## 🛠️ Configuração Inicial

### 1. Clonar o repositório

```bash
git clone <repository-url>
cd finanças
```

### 2. Configurar Google OAuth 2.0

1. Acesse o [Google Cloud Console](https://console.cloud.google.com/)
2. Crie um novo projeto ou selecione um existente
3. Ative a **Google+ API**
4. Vá para **Credenciais** > **Criar Credenciais** > **ID do cliente OAuth 2.0**
5. Configure a tela de consentimento OAuth
6. Adicione as origens autorizadas:
   - `http://localhost:3000`
   - `http://localhost:8080`
7. Adicione os URIs de redirecionamento autorizados:
   - `http://localhost:3000/auth/callback`
8. Copie o **Client ID** e **Client Secret**

### 3. Configurar variáveis de ambiente

Copie o arquivo `.env.example` para `.env`:

```bash
cp .env.example .env
```

Edite o arquivo `.env` e preencha as variáveis:

```env
# JWT Configuration (gere secrets fortes!)
JWT_SECRET=seu-secret-jwt-aqui-minimo-256-bits
JWT_REFRESH_SECRET=seu-secret-refresh-aqui-minimo-256-bits

# Google OAuth 2.0
GOOGLE_CLIENT_ID=seu-google-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=seu-google-client-secret
GOOGLE_REDIRECT_URI=http://localhost:3000/auth/callback

# URLs
FRONTEND_URL=http://localhost:3000
VITE_API_URL=http://localhost:8080
```

**⚠️ IMPORTANTE:** Para gerar secrets seguros, use:

```bash
# Linux/Mac
openssl rand -hex 32

# Windows (PowerShell)
[Convert]::ToBase64String((1..32 | ForEach-Object { Get-Random -Maximum 256 }))
```

---

## 🐳 Rodando com Docker (Recomendado)

### Iniciar todos os serviços

```bash
docker-compose up -d
```

Este comando irá:
- ✅ Criar o banco de dados PostgreSQL
- ✅ Construir e iniciar o backend Go
- ✅ Construir e iniciar o frontend Vue
- ✅ Executar migrations e seeds automáticos

### Acessar a aplicação

- **Frontend:** http://localhost:3000
- **Backend API:** http://localhost:8080
- **Health Check:** http://localhost:8080/health

### Visualizar logs

```bash
# Todos os serviços
docker-compose logs -f

# Apenas backend
docker-compose logs -f backend

# Apenas frontend
docker-compose logs -f frontend
```

### Parar os serviços

```bash
docker-compose down
```

### Resetar o banco de dados

```bash
docker-compose down -v
docker-compose up -d
```

---

## 💻 Desenvolvimento Local (sem Docker)

### Backend

```bash
cd backend

# Instalar dependências
go mod download

# Criar arquivo .env
cp .env.example .env
# (Edite o .env com suas configurações)

# Rodar PostgreSQL localmente (ou use Docker)
docker run -d \
  --name financas-postgres \
  -e POSTGRES_USER=financas \
  -e POSTGRES_PASSWORD=financas_dev_password \
  -e POSTGRES_DB=financas_db \
  -p 5432:5432 \
  postgres:15-alpine

# Rodar a aplicação
go run cmd/api/main.go
```

### Frontend

```bash
cd frontend

# Instalar dependências
npm install

# Criar arquivo .env
cp .env.example .env.local
# (Edite o .env.local com suas configurações)

# Rodar em modo desenvolvimento
npm run dev
```

---

## 📁 Estrutura do Projeto

```
finanças/
├── backend/                    # Backend Go
│   ├── cmd/api/               # Entry point
│   ├── internal/              # Código interno
│   │   ├── config/           # Configurações
│   │   ├── database/         # DB connection & migrations
│   │   ├── middleware/       # Middlewares (auth, CORS)
│   │   ├── models/           # Models do banco
│   │   ├── handlers/         # HTTP handlers
│   │   ├── services/         # Lógica de negócio
│   │   └── routes/           # Definição de rotas
│   └── pkg/utils/            # Utilitários (JWT, response)
│
├── frontend/                   # Frontend Vue 3
│   ├── src/
│   │   ├── assets/           # Estilos e assets
│   │   ├── components/       # Componentes Vue
│   │   │   ├── ui/          # Componentes base
│   │   │   ├── layout/      # Layout components
│   │   │   ├── dashboard/   # Dashboard widgets
│   │   │   ├── transactions/
│   │   │   ├── ofx/
│   │   │   └── categories/
│   │   ├── composables/      # Composables Vue
│   │   ├── stores/           # Pinia stores
│   │   ├── router/           # Vue Router
│   │   └── views/            # Páginas
│   └── public/
│
├── docker-compose.yml          # Orquestração Docker
├── .env.example               # Exemplo de variáveis
└── README.md                  # Este arquivo
```

---

## 🎨 Funcionalidades

### ✅ Implementadas (100%)

- [x] Estrutura base do projeto (backend + frontend)
- [x] Configuração Docker Compose
- [x] Models do banco de dados (User, Category, Transaction)
- [x] Migrations automáticas
- [x] Seed de 13 categorias padrão
- [x] Layout responsivo com sidebar e topbar
- [x] Tema dark mode premium
- [x] Rotas protegidas com autenticação
- [x] **Autenticação completa com Google OAuth 2.0**
- [x] **CRUD de categorias completo**
- [x] **Toast notifications**
- [x] **Loading states e skeletons**
- [x] **ColorPicker com 16 cores preset**
- [x] **EmojiPicker com 100+ emojis**

### 🚧 Em Desenvolvimento

- [ ] CRUD de transações com filtros
- [ ] Dashboard com métricas e gráficos (Chart.js)
- [ ] Importação de arquivos OFX
- [ ] Exportação para CSV

### ⚠️ Removido

- ~~CRUD de contas bancárias~~ (não será implementado - transações simplificadas)

---

## 🎨 Paleta de Cores (Dark Mode)

```css
Background principal: #0f0f13
Cards/Painéis: #1a1a24
Sidebar: #13131c
Borda sutil: #2a2a3a

Primária (azul elétrico): #4f8ef7
Secundária (ciano): #00d4ff
Accent (roxo): #a259ff

Verde (positivo): #00e676
Vermelho (negativo): #ff4d4d

Texto primário: #e8e8f0
Texto secundário: #7a7a9a
```

---

## 📝 API Endpoints

### Autenticação (Públicas)

```
POST   /api/auth/google          # Login com Google OAuth
POST   /api/auth/refresh         # Renovar access token
POST   /api/auth/logout          # Logout (protegida)
```

### Categorias (Protegidas)

```
GET    /api/categories           # Listar todas
POST   /api/categories           # Criar nova
PUT    /api/categories/:id       # Atualizar
DELETE /api/categories/:id       # Deletar
```

### ~~Contas~~ (Removido - não implementado)

### Transações (Protegidas)

```
GET    /api/transactions         # Listar com filtros
POST   /api/transactions         # Criar nova
PUT    /api/transactions/:id     # Atualizar
DELETE /api/transactions/:id     # Deletar
GET    /api/transactions/export  # Exportar CSV
```

### Dashboard (Protegidas)

```
GET    /api/dashboard/summary             # Resumo do mês
GET    /api/dashboard/by-category         # Por categoria
GET    /api/dashboard/monthly-flow        # Fluxo mensal
GET    /api/dashboard/recent-transactions # Últimas transações
```

### OFX (Protegidas)

```
POST   /api/ofx/upload           # Upload e parse OFX
POST   /api/ofx/import           # Importar transações
```

---

## 🧪 Testes

```bash
# Backend
cd backend
go test ./...

# Frontend
cd frontend
npm run test
```

---

## 🔒 Segurança

- ✅ JWT com expiração de 1 hora
- ✅ Refresh tokens com expiração de 30 dias
- ✅ Cookies httpOnly para refresh tokens
- ✅ CORS configurado
- ✅ Validação de propriedade (user só acessa seus dados)
- ✅ Prepared statements (proteção contra SQL injection)
- ✅ Valores monetários em centavos (evita problemas com float)

---

## 📦 Build para Produção

```bash
# Build de todos os serviços
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build

# Ou manualmente:

# Backend
cd backend
CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Frontend
cd frontend
npm run build
```

---

## 🐛 Troubleshooting

### Erro ao conectar ao banco de dados

Verifique se o PostgreSQL está rodando:

```bash
docker-compose ps
```

Verifique os logs do PostgreSQL:

```bash
docker-compose logs postgres
```

### Erro de CORS no frontend

Verifique se a variável `FRONTEND_URL` no backend está correta no `.env`.

### Erro de token inválido

1. Verifique se `JWT_SECRET` e `JWT_REFRESH_SECRET` estão configurados
2. Limpe o localStorage do browser
3. Faça logout e login novamente

---

## 📄 Licença

Este projeto está sob a licença MIT.

---

## 👥 Contribuindo

Contribuições são bem-vindas! Sinta-se à vontade para abrir issues e pull requests.

---

## 📞 Suporte

Para dúvidas e suporte, abra uma issue no repositório.

---

**Desenvolvido com ❤️ usando Go, Vue.js e Tailwind CSS**
#   f i n a n - a s  
 