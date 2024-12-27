package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RequestData struct {
	Message string `json:"message"`
}

type ResponseData struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var requestData RequestData
		err := json.NewDecoder(r.Body).Decode(&requestData)
		if err != nil || requestData.Message == "" {
			response := ResponseData{
				Status:  "fail",
				Message: "Некорректное JSON-сообщение",
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		fmt.Println("Получено сообщение:", requestData.Message)

		response := ResponseData{
			Status:  "success",
			Message: "Данные успешно приняты",
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	if r.Method == http.MethodGet {
		response := ResponseData{
			Status:  "success",
			Message: "Это GET-ответ от сервера",
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func main() {
	// Обрабатываем запросы к /api с поддержкой CORS
	http.Handle("/api", enableCors(http.HandlerFunc(apiHandler)))

	// Обслуживаем статические файлы (например, index.html)
	http.Handle("/", enableCors(http.FileServer(http.Dir("."))))

	fmt.Println("Сервер запущен на порту 8080...")
	http.ListenAndServe(":8080", nil)
}
