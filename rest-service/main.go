package main

import (
	"log"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// Substitua pelo caminho do módulo correto configurado no seu go.mod do rest-service
	"rest-service/pb"
	"rest-service/rest"
)

func main() {
	// Conecta ao serviço gRPC utilizando o nome do serviço definido no docker-compose
	conn, err := grpc.Dial("grpc-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Não foi possível conectar ao gRPC: %v", err)
	}
	defer conn.Close()

	// 1. Instancia o cliente gRPC gerado pelo proto
	grpcClient := pb.NewTodoServiceClient(conn)

	// 2. Instancia o Application Service injetando o cliente gRPC
	appService := rest.NewTodoApplicationService(grpcClient)

	// 3. Instancia o HTTP Controller injetando o Application Service
	todoController := rest.NewTodoController(appService)

	// Definição dos dois endpoints solicitados
	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			todoController.ListTodosHandler(w, r)
		case http.MethodPost:
			todoController.CreateTodoHandler(w, r)
		default:
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "OK", http.StatusOK)
	})

	log.Println("Servidor REST rodando na porta :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
