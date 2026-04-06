#!/bin/bash
# Script helper para configurar Ollama no Docker

set -e

echo "🤖 Setup Ollama - Assistente de IA Local"
echo "========================================"
echo ""

# Verifica se docker-compose está rodando
if ! docker-compose ps | grep -q "financas-ollama"; then
    echo "⚠️  Ollama não está rodando. Iniciando serviços..."
    docker-compose up -d ollama
    echo "⏳ Aguardando Ollama iniciar (30 segundos)..."
    sleep 30
fi

# Verifica health do Ollama
echo "🔍 Verificando status do Ollama..."
if ! docker-compose ps ollama | grep -q "healthy"; then
    echo "⚠️  Ollama ainda não está saudável. Aguardando..."
    sleep 15
fi

# Lista modelos já instalados
echo ""
echo "📦 Modelos instalados:"
docker-compose exec ollama ollama list || echo "Nenhum modelo instalado ainda"

echo ""
echo "Escolha o modelo para baixar:"
echo "1) llama3.1:8b (Recomendado - 4.7GB) - Melhor custo/benefício"
echo "2) phi3:medium (Leve - 7.9GB) - Mais rápido"
echo "3) mistral:7b (Alternativa - 4.1GB) - Boa para análises"
echo "4) llama3.1:70b (Avançado - 40GB) - Máxima qualidade (requer muita RAM/GPU)"
echo "5) Pular (já tenho modelo instalado)"
echo ""
read -p "Opção [1-5]: " choice

case $choice in
    1)
        MODEL="llama3.1:8b"
        ;;
    2)
        MODEL="phi3:medium"
        ;;
    3)
        MODEL="mistral:7b"
        ;;
    4)
        MODEL="llama3.1:70b"
        ;;
    5)
        echo "✅ Pulando download"
        exit 0
        ;;
    *)
        echo "❌ Opção inválida"
        exit 1
        ;;
esac

echo ""
echo "📥 Baixando modelo $MODEL..."
echo "⏱️  Isso pode demorar 5-15 minutos dependendo da sua internet"
echo ""

docker-compose exec ollama ollama pull "$MODEL"

echo ""
echo "✅ Modelo $MODEL baixado com sucesso!"
echo ""
echo "📝 Atualize seu arquivo .env:"
echo "   OLLAMA_MODEL=$MODEL"
echo ""
echo "🔄 Depois reinicie o backend:"
echo "   docker-compose restart backend"
echo ""
echo "🧪 Teste o modelo:"
echo "   docker-compose exec ollama ollama run $MODEL 'Olá!'"
echo ""
echo "🎉 Setup completo!"
