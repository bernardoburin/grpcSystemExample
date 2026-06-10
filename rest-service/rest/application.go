package rest

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// Entidade de Domínio
type Todo struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

// Regra de validação de domínio simples
func (t *Todo) Validate() error {
	if t.Title == "" {
		return errors.New("o título do todo list não pode ser vazio")
	}
	return nil
}

// Estruturas de Dados de Entrada (DTOs) da Aplicação
type CreateTodoInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Serviço de Aplicação (Application Service)
type TodoApplicationService struct {
	repo TodoRepository
}

func NewTodoApplicationService(repo TodoRepository) *TodoApplicationService {
	return &TodoApplicationService{repo: repo}
}

// Regra de Negócio: Criar uma nova lista/item ToDo
func (s *TodoApplicationService) CreateTodo(ctx context.Context, input CreateTodoInput) (*Todo, error) {
	todo := &Todo{
		ID:          uuid.New().String(),
		Title:       input.Title,
		Description: input.Description,
		Completed:   false,
	}

	if err := todo.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, todo); err != nil {
		return nil, err
	}

	return todo, nil
}

// Regra de Negócio: Listar todos os ToDos
func (s *TodoApplicationService) ListTodos(ctx context.Context) ([]*Todo, error) {
	return s.repo.FindAll(ctx)
}
