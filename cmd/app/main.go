package main

import (
	"WIND/internal/database"
	"WIND/internal/handlers"
	"WIND/internal/taskService"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	database.InitDB()
	database.DB.AutoMigrate(&taskService.Task{})

	repo := taskService.NewTaskRepository(database.DB)
	service := taskService.NewService(repo)

	handler := handlers.NewHandler(service)

	router := mux.NewRouter()
	router.HandleFunc("/api/task", handler.GetTasksHandler).Methods("GET")
	router.HandleFunc("/api/task", handler.CreateTaskHandler).Methods("POST")
	router.HandleFunc("/api/task/{id}", handler.UpdateTaskHandler).Methods("PATCH")
	router.HandleFunc("/api/task/{id}", handler.DeleteTaskHandler).Methods("DELETE")
	http.ListenAndServe(":8080", router)
}
