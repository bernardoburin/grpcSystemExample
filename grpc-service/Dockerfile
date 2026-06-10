# ==========================================
# Estágio 1: Build (Compilação)
# ==========================================
# CORREÇÃO: "golang:alpine" é a tag oficial para a versão mais recente com Alpine
FROM golang:alpine AS builder

# Instala ferramentas essenciais para download de dependências
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copia os arquivos de definição de dependências primeiro (otimiza o cache do Docker)
COPY go.mod go.sum ./
RUN go mod download

# Copia todo o restante do código-fonte do gRPC (incluindo main.go e repository.go)
COPY . .

# Compila o DIRETÓRIO inteiro com CGO desativado
RUN if [ -f "./main.go" ]; then \
        CGO_ENABLED=0 GOOS=linux go build -o grpc-service .; \
    else \
        CGO_ENABLED=0 GOOS=linux go build -o grpc-service ./cmd/grpc-service; \
    fi

# ==========================================
# Estágio 2: Execução (Imagem Final Leve)
# ==========================================
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Cria a pasta onde o SQLite guardará os dados persistidos
RUN mkdir -p /app/data

# Copia o binário compilado no Estágio 1 para esta imagem limpa
COPY --from=builder /app/grpc-service .

# Expõe a porta padrão do gRPC
EXPOSE 50051

# Comando para rodar a aplicação
CMD ["./grpc-service"]