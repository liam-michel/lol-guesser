package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	//"github.com/joho/godotenv"
)

//use godotenv to read the env variables from a different directory

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var (
	jwtSecret       = []byte(os.Getenv("JWT_SECRET"))
	jwtLifeSpan     = time.Hour * 24
	refreshLifeSpan = time.Hour * 24 * 7
)

func isAuthTokenValid(claims *Claims) (bool, error) {
	if claims == nil {
		return false, fmt.Errorf("claims cannot be nil")
	}

	// Check expiration
	if claims.ExpiresAt.Before(time.Now()) {
		return false, nil // token is expired
	}

	return true, nil // claims are valid and not expired
}

func isRefreshTokenValid(claims *Claims, username string) {

}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "No token provided", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid auth header format", http.StatusUnauthorized)
			return
		}
		tokenString := parts[1]

		token, err := ParseToken(tokenString)
		if err != nil {
			valid, _ := isAuthTokenValid(token)
			if !valid {
				//check for the refresh token
				refreshToken, err := r.Cookie("refresh_token")
				if err != nil {
					http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
					return
				}
				//check if the refresh token is valid
				newToken, refreshErr := RefreshAuthToken(tokenString)
				if refreshErr != nil {
					http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
					return
				}
				w.Header().Set("Authorization", "Bearer "+newToken)
			} else {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
		}

		// Add claims to context or attach to request if needed
		// ctx := context.WithValue(r.Context(), "userClaims", claims)
		// r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func GenerateRefreshToken(username string) (string, error) {
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshLifeSpan)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)

}

func addRefreshTokentoResponse(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60,
	})

}

func GenerateAuthToken(username string) (string, error) {
	//define claims
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtLifeSpan)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err

	}
	return signedToken, nil
}

func RefreshAuthToken(tokenString string) (string, error) {
	//parse the passed in token
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	//refresh the token by generating a new one
	newToken, err := GenerateAuthToken(claims.Username)
	if err != nil {
		return "", err
	}
	return newToken, nil

}

func ParseToken(tokenString string) (*Claims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	// Extract and validate claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
