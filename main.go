package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	Time      int64  `json:"time"`
}

// implement db!
var tasks = []Task{}
var nextID = 1

func main() {
	fs := http.FileServer(http.Dir("./src"))
	http.Handle("/", noCache(fs))
	http.HandleFunc("/tasks", handleTasks)
	http.HandleFunc("/tasks/", handleTaskbyID)
	http.ListenAndServe(":8080", nil)
}

func noCache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate, max-age=0")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		h.ServeHTTP(w, r)
	})
}

func handleTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no cache")
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
		newTask.Time = time.Now().Unix()
		nextID++
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
		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/tasks/"))
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		for i := 0; i < len(tasks); i++ {
			if tasks[i].ID == id {
				tasks = append(tasks[:i], tasks[i+1:]...)
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}

		http.Error(w, "task not found", http.StatusNotFound)
	}

}
