package main

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"htmx/todo"
)

func GetID(selectedTodo string) int64 {
	t, err := strconv.ParseInt(selectedTodo, 10, 64)
	if err != nil {
		panic(err)
	}
	return t
}

func main() {

	db, err := sql.Open("sqlite3", "todo.sqlite")
	if err != nil {
		panic(err)
	}
	todoRepository := todo.NewSQLiteRepository(db)
	if err := todoRepository.Migrate(); err != nil {
		panic(err)
	}

	funcMap := template.FuncMap{
		"hasItems": func(todo []todo.ToDo) bool {
			return len(todo) > 0
		},

		"div": func(a time.Duration, b time.Duration) int64 {
			return int64(a / b)
		},
	}

	index_template, err := template.New("index.html").Funcs(funcMap).ParseGlob("templates/*.html")
	if err != nil {
		panic(err.Error())
	}

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Received index request\n")
		todos, err := todoRepository.All()
		if err != nil {
			panic(err)
		}
		if err := index_template.Execute(w, todos); err != nil {
			panic(err)
		}
	})

	mux.HandleFunc("/addTodo", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Received new todo request\n")

		timeEstimate, err := time.ParseDuration(r.PostFormValue("todoTime") + r.PostFormValue("todoUnit"))
		if err != nil {
			panic(err)
		}
		newTodo := todo.ToDo{
			Title:        r.PostFormValue("todoTitle"),
			Started:      false,
			Editing:      false,
			TimeStart:    time.Time{},
			TimeElapsed:  0,
			TimeEstimate: timeEstimate}

		if _, err := todoRepository.Add(newTodo); err != nil {
			panic(err)
		}

		todos, err := todoRepository.All()
		if err != nil {
			panic(err)
		}

		if err := index_template.ExecuteTemplate(w, "display", todos); err != nil {
			panic(err)
		}
	})

	mux.HandleFunc("/delete/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Received delete todo request\n")

		split_path := strings.Split(r.URL.Path, "/")
		idx, err := strconv.Atoi(split_path[len(split_path)-1])
		if err != nil {
			panic(err)
		}

		if err := todoRepository.Delete(int64(idx)); err != nil {
			panic(err)
		}

		todos, err := todoRepository.All()
		if err != nil {
			panic(err)
		}

		if err := index_template.ExecuteTemplate(w, "display", todos); err != nil {
			panic(err)
		}
	})

	mux.HandleFunc("/editTodo", func(w http.ResponseWriter, r *http.Request) {

		selectedTodo := GetID(r.URL.Query().Get("id"))
		fmt.Printf("Received an edit request for %d\n", selectedTodo)

		newTitle := r.PostFormValue("editTitle")

		var newAmount time.Duration
		if len(r.PostFormValue("editAmount")) > 0 {
			newAmount, err = time.ParseDuration(r.PostFormValue("editAmount"))
			if err != nil {
				panic(err)
			}
		}

		todos, err := todoRepository.All()
		if err != nil {
			panic(err)
		}

		for i := range todos {
			todo := &todos[i]
			if todo.ID == selectedTodo {
				if !todo.Editing {
					todo.Editing = true
					todo, err = todoRepository.Update(selectedTodo, *todo)
					if err != nil {
						panic(err)
					}

					if err := index_template.ExecuteTemplate(w, "editTodo", todo); err != nil {
						panic(err)
					}
				} else {
					fmt.Printf("new title: %s\n", newTitle)
					todo.Title = newTitle
					todo.TimeEstimate = newAmount
					todo.Editing = false

					todo, err = todoRepository.Update(selectedTodo, *todo)
					if err != nil {
						panic(err)
					}

					if err := index_template.ExecuteTemplate(w, "todoTimer", todo); err != nil {
						panic(err)
					}
				}
				break
			}
		}
	})

	mux.HandleFunc("/timeIt", func(w http.ResponseWriter, r *http.Request) {

		selectedTodo := GetID(r.URL.Query().Get("id"))
		fmt.Printf("Received a todo request for %d\n", selectedTodo)

		todo, err := todoRepository.GetByID(selectedTodo)
		if err != nil {
			panic(err)
		}

		if todo.Started {
			todoRepository.StopTiming(todo)
		} else {
			todoRepository.StartTiming(todo)
		}

		if err := index_template.ExecuteTemplate(w, "todoTimer", todo); err != nil {
			panic(err)
		}
	})

	err = http.ListenAndServe(":1234", mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
