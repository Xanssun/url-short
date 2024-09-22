package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
)

var urlHashMap = make(map[string]string)

func generateHashUrl(url string) string {
	h := sha256.New()
	h.Write([]byte(url))
	hash := fmt.Sprintf("%x", h.Sum(nil))
	return hash[:8]
}

func shortUrl(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Invalid request method", http.StatusBadRequest)
		return
	}

	// Чтение данных из тела запроса
	responseData, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// Преобразование тела запроса в строку, потому что получаем в байтах
	responseString := string(responseData)

	// Генерация хеша URL
	urlHash := generateHashUrl(responseString)

	// дабавляет в hash map
	if _, ok := urlHashMap[responseString]; !ok {
		urlHashMap[urlHash] = responseString
	}

	// Возвращаем сокращенный URL
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)

	_, err = res.Write([]byte("http://localhost:8080/" + urlHash))
	if err != nil {
		log.Println("Error writing response:", err)
	}
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", shortUrl)

	fmt.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
