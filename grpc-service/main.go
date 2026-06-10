package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	// Ajuste este import de acordo com o nome correto do seu módulo em grpc-service/go.mod
	"grpc-service/pb"
)

// O "server" agora possui uma dependência obrigatória do nosso repositório
type server struct {
	pb.UnimplementedTodoServiceServer
	repo *TodoRepository
}

// CreateTodo recebe a requisição do REST, salva no SQLite e retorna os dados gerados
func (s *server) CreateTodo(ctx context.Context, req *pb.CreateTodoRequest) (*pb.TodoResponse, error) {
	log.Printf("[gRPC] Criando Todo no SQLite: %s", req.GetTitle())

	// Gerando IDs reais e timestamp correto
	newTodo := &TodoRow{
		ID:          uuid.New().String(),
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		Completed:   false,
		CreatedAt:   time.Now().Format(time.RFC3339),
	}

	// Persiste fisicamente no banco de dados
	if err := s.repo.Save(newTodo); err != nil {
		log.Printf("Erro ao salvar no repositório: %v", err)
		return nil, err
	}

	return &pb.TodoResponse{
		Id:          newTodo.ID,
		Title:       newTodo.Title,
		Description: newTodo.Description,
		Completed:   newTodo.Completed,
		CreatedAt:   newTodo.CreatedAt,
	}, nil
}

// ListTodos busca fisicamente os dados armazenados no SQLite
func (s *server) ListTodos(ctx context.Context, req *pb.ListTodosRequest) (*pb.ListTodosResponse, error) {
	log.Println("[gRPC] Listando Todos salvos no SQLite")

	rows, err := s.repo.FindAll()
	if err != nil {
		log.Printf("Erro ao buscar no repositório: %v", err)
		return nil, err
	}

	// Mapeia os dados do banco para o formato protobuf exigido pela resposta
	var todos []*pb.TodoResponse
	for _, r := range rows {
		todos = append(todos, &pb.TodoResponse{
			Id:          r.ID,
			Title:       r.Title,
			Description: r.Description,
			Completed:   r.Completed,
			CreatedAt:   r.CreatedAt,
		})
	}

	return &pb.ListTodosResponse{Todos: todos}, nil
}

func main() {
	// 1. Inicializa o repositório SQLite (o arquivo 'todos.db' será criado localmente no container)
	repo, err := NewTodoRepository("/app/data/todos.db")
	if err != nil {
		log.Fatalf("Não foi possível inicializar o SQLite: %v", err)
	}
	defer repo.Close()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Erro ao escutar porta: %v", err)
	}

	s := grpc.NewServer()

	// 2. Registra o serviço passando o repositório criado
	pb.RegisterTodoServiceServer(s, &server{repo: repo})

	log.Println("Servidor gRPC com SQLite ativo na porta :50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Erro ao rodar servidor gRPC: %v", err)
	}
}
