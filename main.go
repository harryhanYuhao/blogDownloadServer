// Description: This is the main file of the blogDownloadServer.
package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/harryhanYuhao/blogDownloadServer/crypto"
	"github.com/harryhanYuhao/blogDownloadServer/execBash"
	"github.com/harryhanYuhao/blogDownloadServer/readDirRecurse"
)

const fileDir = "/tmp/blogDownloadServer/"
const logFileDir = "/tmp/blogDownloadServer_log/"
const gitURL = "https://github.com/harryhanYuhao/blogPosts.git"

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

func formatedCurTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func writeLog(path string, filename string, input string) {

	filepath := path + filename
	// check if file exist
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		execBash.Execute("mkdir -p " + path)
	}

	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Fail to Open " + filepath)
	}
	defer f.Close()

	if _, err := f.WriteString(input + "\n"); err != nil {
		fmt.Println("Fail to Write to " + filepath)
	}
}

func logSyncBlog(input string, ip string) {
	message := formatedCurTime() + " " + ip + " " + input
	writeLog(logFileDir, "sync_log", message)
}

func logDownload(input string, ip string) {
	message := formatedCurTime() + " " + ip + " " + input
	writeLog(logFileDir, "download_log", message)
}

func logList(input string, ip string) {
	message := formatedCurTime() + " " + ip + " " + input
	writeLog(logFileDir, "list_log", message)
}

func syncBlog(ip string) {
	_, err := os.Stat(fileDir)
	if os.IsNotExist(err) {
		// TODO: add err check
		execBash.Execute("mkdir -p " + fileDir)
		execBash.Execute("git -C " + fileDir + " clone " + gitURL + " " + fileDir)
		logSyncBlog("Repo "+fileDir+" Created ", ip)
	} else {
		execBash.Execute("git -C " + fileDir + " pull")
		logSyncBlog("Repo "+fileDir+" Updated ", ip)
	}
}

func handleSync(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	idsum := sha256.Sum256([]byte(id))
	// check key
	if crypto.Key == fmt.Sprintf("%x", idsum) {
		w.Write([]byte("Authorized"))
		syncBlog(ReadUserIP(r))
	} else {
		logSyncBlog("Unauthorized sync request", ReadUserIP(r))
		w.Write([]byte("Unauthorized"))
	}
}

func handleList(w http.ResponseWriter, r *http.Request) {
	list, err := readDirRecurse.ReadDirRecurse(fileDir)
	if err != nil {
		w.Write([]byte("Error"))
	}
	for i := range list {
		list[i] = list[i][len(fileDir):]
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(list)
	logList("List Request", ReadUserIP(r))
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Path[len("/download/"):]
	http.ServeFile(w, r, fileDir+filename)
	logDownload("Request for "+filename, ReadUserIP(r))
}

func main() {
	syncBlog("localhost")
	mux := http.NewServeMux()
	mux.HandleFunc("/sync/", handleSync)
	mux.HandleFunc("/list/", handleList)
	mux.HandleFunc("/download/", handleDownload)
	log.Fatal(http.ListenAndServe(":10001", mux))
}
