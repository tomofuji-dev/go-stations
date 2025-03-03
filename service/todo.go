package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/TechBowl-japan/go-stations/model"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT * FROM todos WHERE id = ?`
	)

	// Prepare the SQL statement
	stmt, err := s.db.PrepareContext(ctx, insert)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the statement with provided subject and description
	res, err := stmt.ExecContext(ctx, subject, description)
	if err != nil {
		return nil, err
	}

	// Retrieve the ID of the newly inserted TODO
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Confirm and retrieve the newly inserted TODO
	var todo model.TODO
	err = s.db.QueryRowContext(ctx, confirm, id).Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		allRead    = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC`
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)

	var rows *sql.Rows
	var err error

	// Query with or without the prevID
	if prevID == 0 && size == 0 {
		rows, err = s.db.QueryContext(ctx, allRead)
	} else if prevID > 0 {
		rows, err = s.db.QueryContext(ctx, readWithID, prevID, size)
	} else {
		rows, err = s.db.QueryContext(ctx, read, size)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := make([]*model.TODO, 0)

	for rows.Next() {
		var todo model.TODO
		err := rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	// Prepare the SQL statement for update
	stmt, err := s.db.PrepareContext(ctx, update)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the update statement
	res, err := stmt.ExecContext(ctx, subject, description, id)
	if err != nil {
		return nil, err
	}

	// Check if any rows were updated
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, &model.ErrNotFound{Message: "No TODO were updated"}
	}

	// Confirm and retrieve the updated TODO
	var todo model.TODO
	err = s.db.QueryRowContext(ctx, confirm, id).Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return nil, err
	}
	todo.ID = id

	return &todo, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	if len(ids) == 0 {
		return nil
	}

	// prepare SQL query with placeholders for ids
	query := fmt.Sprintf(deleteFmt, strings.Repeat(",?", len(ids)-1))
	query = strings.TrimLeft(query, ",")

	// prepare arguments for query
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	// execute query
	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	// check if any row was affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return &model.ErrNotFound{
			Message: "There is no deleted TODO",
		}
	}

	return nil
}
