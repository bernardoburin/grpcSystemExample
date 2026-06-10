package main

import (
	"database/sql"
	"fmt"

	_ "github.com/glebarez/go-sqlite" // Driver SQLite puro Go (não precisa de CGO)
)

// TodoRow representa a estrutura da tabela no banco de dados
type TodoRow struct {
	ID          string
	Title       string
	Description string
	Completed   bool
	CreatedAt   string
}

type TodoRepository struct {
	db *sql.DB
}

// NewTodoRepository inicializa o banco SQLite e garante a criação da tabela
func NewTodoRepository(dbPath string) (*TodoRepository, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir sqlite: %w", err)
	}

	// Criar a tabela se não existir
	schema := `
	CREATE TABLE IF NOT EXISTS todos (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT,
		completed INTEGER DEFAULT 0,
		created_at TEXT NOT NULL
	);`

	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("falha ao criar tabela: %w", err)
	}

	return &TodoRepository{db: db}, nil
}

func (r *TodoRepository) Close() error {
	return r.db.Close()
}

// Save insere um novo Todo no banco de dados
func (r *TodoRepository) Save(todo *TodoRow) error {
	query := `INSERT INTO todos (id, title, description, completed, created_at) VALUES (?, ?, ?, ?, ?)`
	compInt := 0
	if todo.Completed {
		compInt = 1
	}

	_, err := r.db.Exec(query, todo.ID, todo.Title, todo.Description, compInt, todo.CreatedAt)
	if err != nil {
		return fmt.Errorf("erro ao salvar todo: %w", err)
	}
	return nil
}

// FindAll busca todos os registros do banco de dados
func (r *TodoRepository) FindAll() ([]*TodoRow, error) {
	query := `SELECT id, title, description, completed, created_at FROM todos ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar todos: %w", err)
	}
	defer rows.Close()

	var todos []*TodoRow
	for rows.Next() {
		var row TodoRow
		var compInt int
		err := rows.Scan(&row.ID, &row.Title, &row.Description, &compInt, &row.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler linha do banco: %w", err)
		}
		row.Completed = (compInt == 1)
		todos = append(todos, &row)
	}

	return todos, nil
}
