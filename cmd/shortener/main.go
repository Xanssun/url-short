package main

import (
	"crypto/sha256"
	"fmt"
	"net/http"
)

var urlHash = make(map[string]string)

// Функция для хеширования URL
func HashUrl(url string) string {
	hash := sha256.New()
	hash.Write([]byte(url))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// Проверка, существует ли хеш
func isValidHash(hash string) bool {
	_, ok := urlHash[hash]
	return ok
}

// Генерация короткого URL
func ShortUrl(url string) string {
	hash := HashUrl(url)
	if !isValidHash(hash) {
		urlHash[hash] = url
	}
	return fmt.Sprintf("http://localhost:8080/%s", hash)
}

// Обработчик запросов
func handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	url := r.FormValue("url")
	if url == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	shortened := ShortUrl(url)
	fmt.Fprintf(w, "Shortened URL: %s\n", shortened)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", handleShorten)

	fmt.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
