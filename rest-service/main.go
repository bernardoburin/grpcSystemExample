package rest

import (
	"log"
	"net/http"
	"rest-service/rest" // substitua pelo caminho do módulo correto do seu projeto
)

func main() {
	// 1. Inicializa a infraestrutura de banco de dados (SQLite Factory)
	db, err := rest.NewSQLiteFactory()
	if err != nil {
		log.Fatalf("Falha ao conectar no SQLite: %v", err)
	}
	defer db.Close()

	// 2. Injeta as dependências seguindo o fluxo DDD (Infra -> App -> Controller)
	todoRepo := rest.NewSQLiteTodoRepository(db)
	appService := rest.NewTodoApplicationService(todoRepo)
	todoController := rest.NewTodoController(appService)

	// 3. Define as Rotas/Endpoints para o Front-end
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

	// 4. Inicia o servidor HTTP na porta 8080 configurada no projeto
	log.Println("REST API Server rodando com sucesso na porta :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
