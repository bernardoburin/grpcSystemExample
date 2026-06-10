package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	// substitua pelo caminho correto do seu pacote pb
	"grpc-service/pb"
)

type server struct {
	pb.UnimplementedTodoServiceServer
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Erro ao escutar porta: %v", err)
	}
	s := grpc.NewServer()
	// pb.RegisterTodoServiceServer(s, &server{})
	log.Println("Servidor gRPC rodando na porta :50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Erro ao iniciar gRPC: %v", err)
	}
}
