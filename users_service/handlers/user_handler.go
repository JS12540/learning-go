package handlers

import (
	"encoding/json"
	"net/http"
	"users_service/models"
)

var users = []models.User{
	{ID: "1", Name: "Jay", Age: 25},
	{ID: "2", Name: "Rahul", Age: 30},
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
