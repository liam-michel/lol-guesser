package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"lol-guesser/database"
	"lol-guesser/lol_data"

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

// func createUserHandler(w http.ResponseWriter, r *http.Request) {
// 	// Parse the request body
// 	decoder := json.NewDecoder(r.Body)
// 	var user database.User
// 	err := decoder.Decode(&user)
// 	if err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}
// 	// Add the user to the database
// 	err = database.AddUser(user)
// 	if err != nil {
// 		http.Error(w, "Error adding user to the database", http.StatusInternalServerError)
// 		return
// 	}
// 	// Return a success message
// 	response := Response{Message: "User added successfully"}
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }


func setupRoutes(mux *http.ServeMux) {
	// Serve static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../static"))))
	mux.HandleFunc("/api/randomchampion", lol_data.GetRandomChampionHandler)
	mux.HandleFunc("/api/getuser" ,database.GetUserHandler)
	mux.HandleFunc("/api/test", func(w http.ResponseWriter, r *http.Request){
		fmt.Println("Hit tester")
		w.Header().Set("Content-Type", "application/json")
		response := Response{Message: "Test success"}
		json.NewEncoder(w).Encode(response)
	})
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

	// Set up routes
	setupRoutes(mux)

	// Apply CORS middleware to all routes by wrapping the mux handler
	handlerWithCors := enableCors(mux)

	output := fmt.Sprintf(`Server is running on port %s`, FULLPORT)
	// Start the server with CORS applied
	fmt.Println(output)
	http.ListenAndServe(FULLPORT, handlerWithCors)
}
