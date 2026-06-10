# --- Estágio 1: Compilação (Build) ---
FROM golang:1.26.4-alpine AS builder

# Instala dependências de compilação C para o SQLite
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copia os arquivos de definição de módulos e baixa as dependências
COPY go.mod go.sum ./
RUN go mod download

# Copia todo o conteúdo da pasta local do serviço para dentro do container
COPY . .

# COMPILAÇÃO:
# Compila estaticamente as partes possíveis e garante o build do arquivo correto
RUN CGO_ENABLED=1 GOOS=linux go build -o rest-service \
    $(if [ -f "./main.go" ]; then echo "./main.go"; else echo "./cmd/rest-service/main.go"; fi)

# Força a permissão de execução no binário gerado
RUN chmod +x rest-service

# --- Estágio 2: Execução ---
FROM alpine:latest

# libc6-compat e gcompat ajudam a resolver problemas de linkagem dinâmica do SQLite no Alpine
RUN apk add --no-cache ca-certificates libc6-compat gcompat

WORKDIR /app

# Copia o binário compilado e com permissões do estágio anterior
COPY --from=builder /app/rest-service .

EXPOSE 8080

CMD ["./rest-service"]