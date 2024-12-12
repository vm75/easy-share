package server

import (
	"encoding/json"
	"easy-share/utils"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
)

var (
	staticDir        = "./static"
	nwChangedChannel = make(chan string) // Channel for sending status updates
)

type ModuleStatus struct {
	Running bool                   `json:"running"`
	Info    map[string]interface{} `json:"info"`
}

func queryParams(r *http.Request) map[string]string {
	params := make(map[string]string)
	for k, v := range r.URL.Query() {
		if len(v) == 0 {
			continue
		}
		params[k] = v[0]
	}
	return params
}

func getStatus() map[string]interface{} {
	status := make(map[string]interface{})
	status["config"] = GetConfig()
	status["nmbdRunning"] = IsNmbdRunning()
	status["smbdRunning"] = IsSmbdRunning()

	return status
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create a channel to close the connection on client disconnect
	clientGone := r.Context().Done()

	data, _ := json.Marshal(getStatus())
	fmt.Fprintf(w, "data: %s\n\n", data)
	w.(http.Flusher).Flush() // Ensure data is sent immediately

	for {
		select {
		case event := <-nwChangedChannel: // Receive new status from the channel
			utils.LogLn("Received event:", event)
			data, _ := json.Marshal(getStatus())
			fmt.Fprintf(w, "data: %s\n\n", data)
			w.(http.Flusher).Flush() // Ensure data is sent immediately
		case <-clientGone: // Client disconnected
			return
		}
	}
}

// Enable nmbd
func enableNmbdHandler(w http.ResponseWriter, r *http.Request) {
	err := EnableNmbd()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

// Disable nmbd
func disableNmbdHandler(w http.ResponseWriter, r *http.Request) {
	err := DisableNmbd()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func addUserHandler(w http.ResponseWriter, r *http.Request) {
	var config UserConfig
	err := utils.GetContent(r, &config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = AddUser(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func delUserHandler(w http.ResponseWriter, r *http.Request) {
	var config UserConfig
	err := utils.GetContent(r, &config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = DelUser(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func addSambaShareHandler(w http.ResponseWriter, r *http.Request) {
	var config SambaShareConfig
	err := utils.GetContent(r, &config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = AddSambaShare(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func delSambaShareHandler(w http.ResponseWriter, r *http.Request) {
	var config SambaShareConfig
	err := utils.GetContent(r, &config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = DelSambaShare(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func addNfsShareHandler(w http.ResponseWriter, r *http.Request) {
	var config NfsShareConfig
	err := utils.GetContent(r, &config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = AddNfsShare(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func delNfsShareHandler(w http.ResponseWriter, r *http.Request) {
	var config NfsShareConfig
	err := utils.GetContent(r, &config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = DelNfsShare(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type FileInfo struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	IsDir bool   `json:"isDir"`
}

// Helper function to sanitize paths
func sanitizePath(inputPath string) string {
	// Prevent navigating outside baseDir
	cleanPath := filepath.Clean("/" + inputPath) // Ensure the path starts with a slash
	return filepath.Join(VarDir, cleanPath)
}

func listFilesHandler(w http.ResponseWriter, r *http.Request) {
	relPath := r.URL.Query().Get("path")
	absPath := sanitizePath(relPath)

	files, err := os.ReadDir(absPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var fileInfos []FileInfo
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".auth") {
			continue
		}
		fileInfos = append(fileInfos, FileInfo{
			Name:  file.Name(),
			Path:  filepath.Join(relPath, file.Name()), // Preserve the relative path for the client
			IsDir: file.IsDir(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fileInfos)
}

func fileContentHandler(w http.ResponseWriter, r *http.Request) {
	relPath := r.URL.Query().Get("path")
	absPath := sanitizePath(relPath)

	// Ensure the requested path is inside the base directory
	if !strings.HasPrefix(absPath, VarDir) || strings.HasSuffix(absPath, ".auth") {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	content, err := os.ReadFile(absPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(content)
}

// Separate function to handle static files
func handleStaticFiles(r *mux.Router) {
	// Serve static files from /static and root (/)
	fs := http.FileServer(http.Dir(staticDir))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	r.PathPrefix("/").Handler(http.StripPrefix("/", fs)) // Serve "/" from staticDir
}

func WebServer(port string) {
	// Create a new Gorilla Mux router
	r := mux.NewRouter()

	// Config-related routes
	r.HandleFunc("/api/status", statusHandler).Methods("GET")

	// Handle file
	r.HandleFunc("/api/files", listFilesHandler)
	r.HandleFunc("/api/file", fileContentHandler)

	// Module
	r.HandleFunc("/api/nmbd/enable", enableNmbdHandler).Methods("POST")
	r.HandleFunc("/api/nmbd/disable", disableNmbdHandler).Methods("POST")
	r.HandleFunc("/api/user/add", addUserHandler).Methods("POST")
	r.HandleFunc("/api/user/del", delUserHandler).Methods("POST")
	r.HandleFunc("/api/samba/add", addSambaShareHandler).Methods("POST")
	r.HandleFunc("/api/samba/del", delSambaShareHandler).Methods("POST")
	r.HandleFunc("/api/nfs/add", addNfsShareHandler).Methods("POST")
	r.HandleFunc("/api/nfs/del", delNfsShareHandler).Methods("POST")

	// Serve static files
	handleStaticFiles(r)

	// Start the server
	utils.LogF("Server starting on port %s\n", port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		utils.LogFatal(err)
	}
}
