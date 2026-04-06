#!/bin/bash
# Script para inicializar Ollama e baixar modelo automaticamente

echo "🤖 Aguardando Ollama iniciar..."
sleep 5

# Verifica se o modelo já está baixado
MODEL="${OLLAMA_MODEL:-llama3.1:8b}"

if ollama list | grep -q "$MODEL"; then
    echo "✅ Modelo $MODEL já está disponível"
else
    echo "📥 Baixando modelo $MODEL (isso pode demorar na primeira vez)..."
    ollama pull "$MODEL"
    echo "✅ Modelo $MODEL baixado com sucesso!"
fi

echo "🎉 Ollama pronto para uso!"
