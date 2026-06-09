# gRPC System Example - API Documentation

## Overview

Este projeto implementa uma arquitetura de integração entre dois serviços Golang:
- **Serviço gRPC**: Comunica via protocolo gRPC (Protocol Buffers)
- **Serviço REST**: API REST tradicional com endpoints HTTP

A arquitetura permite que ambos os serviços se comuniquem de forma eficiente e escalável.

---

## Architecture

```
┌─────────────────┐         gRPC          ┌──────────────────┐
│   REST Client   │◄──────────────────►│  gRPC Service    │
└─────────────────┘                      └──────────────────┘
        │                                        │
        │ HTTP/REST                              │ Protocol Buffers
        ▼                                        ▼
┌─────────────────┐                      ┌──────────────────┐
│   REST Service  │◄──────────────────►│  Client (gRPC)   │
└─────────────────┘      gRPC            └──────────────────┘
```

---

## Services

### 1. REST API Service
Serviço que expõe endpoints HTTP tradicionais.

**Port**: `8080`

**Endpoints**:
- `GET /api/v1/health` - Verifica status do serviço
- `POST /api/v1/messages` - Envia mensagem via gRPC
- `GET /api/v1/messages/{id}` - Recupera mensagem

### 2. gRPC Service
Serviço que comunica via gRPC usando Protocol Buffers.

**Port**: `50051`

**Services**:
- `MessageService.SendMessage` - Envia mensagem
- `MessageService.GetMessage` - Recupera mensagem
- `MessageService.ListMessages` - Lista todas as mensagens

---

## Prerequisites

- Go 1.19+
- Protocol Buffer Compiler (protoc)
- Go gRPC plugins:
  ```bash
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
  ```

---

## Installation

### 1. Clone o repositório
```bash
git clone <repository-url>
cd grpcSystemExample
```

### 2. Instale as dependências
```bash
go mod download
```

### 3. Gere os arquivos gRPC
```bash
protoc --go_out=. --go-grpc_out=. proto/*.proto
```

---

## Running the Services

### Start gRPC Service
```bash
go run cmd/grpc-service/main.go
```

### Start REST API Service
```bash
go run cmd/rest-service/main.go
```

---

## API Endpoints

### REST API Examples

#### Send Message via gRPC
```bash
curl -X POST http://localhost:8080/api/v1/messages \
  -H "Content-Type: application/json" \
  -d '{"text": "Hello from REST", "sender": "user123"}'
```

**Response**:
```json
{
  "id": "msg-001",
  "text": "Hello from REST",
  "sender": "user123",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

#### Get Message
```bash
curl http://localhost:8080/api/v1/messages/msg-001
```

#### Health Check
```bash
curl http://localhost:8080/api/v1/health
```

---

## gRPC Proto Definition

```protobuf
syntax = "proto3";

package message;

option go_package = "github.com/youruser/grpcSystemExample/proto";

message Message {
  string id = 1;
  string text = 2;
  string sender = 3;
  string timestamp = 4;
}

message SendMessageRequest {
  string text = 1;
  string sender = 2;
}

message GetMessageRequest {
  string id = 1;
}

service MessageService {
  rpc SendMessage(SendMessageRequest) returns (Message);
  rpc GetMessage(GetMessageRequest) returns (Message);
  rpc ListMessages(EmptyRequest) returns (MessageList);
}

message EmptyRequest {}

message MessageList {
  repeated Message messages = 1;
}
```

---

## Project Structure

```
grpcSystemExample/
├── cmd/
│   ├── grpc-service/
│   │   └── main.go
│   └── rest-service/
│       └── main.go
├── proto/
│   └── message.proto
├── internal/
│   ├── grpc/
│   │   ├── server.go
│   │   └── client.go
│   └── rest/
│       ├── handler.go
│       └── server.go
├── go.mod
├── go.sum
└── README.MD
```

---

## Configuration

### Environment Variables

```bash
# gRPC Service
GRPC_PORT=50051

# REST API Service
REST_PORT=8080
GRPC_HOST=localhost
GRPC_PORT=50051
```

---

## Error Handling

### Common Errors

| Error | Cause | Solution |
|-------|-------|----------|
| Connection refused | gRPC service não está rodando | Inicie o serviço gRPC |
| Port already in use | Porta já está em uso | Mude a porta ou libere a atual |
| gRPC: code = Unavailable | Serviço indisponível | Verifique conectividade |

---

## Development

### Running Tests
```bash
go test ./...
```

### Building
```bash
go build -o bin/grpc-service ./cmd/grpc-service
go build -o bin/rest-service ./cmd/rest-service
```

---

## Performance Considerations

- **gRPC**: Mais eficiente para comunicação inter-serviços (serialização binária)
- **REST**: Melhor para clientes externos e integrações públicas
- Use connection pooling para gRPC em produção
- Implemente rate limiting na API REST

---

## Security

- Use TLS/SSL em produção
- Implemente autenticação (JWT tokens)
- Valide todas as entradas
- Use CORS apropriadamente

```bash
# Gerar certificados TLS (auto-signed)
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes
```

---

## License

MIT

---

## Support

Para dúvidas ou problemas, abra uma issue no repositório.
