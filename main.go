package main

import (
	"html/template"
	"net/http"
	"strconv"
)

type Task struct {
	ID        int
	Text      string
	Completed bool
}

var tasks = make(map[int]Task)
var idCounter int

func main() {

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/delete", deleteHandler)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Запуск сервера
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/index.html",
	))

	taskSlice := make([]Task, 0, len(tasks))
	for _, task := range tasks {
		taskSlice = append(taskSlice, task)
	}

	tmpl.Execute(w, struct{ Tasks []Task }{taskSlice})
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	text := r.FormValue("task")
	idCounter++

	tasks[idCounter] = Task{
		ID:        idCounter,
		Text:      text,
		Completed: false,
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	delete(tasks, id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
