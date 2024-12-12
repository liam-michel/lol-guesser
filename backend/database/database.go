package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type UserDetails struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRequest struct {
	Username string `json:"username"`
}

type UserResponse struct {
	Username string `json:"username"`
	Success  bool   `json:"success"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func jsonErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := map[string]string{"error": message}
	json.NewEncoder(w).Encode(response)
}

// Global DB variable
var db *sql.DB

// Function to initialize the DB connection
func InitDB() {
	var err error
	err = godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file in database.go")
	}

	// Get database credentials from environment variables or hardcode for now
	dbPassword := os.Getenv("MYSQLPASSWORD") // Replace with your real password or use os.Getenv for environment variable
	dbName := "lol_users"                    // Replace with your database name

	// DSN: Data Source Name (replace placeholders with your actual DB credentials)
	dsn := fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/%s", dbPassword, dbName)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	// Ping to make sure the connection is valid
	if err := db.Ping(); err != nil {
		log.Fatal("Error pinging the database:", err)
	}

	fmt.Println("Successfully connected to the database!")
}

func CreateUser(username string, password string) error {
	//start by checking if the user already exists in the database

	exists, err := checkUserExists(username)
	if err != nil {
		return fmt.Errorf("error getting user from database: %v", err)
	}
	if exists {

		return fmt.Errorf("user already exists in the database")
	}
	//start with hashing the password with bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %v", err)
	}
	//print the hashed password
	fmt.Println("Hashed password: ", string(hash))

	//insert the user into the database
	_, err = db.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", username, hash)
	if err != nil {
		return fmt.Errorf("error inserting user into database: %v", err)
	}
	fmt.Println("User added to database")
	return nil
}

// Function to get users from the database (example)
func GetUsers() ([]string, error) {
	rows, err := db.Query("SELECT username FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, err
		}
		users = append(users, username)
	}

	return users, nil
}

func checkUserExists(username string) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", username).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking user existence: %v", err)
	}
	return exists, nil
}

func GetUserName(username string) (string, error) {
	var retrievedUserName string
	err := db.QueryRow("SELECT username from users where username = ?", username).Scan(&retrievedUserName)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle case where no rows are returned
			return "", fmt.Errorf("user not found")
		}
		// Handle other errors
		return "", err
	}
	return retrievedUserName, nil

}

func GetUserDetails(username string) (UserDetails, error) {
	var user UserDetails
	err := db.QueryRow("SELECT username, password_hash FROM users WHERE username = ?", username).Scan(&user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return UserDetails{}, fmt.Errorf("user not found")
		}
		return UserDetails{}, err
	}
	return user, nil
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hit create user handler")
	var createUserRequest CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&createUserRequest)
	if err != nil {
		jsonErrorResponse(w, "Invalid request body, please include username and password", http.StatusBadRequest)
		return
	}
	username, password := createUserRequest.Username, createUserRequest.Password

	err = CreateUser(username, password)
	if err != nil {
		if err.Error() == "user already exists in the database" {
			jsonErrorResponse(w, "Username is already taken", http.StatusConflict)
		} else {
			log.Printf("Error adding user to the database: %v", err)
			jsonErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response := UserResponse{Username: username, Success: true}
	json.NewEncoder(w).Encode(response)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hit get user handler")
	var userRequest UserRequest
	err := json.NewDecoder(r.Body).Decode(&userRequest)
	if err != nil {
		jsonErrorResponse(w, "Invalid request body, please include username", http.StatusBadRequest)
		return
	}
	receivedUsername := userRequest.Username
	user, err := GetUserName(receivedUsername)
	if err != nil {
		if err.Error() == "user not found" {
			jsonErrorResponse(w, "User not found", http.StatusNotFound)
		} else {
			log.Printf("Error getting user from database: %v", err)
			jsonErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response := UserResponse{Username: user, Success: true}
	json.NewEncoder(w).Encode(response)
}

// handler for logging in
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user UserDetails
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		jsonErrorResponse(w, "Invalid request body, please include username and password", http.StatusBadRequest)
		return
	}
	username, password := user.Username, user.Password

	storedUser, err := GetUserDetails(username)
	if err != nil {
		if err.Error() == "user not found" {
			jsonErrorResponse(w, "Invalid username or password", http.StatusUnauthorized)
		} else {
			log.Printf("Error fetching user: %v", err)
			jsonErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(password))
	if err != nil {
		jsonErrorResponse(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response := UserResponse{Username: username, Success: true}
	json.NewEncoder(w).Encode(response)
}
