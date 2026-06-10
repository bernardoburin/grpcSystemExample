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
# Forçamos o CGO e adicionamos flags para garantir um binário estático compatível com o Alpine.
# O uso de './...' compila o pacote atual (onde está o main.go) de forma limpa.
RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags "-extldflags '-static'" \
    -o rest-service .

# --- Estágio 2: Execução ---
FROM alpine:latest

# libc6-compat e gcompat ajudam a resolver problemas de linkagem dinâmica do SQLite no Alpine
RUN apk add --no-cache ca-certificates libc6-compat gcompat

WORKDIR /app

# Copia o binário compilado do estágio anterior
COPY --from=builder /app/rest-service .

# GARANTE A PERMISSÃO: Força o binário a ser executável no estágio final
RUN chmod +x ./rest-service

EXPOSE 8080

CMD ["./rest-service"]