package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	//"io/ioutil"
	//"math/rand"
	//"time"
	"github.com/joho/godotenv"
)

type Response struct {
	Message string `json:"message"`
}

type ChampionResponse struct {
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

func testAPIfunction(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Test API function hit")
	w.Header().Set("Content-Type", "application/json")
	response := Response{Message: "Hello World from GO"}
	json.NewEncoder(w).Encode(response)
}

func getRandomChampion(w http.ResponseWriter, r *http.Request) {
	//read in the json names
	name, url, err := PickRandomChampion()
	fullURL := fmt.Sprintf("http://localhost:%s/static/images/%s", os.Getenv("VITE_GOLANG_PORT"), url)
	if err != nil {
		http.Error(w, "Error picking random champion", http.StatusInternalServerError)
		return
	}
	fmt.Println("Name: ", name)
	fmt.Println("URL: ", fullURL)
	w.Header().Set("Content-Type", "application/json")
	response := ChampionResponse{Name: name, URL: fullURL}
	json.NewEncoder(w).Encode(response)
}
func setupRoutes(mux *http.ServeMux) {
	// Serve static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../static"))))
	mux.HandleFunc("/api/testAPI", testAPIfunction)
	mux.HandleFunc("/api/randomchampion", getRandomChampion)
}

func main() {

	//load env variables
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	FULLPORT := ":" + os.Getenv("VITE_GOLANG_PORT")
	// Create a new router (mux)
	mux := http.NewServeMux()
	ReadChampionsJSON()

	// Set up routes
	setupRoutes(mux)

	// Apply CORS middleware to all routes by wrapping the mux handler
	handlerWithCors := enableCors(mux)

	output := fmt.Sprintf(`Server is running on port %s`, FULLPORT)
	// Start the server with CORS applied
	fmt.Println(output)
	http.ListenAndServe(FULLPORT, handlerWithCors)
}
