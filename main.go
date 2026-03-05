package main

import (
	"bleh/services"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	ID        int    `json:"id"`
	OwnerID   int    `json:"owner_id"` // expose API response
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	Time      int64  `json:"time"`
}

// implement db!
var tasks = []Task{}
var nextID = 1

func main() {
	tasks = append(tasks,
		Task{ID: nextID, OwnerID: 1, Title: "cat's secret note", Completed: false, Time: time.Now().Unix()},
	)

	mux := http.NewServeMux()
	loginfs := http.FileServer(http.Dir("./src/pages/login"))
	homefs := http.FileServer(http.Dir("./src/pages/home"))
	mux.Handle("/", homefs)
	mux.Handle("/login/", http.StripPrefix("/login/", loginfs))
	mux.HandleFunc("/api/login", services.HandleLogin)
	mux.HandleFunc("/tasks", handleTasks)
	mux.HandleFunc("/tasks/", handleTaskbyID)
	protectedHandler := authMiddle(noCache(mux))
	http.ListenAndServe(":8080", protectedHandler)
}

func noCache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate, max-age=0")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		h.ServeHTTP(w, r)
	})
}

func getSessionUserID(r *http.Request) (int, bool) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return 0, false
	}
	sessionData, exists := services.Sessions[cookie.Value]
	if !exists || time.Now().After(sessionData.ExpiresAt) {
		return 0, false
	}
	return sessionData.UserID, true
}

func authMiddle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/login/" || strings.HasPrefix(r.URL.Path, "/login/") || r.URL.Path == "/api/login" {
			next.ServeHTTP(w, r)
			return
		}
		session, err := r.Cookie("session_id")
		if err != nil {
			http.Redirect(w, r, "/login/", http.StatusFound)
			return
		}
		sessionData, exists := services.Sessions[session.Value]
		if !exists || time.Now().After(sessionData.ExpiresAt) {
			http.Redirect(w, r, "/login/", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func handleTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no cache")

	callerID, _ := getSessionUserID(r)

	if r.Method == http.MethodGet {
		for _, t := range tasks {
			if t.OwnerID != callerID {
				printIDORWarning("GET /tasks", callerID, t.OwnerID, t.ID,
					"user is reading another user's task with no authorization check")
			}
		}
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
		if newTask.OwnerID != callerID && newTask.OwnerID != 0 {
			printIDORWarning("POST /tasks", callerID, newTask.OwnerID, newTask.ID,
				"user is creating a task with a spoofed owner_id")
		}
		if newTask.OwnerID == 0 {
			newTask.OwnerID = callerID
		}
		newTask.ID = nextID
		newTask.Time = time.Now().Unix()
		// A user can set owner_id to any value and claim ownership
		// of tasks on behalf of other users.
		if newTask.OwnerID != callerID && newTask.OwnerID != 0 {
			printIDORWarning("POST /tasks", callerID, newTask.OwnerID, newTask.ID,
				"user is creating a task with a spoofed owner_id")
		}
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

	callerID, _ := getSessionUserID(r)

	if r.Method == http.MethodPut {
		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/tasks/"))

		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
		}

		found := false
		for i := 0; i < len(tasks); i++ {
			if tasks[i].ID == id {
				// No ownership check — any user can modify any task by ID.
				if tasks[i].OwnerID != callerID {
					printIDORWarning("PUT /tasks/"+strconv.Itoa(id), callerID, tasks[i].OwnerID, id,
						"user is modifying a task they do not own")
				}
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
		}
		if !found {
			http.Error(w, "task not found", http.StatusNotFound)
		}
	} else if r.Method == http.MethodDelete {
		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/tasks/"))
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		for i := 0; i < len(tasks); i++ {
			if tasks[i].ID == id {
				// No ownership check — any user can delete any task by ID.
				if tasks[i].OwnerID != callerID {
					printIDORWarning("DELETE /tasks/"+strconv.Itoa(id), callerID, tasks[i].OwnerID, id,
						"user is deleting a task they do not own")
				}
				tasks = append(tasks[:i], tasks[i+1:]...)
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}

		http.Error(w, "task not found", http.StatusNotFound)
	}

}

func printIDORWarning(endpoint string, calledID, ownerID, taskID int, reason string) {
	fmt.Println()
	fmt.Println("!! IDOR VULNERABILITY TRIGGERED !!")
	fmt.Printf("endpoint: %-39s\n", endpoint)
	fmt.Printf("Called ID: %-39d\n", calledID)
	fmt.Printf("owner ID: %-39d\n", ownerID)
	fmt.Printf("task ID: %-39d\n", taskID)
	fmt.Printf("reason : %-39s\n", reason)
}
