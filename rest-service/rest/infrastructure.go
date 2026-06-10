package rest

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

// Entidade de domínio Todo
type Todo struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

// Factory de conexão com o banco de dados SQLite
func NewSQLiteFactory() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./todos.db")
	if err != nil {
		return nil, err
	}

	// Criação da tabela para testes caso não exista
	statement, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS todos (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT,
			completed INTEGER NOT NULL
		);
	`)
	if err != nil {
		return nil, err
	}
	_, err = statement.Exec()
	return db, err
}

// Interface de Infraestrutura/Repositório (definida no domínio/aplicação, implementada aqui)
type TodoRepository interface {
	Save(ctx context.Context, todo *Todo) error
	FindAll(ctx context.Context) ([]*Todo, error)
}

// Implementação do Repositório usando SQLite
type SQLiteTodoRepository struct {
	db *sql.DB
}

func NewSQLiteTodoRepository(db *sql.DB) *SQLiteTodoRepository {
	return &SQLiteTodoRepository{db: db}
}

func (r *SQLiteTodoRepository) Save(ctx context.Context, todo *Todo) error {
	// Gera ID se não existir
	if todo.ID == "" {
		todo.ID = uuid.New().String()
	}
	query := `INSERT INTO todos (id, title, description, completed) VALUES (?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, todo.ID, todo.Title, todo.Description, todo.Completed)
	return err
}

func (r *SQLiteTodoRepository) FindAll(ctx context.Context) ([]*Todo, error) {
	query := `SELECT id, title, description, completed FROM todos`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*Todo
	for rows.Next() {
		var t Todo
		var completedInt int
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &completedInt); err != nil {
			return nil, err
		}
		t.Completed = completedInt == 1
		todos = append(todos, &t)
	}
	return todos, nil
}
