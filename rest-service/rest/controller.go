package rest

import (
	"encoding/json"
	"net/http"
)

// HTTP Controller (Interface Adapter)
type TodoController struct {
	appService *TodoApplicationService
}

func NewTodoController(appService *TodoApplicationService) *TodoController {
	return &TodoController{appService: appService}
}

// POST /api/v1/todos
func (c *TodoController) CreateTodoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var input CreateTodoInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	todo, err := c.appService.CreateTodo(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

// GET /api/v1/todos
func (c *TodoController) ListTodosHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	todos, err := c.appService.ListTodos(r.Context())
	if err != nil {
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if todos == nil {
		todos = []*Todo{} // Evita retornar 'null' no JSON do front-end
	}
	json.NewEncoder(w).Encode(todos)
}
