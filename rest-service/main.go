package main

import (
	"log"
	"net/http"
	"rest-service/rest"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 1. Conecta ao microsserviço gRPC usando o nome do container definido no docker-compose
	conn, err := grpc.Dial("grpc-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Não foi possível conectar ao gRPC service: %v", err)
	}
	defer conn.Close()

	// 2. Cria o cliente gRPC baseado no contrato proto
	grpcClient := pb.NewTodoServiceClient(conn)

	// 3. Injeta o cliente gRPC no Application Service (DDD)
	appService := rest.NewTodoApplicationService(grpcClient)
	todoController := rest.NewTodoController(appService)

	// 4. Mantém as rotas HTTP originais intactas para o Frontend
	http.HandleFunc("/api/v1/todos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			todoController.CreateTodoHandler(w, r)
		case http.MethodGet:
			todoController.ListTodosHandler(w, r)
		default:
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		}
	})

	log.Println("REST API Gateway rodando com sucesso na porta :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Erro ao iniciar o servidor HTTP: %v", err)
	}
}
