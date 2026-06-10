package main

import (
	"context"
	"grpc-service/pb"
	"grpc-service/rest" // Reaproveitando a lógica de infraestrutura que você já criou
	"log"
	"net"

	"google.golang.org/grpc"
)

// Servidor que implementa a interface gerada pelo proto
type grpcServer struct {
	pb.UnimplementedTodoServiceServer
	repo *rest.SQLiteTodoRepository // O repositório do banco fica aqui agora!
}

func (s *grpcServer) CreateTodo(ctx context.Context, req *pb.CreateTodoRequest) (*pb.TodoResponse, error) {
	// Cria a entidade interna usando a lógica que você já tinha
	todo := &rest.Todo{
		Title:       req.Title,
		Description: req.Description,
	}

	err := s.repo.Save(ctx, todo)
	if err != nil {
		return nil, err
	}

	return &pb.TodoResponse{
		Id:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		Completed:   todo.Completed,
	}, nil
}

func (s *grpcServer) ListTodos(ctx context.Context, req *pb.ListTodosRequest) (*pb.ListTodosResponse, error) {
	todos, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var pbTodos []*pb.TodoResponse
	for _, t := range todos {
		pbTodos = append(pbTodos, &pb.TodoResponse{
			Id:          t.ID,
			Title:       t.Title,
			Description: t.Description,
			Completed:   t.Completed,
		})
	}

	return &pb.ListTodosResponse{Todos: pbTodos}, nil
}

func main() {
	// Inicializa o SQLite no lado do servidor gRPC
	db, err := rest.NewSQLiteFactory()
	if err != nil {
		log.Fatalf("Erro banco: %v", err)
	}
	defer db.Close()

	repo := rest.NewSQLiteTodoRepository(db)

	// Inicia o listener TCP na porta 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Falha ao escutar porta 50051: %v", err)
	}

	baseServer := grpc.NewServer()
	pb.RegisterTodoServiceServer(baseServer, &grpcServer{repo: repo})

	log.Println("gRPC Server rodando com sucesso na porta :50051...")
	if err := baseServer.Serve(lis); err != nil {
		log.Fatalf("Falha ao rodar gRPC: %v", err)
	}
}
