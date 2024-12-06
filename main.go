package main

import (
	"fmt"
	"path/filepath"
	"net/http"
	"os"
	"encoding/json"
	"strings"
)

type Response struct {
	Message string `json:"message"`
}

type ChampionResponse struct{
	Name string `json:"name"`
	URL  string `json:"url"`
}

func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}

func extractChampionName(path string) string {
	return strings.TrimPrefix(path, "/api/champion/")
}

func getChampionImageURL(championName string) string {
	return filepath.Join("./static/images", championName+".png")
}

func imageExists(imagePath string) bool {
	_, err := os.Stat(imagePath)
	return err == nil
}

func handleChampionImage(w http.ResponseWriter, r *http.Request) {
	championName := extractChampionName(r.URL.Path)

	if championName == "" {
		http.Error(w, "Champion name is required to hit this API", http.StatusBadRequest)
		return
	}
	imagePath := getChampionImageURL(championName)
	if !imageExists(imagePath) {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, imagePath)
}

func testAPIfunction(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Test API function hit")
	w.Header().Set("Content-Type", "application/json")
	response := Response{Message: "Hello World from GO"}
	json.NewEncoder(w).Encode(response)
}

func setupRoutes(mux *http.ServeMux) {
	// Serve static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("/api/champion/", handleChampionImage)
	mux.HandleFunc("/testAPI", testAPIfunction)
}

func main() {
	// Create a new router (mux)
	mux := http.NewServeMux()

	// Set up routes
	setupRoutes(mux)

	// Apply CORS middleware to all routes by wrapping the mux handler
	handlerWithCors := enableCors(mux)

	// Start the server with CORS applied
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", handlerWithCors)
}
