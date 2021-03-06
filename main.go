package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"log"
	"net/http"
	"os"
	"strconv"
)

var db *sql.DB
var logger = log.New(os.Stdout, "vtodo:", log.Ldate+log.Ltime+log.Lshortfile)

type Todo struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	IsCompleted bool   `json:"isCompleted"`
}

type EmberTodo struct {
	Todo Todo `json:"todo"`
}

type EmberTodos struct {
	Todos []Todo `json:"todos"`
}

func init() {
	var err error
	db, err = sql.Open("postgres", "user=postgres password =Sunithak*1247 host=localhost port=8080 dbname=postgres sslmode=disable")
	if err != nil {
		logger.Println("Open connection failed", err.Error())
	}
}

func main() {
	goji.Get("/", index)
	goji.Get("/assets/*", http.FileServer(http.Dir("./dist")))
	goji.Get("/api/todos", todos)
	goji.Post("/api/todos", newTodo)
	goji.Put("/api/todos/:id", putTodo)
	goji.Delete("/api/todos/:id", delTodo)

	goji.NotFound(index)

	goji.Serve()
}

func index(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "./dist/index.html")
}

func todos(w http.ResponseWriter, req *http.Request) {
	var mytodos []Todo
	stmt, err := db.Prepare("Select id, name, isCompleted from todos")
	if err != nil {
		logger.Println("Prepare failed:", err.Error())
		return
	}

	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		logger.Println(err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		t := Todo{}
		err := rows.Scan(&t.Id, &t.Name, &t.IsCompleted)
		if err != nil {
			logger.Println(err)
		}
		mytodos = append(mytodos, t)
	}
	err = rows.Err()
	if err != nil {
		logger.Println(err)
		return
	}

	etodos := EmberTodos{
		Todos: mytodos,
	}
	j, err := json.Marshal(etodos)
	if err != nil {
		logger.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
	return
}

func newTodo(w http.ResponseWriter, req *http.Request) {
	var todo EmberTodo
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&todo)
	if err != nil {
		logger.Printf("json decode failed:%s", err.Error())
		err := fmt.Errorf("Error")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	stmt, err := db.Prepare("INSERT INTO todos(name,isCompleted) values($1,$2) RETURNING id")
	if err != nil {
		logger.Println("Prepare failed:", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(todo.Todo.Name, todo.Todo.IsCompleted).Scan(&todo.Todo.Id)
	if err != nil {
		logger.Println("Insert failed:", err.Error())
		err := fmt.Errorf("Error")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	j, err := json.Marshal(todo)
	if err != nil {
		logger.Println(err)
		err := fmt.Errorf("Error")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
	return
}

func putTodo(c web.C, w http.ResponseWriter, req *http.Request) {
	id := c.URLParams["id"]
	var todo EmberTodo
	var err error
	todo.Todo.Id, err = strconv.Atoi(id)

	if err != nil {
		logger.Printf("Json decode failed:%s", err.Error())
		err := fmt.Errorf("Error")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&todo)
	if err != nil {
		logger.Printf("Json decode failed:%s", err.Error())
		err := fmt.Errorf("Error")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("UPDATE todos SET Name= $1, IsCompleted = $2 where id = $3")
	if err != nil {
		logger.Println("Prepared failed:", err.Error())
		err := fmt.Errorf("Error")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(todo.Todo.Name, todo.Todo.IsCompleted, id)
	if err != nil {
		logger.Println("Update failed", err.Error())
		err := fmt.Errorf("Error")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	affect, err := res.RowsAffected()
	if err != nil {
		logger.Println("Rows affected failed", err.Error())
		err := fmt.Errorf("Error")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Println(affect, "rows affected")

	j, err := json.Marshal(todo)
	if err != nil {
		logger.Println(err)
		err := fmt.Errorf("Error")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
	return
}

func delTodo(c web.C, w http.ResponseWriter, req *http.Request) {
	id := c.URLParams["id"]
	logger.Println(id)
	stmt, err := db.Prepare("DELETE from todos where id = $1")
	if err != nil {
		logger.Println(err)
		return
	}
	res, err := stmt.Exec(id)
	if err != nil {
		logger.Println(err)
		return
	}
	logger.Println(res)
}
