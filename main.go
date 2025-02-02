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

func UpdateMessage(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var message Message
	if err := DB.First(&message, id).Error; err != nil {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	var requestBody requestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if requestBody.Task != "" {
		message.Task = requestBody.Task
	}
	message.IsDone = requestBody.IsDone

	if err := DB.Save(&message).Error; err != nil {
		http.Error(w, "Failed to update message", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseBody{
		ID:     message.ID,
		Task:   message.Task,
		IsDone: message.IsDone,
	})
}

func DeleteMessage(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var message Message
	if err := DB.First(&message, id).Error; err != nil {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	if err := DB.Delete(&message).Error; err != nil {
		http.Error(w, "Failed to delete message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	InitDB()

	DB.AutoMigrate(&Message{})

	router := mux.NewRouter()
	router.HandleFunc("/api/messages", CreateMessage).Methods("POST")
	router.HandleFunc("/api/messages", GetMessages).Methods("GET")
	router.HandleFunc("/api/messages/{id}", UpdateMessage).Methods("PATCH")
	router.HandleFunc("/api/messages/{id}", DeleteMessage).Methods("DELETE")
	http.ListenAndServe(":8080", router)
}
