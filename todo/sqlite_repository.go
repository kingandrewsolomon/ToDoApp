package todo

import (
	"database/sql"
	"errors"
	"time"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		db: db,
	}
}

func (r *SQLiteRepository) Migrate() error {
	query := `
	create table if not exists todos(
		id integer primary key autoincrement,
		title text not null,
		started boolean,
		editing boolean,
		timestart datetime,
		timeelapsed integer,
		timeestimate integer 
	);`
	_, err := r.db.Exec(query)
	return err
}

func (r *SQLiteRepository) Add(todo ToDo) (*ToDo, error) {
	res, err := r.db.Exec("insert into todos (title, started, editing, timestart, timeelapsed, timeestimate) values(?, ?, ?, ?, ?, ?)", todo.Title, todo.Started, todo.Editing, todo.TimeStart, todo.TimeElapsed, todo.TimeEstimate)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	todo.ID = id
	return &todo, nil
}

func (r *SQLiteRepository) All() ([]ToDo, error) {
	rows, err := r.db.Query("select * from todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []ToDo
	for rows.Next() {
		var todo ToDo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Started, &todo.Editing, &todo.TimeStart, &todo.TimeElapsed, &todo.TimeEstimate); err != nil {
			return nil, err
		}

		all = append(all, todo)
	}
	return all, nil
}

func (r *SQLiteRepository) GetByID(id int64) (*ToDo, error) {
	row := r.db.QueryRow("select * from todos where id = ?", id)

	var todo ToDo
	if err := row.Scan(&todo.ID,&todo.Title, &todo.Started, &todo.Editing, &todo.TimeStart, &todo.TimeElapsed, &todo.TimeEstimate); err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *SQLiteRepository) Update(id int64, updatedTodo ToDo) (*ToDo, error) {
	res, err := r.db.Exec("Update todos set title = ?, started = ?, editing = ?, timestart = ?, timeelapsed = ?, timeestimate = ? where id = ?", updatedTodo.Title, updatedTodo.Started, updatedTodo.Editing, updatedTodo.TimeStart, updatedTodo.TimeElapsed, updatedTodo.TimeEstimate, id)
	if err != nil {
		return nil, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, errors.New("update failed")
	}

	return &updatedTodo, nil
}

func (r *SQLiteRepository) Delete(id int64) error {
	res, err := r.db.Exec("delete from todos where id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("delete failed")
	}

	return nil
}

func (r *SQLiteRepository) StartTiming(t *ToDo) (*ToDo, error) {
	t.Started = true
	t.TimeStart = time.Now()
	todo, err := r.Update(t.ID, *t)
	return todo, err
}

func (r *SQLiteRepository) StopTiming(t *ToDo) (*ToDo, error) {
	t.Started = false
	t.TimeElapsed += time.Since(t.TimeStart)
	todo, err := r.Update(t.ID, *t)
	return todo, err
}
