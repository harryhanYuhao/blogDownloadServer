package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"

	"github.com/harryhanYuhao/blogDownloadServer/crypto"
)

func syncBlog(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	idsum := sha256.Sum256([]byte(id))
	if crypto.Key == fmt.Sprintf("%x", idsum) {
		w.Write([]byte("Authorized"))
	} else {
		w.Write([]byte("Unauthorized"))
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/sync/", syncBlog)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
