package main

import (
	"log"
	"net/http"
	"os"
	"project/internal/database"
	"project/internal/handlers"
	"strings" // Добавили strings для проверки ID в URL
)

func main() {
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		databaseUrl = "postgres://taskuser:taskpass@localhost:5432/taskdb?sslmode=disable"
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	log.Printf("Starting server on port... %s", serverPort)

	db, err := database.Connect(databaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	log.Println("Успешно подключено к бд")

	taskStore := database.NewTaskStore(db)

	handler := handlers.NewHandlers(taskStore)

	mux := http.NewServeMux()

	// Раздаём фронтенд из корня проекта (запускай: go run ./cmd/api из корня)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "index.html")
	})

	// 1. Оставляем роут для создания задачи
	mux.HandleFunc("/tasks/create", methodHandlear(handler.CreateTask, "POST"))

	// 2. Делаем этот роут единой точкой входа для всех остальных запросов на /tasks/...
	// Он больше не конфликтует, так как зарегистрирован всего ОДИН раз
	mux.HandleFunc("/tasks/", taskIdHandler(handler))

	loggedMux := loggingMiddleware(corsMiddleware(mux))

	serverAddr := ":" + serverPort

	err = http.ListenAndServe(serverAddr, loggedMux)
	if err != nil {
		log.Fatal(err)
	}
}

func methodHandlear(handlerFunc http.HandlerFunc, allowedMethods string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != allowedMethods {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return // Важно: останавливаем выполнение, если метод не тот
		}
		handlerFunc(w, r)
	}
}

func taskIdHandler(handler *handlers.Handlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Отрезаем префикс. Если пришел запрос на "/tasks/", останется "" (пустая строка).
		// Если пришел запрос на "/tasks/15", останется "15".
		idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")

		switch r.Method {
		case http.MethodGet:
			if idStr == "" {
				// Если ID пустой, значит клиент хочет получить ВСЕ задачи
				handler.GetAllTasks(w, r)
			} else {
				// Если ID есть (например, "15"), получаем конкретную задачу
				handler.GetTask(w, r)
			}
		case http.MethodPut:
			handler.UpdateTask(w, r)
		case http.MethodPost:
			handler.UpdateTask(w, r)
		case http.MethodDelete:
			handler.DeleteTask(w, r)
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

// corsMiddleware позволяет открывать index.html отдельно от API (Live Server / file://)
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
