package main

import (
    "database/sql"
    "fmt"
    "log"
	"github.com/joho/godotenv"
	"os"
	"golang.org/x/crypto/bcrypt"
    _ "github.com/go-sql-driver/mysql"

)

// Global DB variable
var db *sql.DB

// Function to initialize the DB connection
func InitDB() {
    var err error
	err = godotenv.Load("../../.env")
	if(err != nil){	
		log.Fatal("Error loading .env file")
	}

    // Get database credentials from environment variables or hardcode for now
    dbPassword := os.Getenv("MYSQLPASSWORD")// Replace with your real password or use os.Getenv for environment variable
    dbName := "lol_users"         // Replace with your database name

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

func addUser(username string, password string) error {
	//start by checking if the user already exists in the database

	user, err := getUser(username)
	if err != nil {
		fmt.Errorf("Error getting user from database: %v", err)
	}
	if user != ""{
		fmt.Errorf("User already exists in the database")
	}
	//start with hashing the password with bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {	
		fmt.Errorf("Error hashing password: %v", err)
	}
	//print the hashed password
	fmt.Println("Hashed password: ", string(hash))

	//insert the user into the database
	_, err = db.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", username, hash)
	if err != nil {
		fmt.Errorf("Error inserting user into database: %v", err)
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

func getUser(username string) (string, error){
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


func main(){
	InitDB()
	addUser("newtest", "testpassword")
	// users, err := GetUsers()
	// if err != nil {	
	// 	fmt.Errorf("Error getting users: %v", err)
	// }
	// fmt.Println("Users: ", users)


}