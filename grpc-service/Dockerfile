# --- Estágio 1: Compilação (Build) ---
FROM golang:1.26.4-alpine AS builder

# Instala ferramentas essenciais de compilação C
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copia as definições de módulos e baixa as dependências
COPY go.mod ./
# Se houver o arquivo go.sum na pasta do grpc-service, descomente a linha abaixo:
# COPY go.sum ./
RUN go mod download

# Copia todo o código-fonte do serviço gRPC
COPY . .

# DETERMINA O CAMINHO DO MAIN E COMPILA ESTATICAMENTE:
# Verifica se o main.go está na raiz do serviço ou dentro de cmd/grpc-service
RUN MAIN_PATH=$(if [ -f "./main.go" ]; then echo "./main.go"; else echo "./cmd/grpc-service/main.go"; fi) && \
    CGO_ENABLED=1 GOOS=linux go build \
    -ldflags "-extldflags '-static'" \
    -o grpc-service ${MAIN_PATH}

# Garante explicitamente a permissão de execução no binário gerado
RUN chmod +x grpc-service

# --- Estágio 2: Execução ---
FROM alpine:latest

# Instala certificados e pacotes de compatibilidade essenciais para runtime
RUN apk add --no-cache ca-certificates libc6-compat gcompat

WORKDIR /app

# Copia o binário final e limpo do estágio de build
COPY --from=builder /app/grpc-service .

# Expõe a porta padrão do gRPC configurada no projeto
EXPOSE 50051

# Comando para rodar o serviço gRPC
CMD ["./grpc-service"]