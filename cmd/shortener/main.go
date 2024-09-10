package main

import (
	"crypto/sha256"
	"fmt"
	"net/http"
)

var urlHash = make(map[string]string)

func HashUrl(url string) string {
	hash := sha256.New()
	hash.Write([]byte(url))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func main() {
	http.ListenAndServe(":8080", nil)

}
