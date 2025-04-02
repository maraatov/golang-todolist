package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	todo "todolists"
)

type TodoListPostgres struct {
	db *pgxpool.Pool
}

func NewTodoListPostgres(db *pgxpool.Pool) *TodoListPostgres {
	return &TodoListPostgres{db: db}
}

func (r *TodoListPostgres) Create(userId int, list todo.TodoList) (int, error) {
	tx, err := r.db.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return 0, err
	}

	var id int
	err = tx.QueryRow(context.Background(), "INSERT INTO todo_lists (title, description) VALUES ($1, $2) RETURNING id", list.Title, list.Description).Scan(&id)
	if err != nil {
		tx.Rollback(context.Background())
		return 0, err
	}

	_, err = tx.Exec(context.Background(), "INSERT INTO users_lists (user_id, list_id) VALUES ($1, $2)", userId, id)
	if err != nil {
		tx.Rollback(context.Background())
		return 0, err
	}

	if err = tx.Commit(context.Background()); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *TodoListPostgres) Delete(userId, listId int) error {
	query := `DELETE FROM todo_lists tl USING users_lists ul WHERE tl.id = ul.list_id AND ul.user_id = $1 AND ul.list_id = $2`
	_, err := r.db.Exec(context.Background(), query, userId, listId)
	return err
}

func (r *TodoListPostgres) GetById(userId, listId int) (todo.TodoList, error) {
	var list todo.TodoList
	query := `SELECT tl.id, tl.title, tl.description FROM todo_lists tl INNER JOIN users_lists ul ON tl.id = ul.list_id WHERE ul.user_id = $1 AND ul.list_id = $2`
	err := r.db.QueryRow(context.Background(), query, userId, listId).Scan(&list.Id, &list.Title, &list.Description)
	return list, err
}
