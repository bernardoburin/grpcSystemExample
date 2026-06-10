package rest

import (
	"context"
	"rest-service/pb" // Pasta onde ficarão os arquivos gerados do proto
	"time"
)

type TodoApplicationService struct {
	grpcClient pb.TodoServiceClient // Substitui o TodoRepository
}

func NewTodoApplicationService(client pb.TodoServiceClient) *TodoApplicationService {
	return &TodoApplicationService{grpcClient: client}
}

func (s *TodoApplicationService) CreateTodo(ctx context.Context, title, description string) (*Todo, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	// Transforma a chamada local em uma chamada de rede gRPC para o outro microsserviço
	resp, err := s.grpcClient.CreateTodo(ctx, &pb.CreateTodoRequest{
		Title:       title,
		Description: description,
	})
	if err != nil {
		return nil, err
	}

	// Mapeia a resposta do gRPC de volta para a struct do domínio do REST
	return &Todo{
		ID:          resp.Id,
		Title:       resp.Title,
		Description: resp.Description,
		Completed:   resp.Completed,
	}, nil
}

func (s *TodoApplicationService) ListTodos(ctx context.Context) ([]*Todo, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	resp, err := s.grpcClient.ListTodos(ctx, &pb.ListTodosRequest{})
	if err != nil {
		return nil, err
	}

	var todos []*Todo
	for _, t := range resp.Todos {
		todos = append(todos, &Todo{
			ID:          t.Id,
			Title:       t.Title,
			Description: t.Description,
			Completed:   t.Completed,
		})
	}
	return todos, nil
}