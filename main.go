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
)

var db *sql.DB

type Todo struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	IsCompleted bool   `json:"isCompleted"`
}

func init() {
	var err error
	db, err = sql.Open("postgres", "user=postgres password=Sunithak*1247 host=localhost dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
}

type EmberTodo struct {
	Todo Todo `json:"todo"`
}

type EmberTodos struct {
	Todos []Todo `json:"todos"`
}

var mytodos = []Todo{
	{1, "reading books", false},
	{2, "playing cricket", false},
}

func main() {
	goji.Get("/", index)
	goji.Get("/assets/*", http.FileServer(http.Dir("./dist")))
	goji.Get("/api/todos", todos)
	goji.Post("/api/todos/newTodo", newTodo)
	goji.Put("/api/todos/:id", putTodo)
	goji.Delete("/api/todos/:id", delTodo)
	goji.NotFound(index)
	goji.Serve()
}

func index(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "./dist/index.html")
}

func todos(w http.ResponseWriter, req *http.Request) {
	etodos := EmberTodos{
		Todos: mytodos,
	}
	j, err := json.Marshal(etodos)
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-type", "application/json")
	w.Write(j)
}

func newTodo(w http.ResponseWriter, req *http.Request) {

}

func putTodo(c web.C, w http.ResponseWriter, req *http.Request) {
	var etodo EmberTodo
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&etodo)
	if err != nil {
		log.Println("JSON decode failed:", err.Error())
	}
	id, err := strconv.ParseInt(c.URLParams["id"], 10, 64)
	if err != nil {
		log.Println("URLParams[id] failed:", err.Error())
	}
	etodo.Todo.Id = int(id)
	fmt.Println(etodo)

}

func delTodo(c web.C, w http.ResponseWriter, req *http.Request) {

}
