package main

import (
	"log"
	"net/http"
	"os"
	"project/internal/database"
	"project/internal/handlers"
)

func main() {
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		databaseUrl = "postgres://taskuser:taskpass@localhost:5432/tasksdb?sslmode=disable"
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

	mux.HandleFunc("/tasks/", methodHandlear())
}

func methodHandlear(handlerFunc http.HandlerFunc, allowedMethods string) http.HandlerFunc {

}
