package todo

import "time"

type ToDo struct {
	ID           int64
	Title        string
	Started      bool
	Editing      bool
	TimeStart    time.Time
	TimeElapsed  time.Duration
	TimeEstimate time.Duration
}

type ToDoRepository interface {
	Migrate() error
	Add(todo ToDo) (*ToDo, error)
	All() ([]ToDo, error)
	GetByID(id int64) (*ToDo, error)
	Update(id int64, updatedToDo ToDo) (*ToDo, error)
	Delete(id int64) error
}
