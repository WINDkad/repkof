package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type requestBody struct {
	Task   string `json:"task"`
	IsDone bool   `json:"is_done"`
}

type responseBody struct {
	ID     uint   `json:"ID"`
	Task   string `json:"task"`
	IsDone bool   `json:"is_done"`
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	var messages []Message

	if err := DB.Find(&messages).Error; err != nil {
		http.Error(w, "Failed to get message", http.StatusInternalServerError)
		return
	}

	var response []responseBody

	for _, msg := range messages {
		response = append(response, responseBody{
			ID:     msg.ID,
			Task:   msg.Task,
			IsDone: msg.IsDone})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	var message requestBody

	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	task := Message{Task: message.Task, IsDone: message.IsDone}
	if err := DB.Create(&task).Error; err != nil {
		http.Error(w, "Failed to create message", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseBody{
		ID:     task.ID,
		Task:   task.Task,
		IsDone: task.IsDone,
	})
}

func main() {
	InitDB()

	DB.AutoMigrate(&Message{})

	router := mux.NewRouter()
	router.HandleFunc("/api/messages", CreateMessage).Methods("POST")
	router.HandleFunc("/api/messages", GetMessages).Methods("GET")
	http.ListenAndServe(":8080", router)
}
