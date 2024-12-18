package auth

import (
	"database/sql"
	"lol-guesser/common"

	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func SetDB(database *sql.DB) {
	db = database
}

func CreateUser(username string, password string, refreshToken string) error {
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
	_, err = db.Exec("INSERT INTO users (username, password_hash, refresh_token) VALUES (?, ?, ?)", username, hash, refreshToken)
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

func GetUserDetails(username string) (common.UserDetails, error) {
	var user common.UserDetails
	err := db.QueryRow("SELECT username, password_hash FROM users WHERE username = ?", username).Scan(&user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return common.UserDetails{}, fmt.Errorf("user not found")
		}
		return common.UserDetails{}, err
	}
	return user, nil
}

func updateRefreshToken(username string, refreshToken string) error {
	_, err := db.Exec("UPDATE users SET refresh_token = ? WHERE username = ?", refreshToken, username)
	if err != nil {
		return fmt.Errorf("error updating refresh token: %v", err)
	}
	return nil
}
