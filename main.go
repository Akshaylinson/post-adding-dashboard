package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

type Content struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	Image string `json:"image"`
	Link  string `json:"link"`
}

var (
	dataFile = "data.json"
	mu       sync.Mutex
)

// Load content from JSON file
func loadContent() Content {
	file, err := os.ReadFile(dataFile)
	if err != nil {
		return Content{}
	}
	var content Content
	json.Unmarshal(file, &content)
	return content
}

// Save content to JSON file
func saveContent(content Content) {
	mu.Lock()
	defer mu.Unlock()
	data, _ := json.MarshalIndent(content, "", "  ")
	os.WriteFile(dataFile, data, 0644)
}

// Serve Admin page
func adminPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/admin.html")
}

// Serve User page
func userPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/user.html")
}

// API to get content
func getContent(w http.ResponseWriter, r *http.Request) {
	content := loadContent()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(content)
}

// API to update content (admin only)
func updateContent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var content Content
	err := json.NewDecoder(r.Body).Decode(&content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	saveContent(content)
	w.Write([]byte("Content updated successfully"))
}

func main() {
	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes
	http.HandleFunc("/admin", adminPage)
	http.HandleFunc("/user", userPage)
	http.HandleFunc("/api/get", getContent)
	http.HandleFunc("/api/update", updateContent)

	port := ":8080"
	fmt.Println("🚀 Server running at http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, nil))
}
