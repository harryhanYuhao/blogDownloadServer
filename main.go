// Description: This is the main file of the blogDownloadServer.
package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/harryhanYuhao/blogDownloadServer/crypto"
	"github.com/harryhanYuhao/blogDownloadServer/execBash"
)

const fileDir = "/tmp/blogDownloadServer"
const logFileDir = "/tmp/blogDownloadServer_log"
const logFilePath = "/tmp/blogDownloadServer_log/log"
const gitURL = "https://github.com/harryhanYuhao/blogPosts.git"

func formatedCurTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func log_writeLog(input string) {
	execBash.Execute("mkdir -p " + logFileDir)
	f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if _, err := f.WriteString(input + "\n"); err != nil {
		log.Fatal(err)
	}
}

func syncBlog() {
	_, err := os.Stat(fileDir)
	if os.IsNotExist(err) {
		// TODO: add err check
		execBash.Execute("mkdir -p " + fileDir)
		execBash.Execute("git -C " + fileDir + " clone " + gitURL + " " + fileDir)
		log_writeLog("Repo" + fileDir + " Created: " + formatedCurTime())
	} else {
		execBash.Execute("git -C " + fileDir + " pull")
		log_writeLog("Repo " + fileDir + " Updated: " + formatedCurTime())
	}
}

func handleSync(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	idsum := sha256.Sum256([]byte(id))
	// check key
	if crypto.Key == fmt.Sprintf("%x", idsum) {
		w.Write([]byte("Authorized"))
		syncBlog()
	} else {
		w.Write([]byte("Unauthorized"))
	}
}

func main() {
	syncBlog()
	mux := http.NewServeMux()
	mux.HandleFunc("/sync/", handleSync)
	mux.Handle("/download/", http.StripPrefix("/download/", http.FileServer(http.Dir(fileDir))))
	log.Fatal(http.ListenAndServe(":10001", mux))
}
