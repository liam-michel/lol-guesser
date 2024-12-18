package main

import (
	"database/sql"
	"fmt"
	"log"
	"lol-guesser/auth"
	"lol-guesser/lol_data"
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

// Global DB variable

func InitDB() (*sql.DB, error) {
	err := godotenv.Load("../.env")
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	dbPassword := os.Getenv("MYSQLPASSWORD")
	dbName := "lol_users"
	dsn := fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/%s", dbPassword, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging the database: %v", err)
	}

	fmt.Println("Successfully connected to the database!")
	return db, nil
}

func setupRoutes(mux *http.ServeMux) {
	// Serve static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../static"))))
	mux.HandleFunc("/api/randomchampion", lol_data.GetRandomChampionHandler)
	mux.HandleFunc("/api/getuser", auth.GetUserHandler)
	mux.HandleFunc("/api/createuser", auth.CreateUserHandler)
	mux.HandleFunc("/api/login", auth.LoginHandler)
	mux.HandleFunc("/api/refresh", auth.RefreshTokenHandler)
}

func main() {

	//load env variables
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file in main.go")
	}
	db, err := InitDB()
	if err != nil {
		log.Fatalf("Error initializing the database: %v", err)
	}
	auth.SetDB(db)
	defer db.Close()
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
