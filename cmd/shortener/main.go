package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
)

var urlHashMap = make(map[string]string)

func generateHashURL(url string) string {
	h := sha256.New()
	h.Write([]byte(url))
	hash := fmt.Sprintf("%x", h.Sum(nil))
	return hash[:8]
}

func shortURL(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Invalid request method", http.StatusBadRequest)
		return
	}

	responseData, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Error reading request body", http.StatusInternalServerError)
		return
	}

	responseString := string(responseData)
	urlHash := generateHashURL(responseString)

	if _, ok := urlHashMap[urlHash]; !ok {
		urlHashMap[urlHash] = responseString
		fmt.Println("New entry added:", urlHash, responseString)
	} else {
		fmt.Println("Entry already exists:", urlHash, urlHashMap[urlHash])
	}

	// Устанавливаем заголовок Location
	res.Header().Set("Location", "http://localhost:8080/"+urlHash)
	res.WriteHeader(http.StatusCreated)
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", shortURL)

	fmt.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
