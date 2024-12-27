package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

var db *gorm.DB

func initDB() {
	dsn := "host=localhost user=postgres password=ali24815 dbname=chat-app port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	db.AutoMigrate(&User{})
}

func jsonResponse(w http.ResponseWriter, status int, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

func handleGetUsers(w http.ResponseWriter, r *http.Request) {
	var users []User
	db.Find(&users)
	jsonResponse(w, http.StatusOK, Response{
		Status:  "success",
		Message: "User list fetched successfully",
		Data:    users,
	})
}

func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		jsonResponse(w, http.StatusBadRequest, Response{
			Status:  "fail",
			Message: "Invalid JSON format",
		})
		return
	}
	if user.Name == "" || user.Email == "" {
		jsonResponse(w, http.StatusBadRequest, Response{
			Status:  "fail",
			Message: "Name and Email are required",
		})
		return
	}
	db.Create(&user)
	jsonResponse(w, http.StatusCreated, Response{
		Status:  "success",
		Message: "User created successfully",
		Data:    user,
	})
}

func handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, Response{
			Status:  "fail",
			Message: "Invalid ID",
		})
		return
	}
	var user User
	if err := db.First(&user, id).Error; err != nil {
		jsonResponse(w, http.StatusNotFound, Response{
			Status:  "fail",
			Message: "User not found",
		})
		return
	}
	var updatedData User
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		jsonResponse(w, http.StatusBadRequest, Response{
			Status:  "fail",
			Message: "Invalid JSON format",
		})
		return
	}
	db.Model(&user).Updates(updatedData)
	jsonResponse(w, http.StatusOK, Response{
		Status:  "success",
		Message: "User updated successfully",
		Data:    user,
	})
}

func handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, Response{
			Status:  "fail",
			Message: "Invalid ID",
		})
		return
	}
	if err := db.Delete(&User{}, id).Error; err != nil {
		jsonResponse(w, http.StatusNotFound, Response{
			Status:  "fail",
			Message: "User not found",
		})
		return
	}
	jsonResponse(w, http.StatusOK, Response{
		Status:  "success",
		Message: "User deleted successfully",
	})
}

func main() {
	initDB()
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetUsers(w, r)
		case http.MethodPost:
			handleCreateUser(w, r)
		case http.MethodPut:
			handleUpdateUser(w, r)
		case http.MethodDelete:
			handleDeleteUser(w, r)
		default:
			jsonResponse(w, http.StatusMethodNotAllowed, Response{
				Status:  "fail",
				Message: "Method not allowed",
			})
		}
	})
	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
