package main

import (
	"log"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Importante: conecta usando o nome do serviço do docker-compose
	conn, err := grpc.Dial("grpc-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Não foi possível conectar ao gRPC: %v", err)
	}
	defer conn.Close()

	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status": "conectado ao gRPC"}`))
	})

	log.Println("Servidor REST rodando na porta :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
