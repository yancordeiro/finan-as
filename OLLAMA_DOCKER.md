# 🐳 Ollama no Docker - Guia Completo

## 📋 O que foi configurado

O Ollama agora roda como um serviço Docker junto com seu projeto:

```
docker-compose.yml
├─ postgres:5432       ✅ Banco de dados
├─ ollama:11434        ✅ IA Local (NOVO!)
├─ backend:8080        ✅ API Go
└─ frontend:3000       ✅ Interface Vue
```

## 🚀 Como Usar

### 1. Primeira Inicialização (Download do Modelo)

Na primeira vez, você precisa baixar o modelo de IA:

```bash
# 1. Subir todos os serviços
docker-compose up -d

# 2. Aguardar Ollama iniciar (~30 segundos)
docker-compose logs -f ollama

# 3. Baixar o modelo (escolha um):
docker-compose exec ollama ollama pull llama3.1:8b      # Recomendado (4.7GB)
# OU
docker-compose exec ollama ollama pull phi3:medium      # Mais leve (7.9GB)
# OU
docker-compose exec ollama ollama pull mistral:7b       # Alternativa (4.1GB)
```

**⏱️ Tempo estimado:** 5-15 minutos dependendo da sua internet

### 2. Verificar se o Modelo foi Baixado

```bash
# Listar modelos instalados
docker-compose exec ollama ollama list

# Saída esperada:
# NAME              ID              SIZE      MODIFIED
# llama3.1:8b       42182419e950    4.7 GB    2 minutes ago
```

### 3. Testar o Ollama

```bash
# Teste simples
docker-compose exec ollama ollama run llama3.1:8b "Olá, você está funcionando?"

# Teste com contexto financeiro
docker-compose exec ollama ollama run llama3.1:8b "Analise: gastei R$500 em delivery. Como economizar?"
```

### 4. Uso Normal (Após Configuração Inicial)

```bash
# Iniciar tudo
docker-compose up -d

# O modelo já está baixado, inicia instantaneamente!
# Acesse: http://localhost:3000
```

---

## 🔧 Configuração Avançada

### Trocar de Modelo

**1. Baixar novo modelo:**
```bash
docker-compose exec ollama ollama pull mistral:7b
```

**2. Atualizar .env:**
```env
OLLAMA_MODEL=mistral:7b
```

**3. Reiniciar backend:**
```bash
docker-compose restart backend
```

### Usar GPU NVIDIA (Acelera 10-100x)

**1. Instalar NVIDIA Container Toolkit:**
```bash
# Windows com WSL2
# Siga: https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/install-guide.html

# Linux
distribution=$(. /etc/os-release;echo $ID$VERSION_ID)
curl -s -L https://nvidia.github.io/nvidia-docker/gpgkey | sudo apt-key add -
curl -s -L https://nvidia.github.io/nvidia-docker/$distribution/nvidia-docker.list | sudo tee /etc/apt/sources.list.d/nvidia-docker.list
sudo apt-get update && sudo apt-get install -y nvidia-container-toolkit
sudo systemctl restart docker
```

**2. Editar docker-compose.yml:**
```yaml
ollama:
  # ... resto da configuração
  deploy:
    resources:
      reservations:
        devices:
          - driver: nvidia
            count: 1
            capabilities: [gpu]
```

**3. Reiniciar:**
```bash
docker-compose down
docker-compose up -d
```

**4. Verificar GPU em uso:**
```bash
docker-compose exec ollama nvidia-smi
```

---

## 📊 Monitoramento

### Ver logs do Ollama

```bash
# Logs em tempo real
docker-compose logs -f ollama

# Últimas 100 linhas
docker-compose logs --tail=100 ollama
```

### Verificar saúde dos serviços

```bash
docker-compose ps

# Saída esperada (todos "healthy"):
# NAME                 STATUS
# financas-ollama      Up (healthy)
# financas-backend     Up (healthy)
# financas-db          Up (healthy)
# financas-frontend    Up
```

### Verificar uso de recursos

```bash
# CPU, RAM, Network
docker stats financas-ollama
```

---

## 🗂️ Persistência de Dados

Os modelos são salvos em um volume Docker chamado `ollama_data`:

```bash
# Ver volumes
docker volume ls | grep ollama

# Inspecionar volume
docker volume inspect finanças_ollama_data

# Backup do modelo (opcional)
docker run --rm -v finanças_ollama_data:/data -v $(pwd):/backup alpine tar czf /backup/ollama-backup.tar.gz /data
```

**Importante:** Mesmo se você der `docker-compose down`, os modelos não são perdidos!

---

## 🌐 Networking entre Containers

O backend Go acessa Ollama via **nome do serviço**:

```go
// No código Go:
ollamaURL := "http://ollama:11434"  // ← Nome do serviço Docker!
// NÃO use: http://localhost:11434 (não funciona entre containers)
```

Como funciona:
```
backend container → Docker network → ollama container
     (financas-backend)                (financas-ollama)
           ↓
     resolve "ollama" → 172.18.0.3:11434
```

---

## ❓ Troubleshooting

### Problema: "connection refused" ao acessar Ollama

**Causa:** Ollama ainda não iniciou completamente

**Solução:**
```bash
# Verificar status
docker-compose ps ollama

# Se não estiver "healthy", aguarde mais tempo
docker-compose logs -f ollama

# Forçar restart
docker-compose restart ollama
```

### Problema: Modelo não foi baixado

**Verificar:**
```bash
docker-compose exec ollama ollama list
```

**Se vazio, baixe manualmente:**
```bash
docker-compose exec ollama ollama pull llama3.1:8b
```

### Problema: Ollama muito lento

**Opções:**

1. **Usar modelo menor:**
   ```bash
   docker-compose exec ollama ollama pull phi3:medium
   # Atualizar OLLAMA_MODEL no .env
   ```

2. **Alocar mais RAM ao Docker:**
   - Docker Desktop → Settings → Resources → Memory: 8GB+

3. **Usar GPU** (veja seção acima)

### Problema: Erro "model not found" no backend

**Causa:** Variável `OLLAMA_MODEL` no .env não corresponde ao modelo baixado

**Solução:**
```bash
# 1. Ver modelos disponíveis
docker-compose exec ollama ollama list

# 2. Atualizar .env para corresponder
OLLAMA_MODEL=llama3.1:8b  # Use o nome exato

# 3. Reiniciar backend
docker-compose restart backend
```

---

## 🧹 Limpeza

### Remover modelo específico

```bash
docker-compose exec ollama ollama rm llama3.1:8b
```

### Limpar tudo e recomeçar

```bash
# Para serviços
docker-compose down

# Remove volume (CUIDADO: apaga modelos baixados!)
docker volume rm finanças_ollama_data

# Reiniciar do zero
docker-compose up -d
```

---

## 📈 Comparação de Performance

### CPU vs GPU (Llama 3.1 8B)

| Hardware | Tokens/segundo | Tempo resposta (100 tokens) |
|----------|----------------|----------------------------|
| CPU i7 (8 cores) | ~10-20 tokens/s | ~5-10 segundos |
| CPU i9 (16 cores) | ~20-30 tokens/s | ~3-5 segundos |
| GPU RTX 3060 (12GB) | ~80-100 tokens/s | ~1 segundo |
| GPU RTX 4090 (24GB) | ~150-200 tokens/s | ~0.5 segundo |

### Modelos (em CPU i7)

| Modelo | Tamanho | RAM | Velocidade | Qualidade |
|--------|---------|-----|------------|-----------|
| phi3:medium | 7.9GB | 8GB | 25 tok/s ⚡⚡⚡ | ⭐⭐⭐ |
| llama3.1:8b | 4.7GB | 8GB | 15 tok/s ⚡⚡ | ⭐⭐⭐⭐ |
| mistral:7b | 4.1GB | 8GB | 18 tok/s ⚡⚡ | ⭐⭐⭐⭐ |
| llama3.1:70b | 40GB | 64GB | 2 tok/s ⚡ | ⭐⭐⭐⭐⭐ |

---

## 🎯 Próximos Passos

Agora que Ollama está configurado:

1. ✅ Ollama rodando no Docker
2. ✅ Modelo baixado
3. ⏭️ **Implementar código Go para chamar Ollama**
4. ⏭️ **Criar interface de chat no frontend**
5. ⏭️ **Testar análises financeiras**

Continue seguindo o README.md principal!

---

## 📚 Referências

- Ollama Docs: https://github.com/ollama/ollama/blob/main/docs/docker.md
- Ollama API: https://github.com/ollama/ollama/blob/main/docs/api.md
- Modelos disponíveis: https://ollama.com/library
