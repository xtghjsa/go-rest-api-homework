package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Обработчик для отправки задачи на сервер
func postTask(a http.ResponseWriter, b *http.Request) {
	var task Task
	var buf bytes.Buffer
	_, err := buf.ReadFrom(b.Body)
	if err != nil {
		http.Error(a, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(a, err.Error(), http.StatusBadRequest)
		return
	}
	tasks[task.ID] = task
	a.Header().Set("Content-Type", "application/json")
	a.WriteHeader(http.StatusCreated)
}

// Обработчик для получения всех задач
func getTasks(a http.ResponseWriter, b *http.Request) {
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(a, err.Error(), http.StatusInternalServerError)
		return
	}
	a.Header().Set("Content-Type", "application/json")
	a.WriteHeader(http.StatusOK)
	a.Write(resp)
}

// Обработчик для получения задачи по ID
func getTaskId(a http.ResponseWriter, b *http.Request) {
	id := chi.URLParam(b, "id")
	task, ok := tasks[id]
	if !ok {
		http.Error(a, "Task not found", http.StatusBadRequest)
	}
	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(a, err.Error(), http.StatusBadRequest)
		return
	}
	a.Header().Set("Content-Type", "application/json")
	a.WriteHeader(http.StatusOK)
	a.Write(resp)
}

// Обработчик удаления задачи по ID
func delTask(a http.ResponseWriter, b *http.Request) {
	id := chi.URLParam(b, "id")
	_, ok := tasks[id]
	if !ok {
		http.Error(a, "", http.StatusBadRequest)
	}
	delete(tasks, id)
	a.WriteHeader(http.StatusOK)
}
func main() {
	r := chi.NewRouter()
	r.Post("/tasks", postTask)
	r.Get("/tasks", getTasks)
	r.Get("/tasks/{id}", getTaskId)
	r.Delete("/tasks/{id}", delTask)
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
