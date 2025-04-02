package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	todo "todolists"
)

type TodoItemPostgres struct {
	db *pgxpool.Pool
}

func NewTodoItemPostgres(db *pgxpool.Pool) *TodoItemPostgres {
	return &TodoItemPostgres{db: db}
}

func (r *TodoItemPostgres) Create(listId int, item todo.TodoItem) (int, error) {
	tx, err := r.db.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return 0, err
	}

	var itemId int
	createItemQuery := `INSERT INTO todo_items (title, description) VALUES ($1, $2) RETURNING id`
	row := tx.QueryRow(context.Background(), createItemQuery, item.Title, item.Description)
	err = row.Scan(&itemId)
	if err != nil {
		tx.Rollback(context.Background())
		return 0, err
	}

	createListItemsQuery := `INSERT INTO lists_items (list_id, item_id) VALUES ($1, $2)`
	_, err = tx.Exec(context.Background(), createListItemsQuery, listId, itemId)
	if err != nil {
		tx.Rollback(context.Background())
		return 0, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return 0, err
	}

	return itemId, nil
}

func (r *TodoItemPostgres) GetAll(userId, listId int) ([]todo.TodoItem, error) {
	var items []todo.TodoItem
	query := `SELECT ti.id, ti.title, ti.description, ti.done 
			  FROM todo_items ti
			  INNER JOIN lists_items li ON li.item_id = ti.id
			  INNER JOIN users_lists ul ON ul.list_id = li.list_id 
			  WHERE li.list_id = $1 AND ul.user_id = $2`

	rows, err := r.db.Query(context.Background(), query, listId, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item todo.TodoItem
		if err := rows.Scan(&item.Id, &item.Title, &item.Description, &item.Done); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *TodoItemPostgres) GetById(userId, itemId int) (todo.TodoItem, error) {
	var item todo.TodoItem
	query := `SELECT ti.id, ti.title, ti.description, ti.done 
			  FROM todo_items ti
			  INNER JOIN lists_items li ON li.item_id = ti.id
			  INNER JOIN users_lists ul ON ul.list_id = li.list_id 
			  WHERE ti.id = $1 AND ul.user_id = $2`
	err := r.db.QueryRow(context.Background(), query, itemId, userId).Scan(&item.Id, &item.Title, &item.Description, &item.Done)
	if err != nil {
		return item, err
	}

	return item, nil
}

func (r *TodoItemPostgres) Delete(userId, itemId int) error {
	query := `DELETE FROM todo_items ti 
			  USING lists_items li, users_lists ul 
			  WHERE ti.id = li.item_id AND li.list_id = ul.list_id 
			  AND ul.user_id = $1 AND ti.id = $2`
	_, err := r.db.Exec(context.Background(), query, userId, itemId)
	return err
}
