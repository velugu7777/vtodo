package main

import (
	"encoding/json"
	"fmt"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"net/http"
)

type Todo struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	IsDone bool   `json:"isDone"`
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

}

func delTodo(c web.C, w http.ResponseWriter, req *http.Request) {

}
