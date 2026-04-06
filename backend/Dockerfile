# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Instalar dependências de build
RUN apk add --no-cache git

# Copiar go.mod e go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copiar código fonte
COPY . .

# Build da aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar binário do build stage
COPY --from=builder /app/main .

# Expor porta
EXPOSE 8080

# Comando para rodar a aplicação
CMD ["./main"]
