package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var tasks = []Task{}
var nextID = 1

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/tasks", handleTasks)
	http.HandleFunc("/tasks/", handleTaskbyID)
	http.ListenAndServe(":8080", nil)

}

func handleTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodGet {
		json.NewEncoder(w).Encode(tasks)
	} else if r.Method == http.MethodPost {
		var newTask Task
		err := json.NewDecoder(r.Body).Decode(&newTask)
		if err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if len(newTask.Title) > 200 || newTask.Title == "" {
			http.Error(w, "Title too large or empty", http.StatusBadRequest)
			return
		}
		newTask.ID = nextID
		nextID++
		newTask.Completed = false
		tasks = append(tasks, newTask)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newTask)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleTaskbyID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodPut {
		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/tasks/"))
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
		}
		found := false
		for i := 0; i < len(tasks); i++ {
			if tasks[i].ID == id {
				found = true
				var newTask Task
				err := json.NewDecoder(r.Body).Decode(&newTask)
				if err != nil {
					http.Error(w, "invalid json", http.StatusBadRequest)
					return
				}
				tasks[i].Title = newTask.Title
				tasks[i].Completed = newTask.Completed
				json.NewEncoder(w).Encode(tasks[i])
				return
			}
			if !found {
				http.Error(w, "task not found", http.StatusNotFound)
			}
		}
	} else if r.Method == http.MethodDelete {

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
