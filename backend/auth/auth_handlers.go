package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"lol-guesser/common"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func JsonErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := map[string]string{"error": message}
	json.NewEncoder(w).Encode(response)
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hit create user handler")
	var createUserRequest common.CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&createUserRequest)
	if err != nil {
		JsonErrorResponse(w, "Invalid request body, please include username and password", http.StatusBadRequest)
		return
	}
	username, password := createUserRequest.Username, createUserRequest.Password

	//generate a jwt token (once we create the user, we want to log them in)
	authToken, err := GenerateAuthToken(username)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		JsonErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	refreshToken, err := GenerateRefreshToken(username)
	if err != nil {
		log.Printf("Error generating the refresh token %v", err)
		JsonErrorResponse(w, "Internal server error", http.StatusInternalServerError)
	}

	err = CreateUser(username, password, refreshToken)
	if err != nil {
		if err.Error() == "user already exists in the database" {
			JsonErrorResponse(w, "Username is already taken", http.StatusConflict)
		} else {
			log.Printf("Error adding user to the database: %v", err)
			JsonErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	addRefreshTokentoResponse(w, refreshToken)

	//return the response with the token included
	response := common.JwtResponse{
		Username:  username,
		AuthToken: authToken,
		Success:   "true",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handler for logging in
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user common.UserDetails
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		JsonErrorResponse(w, "Invalid request body, please include username and password", http.StatusBadRequest)
		return
	}
	username, password := user.Username, user.Password

	storedUser, err := GetUserDetails(username)
	if err != nil {
		if err.Error() == "user not found" {
			JsonErrorResponse(w, "Invalid username or password", http.StatusUnauthorized)
		} else {
			log.Printf("Error fetching user: %v", err)
			JsonErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(password))
	if err != nil {
		JsonErrorResponse(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	//generate a jwt token
	token, err := GenerateAuthToken(storedUser.Username)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		JsonErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	//replace the refresh token if its expired
	refreshtoken, err := GenerateRefreshToken(username)
	if err != nil {
		log.Printf("Error generating the refresh token %v", err)
		JsonErrorResponse(w, "Internal server error", http.StatusInternalServerError)
	}
	addRefreshTokentoResponse(w, refreshtoken)

	//return the response with the token included
	response := common.JwtResponse{
		Username:  storedUser.Username,
		AuthToken: token,
		Success:   "true",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
// 	//generate a new token
// 	//get the username from the request
// 	var user common.Username
// 	err := json.NewDecoder(r.Body).Decode(&user)
// 	if err != nil {
// 		JsonErrorResponse(w, "Invalid request body, please include username", http.StatusBadRequest)
// 		return
// 	}
// 	username := user.Username
// 	//check if the user is valid
// 	exists, err := checkUserExists(username)
// 	if err != nil {
// 		log.Printf("Error checking if user exists: %v", err)
// 		JsonErrorResponse(w, "Internal server error", http.StatusInternalServerError)
// 		return
// 	}
// 	if !exists {
// 		JsonErrorResponse(w, "User does not exist", http.StatusNotFound)
// 		return
// 	}

// 	newRefreshToken, err := GenerateRefreshToken(username)
// 	if err != nil {
// 		log.Printf("Error generating refresh token: %v", err)
// 		JsonErrorResponse(w, "Internal server error", http.StatusInternalServerError)
// 		return
// 	}
// 	updateRefreshToken(username, newRefreshToken)
// 	addRefreshTokentoResponse(w, newRefreshToken)
// 	//send an empty response, as new refresh token is returned in cookies
// 	w.WriteHeader(http.StatusOK)
// }
