package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"lol-guesser/common"
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

type contextKey string

const usernameKey contextKey = "username"

var (
	jwtSecret       = []byte(os.Getenv("JWT_SECRET"))
	jwtLifeSpan     = time.Hour * 24
	refreshLifeSpan = time.Hour * 24 * 7
)

func ValidateAuthToken(tokenString string) (*Claims, error) {
	// Step 1: Parse and validate the token
	token, err := jwt.ParseWithClaims(
		tokenString, // The token string to validate
		&Claims{},   // The struct to parse the claims into
		func(token *jwt.Token) (interface{}, error) {
			// Step 2: Check signing method is HS256
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || token.Method != jwt.SigningMethodHS256 {
				return nil, common.ErrInvalidSigningMethod
			}
			// Return the secret key used to validate the signature
			return jwtSecret, nil
		},
	)

	// Step 3: Handle parsing errors
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			// Token has expired
			return nil, common.ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			// Signature doesn't match - token has been tampered with
			return nil, common.ErrInvalidSignature
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			// Token isn't in the correct JWT format
			return nil, common.ErrMalformedToken
		}
		return nil, err
	}

	// Step 4: Verify we got valid claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, common.ErrInvalidClaims
	}

	// Token passed all checks
	return claims, nil
}

func ValidateRefreshToken(tokenString string, username string) (*Claims, error) {
	// First check if this refresh token matches what's in the database
	var storedToken string
	err := db.QueryRow("SELECT refresh_token FROM users WHERE username = ?", username).Scan(&storedToken)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	if tokenString != storedToken {
		return nil, common.ErrInvalidSignature
	}

	// Now validate the token structure and claims
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || token.Method != jwt.SigningMethodHS256 {
			return nil, common.ErrInvalidSigningMethod
		}
		return jwtSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, common.ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return nil, common.ErrInvalidSignature
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, common.ErrMalformedToken
		}
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, common.ErrInvalidClaims
	}

	return claims, nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// First check the auth token
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			JsonErrorResponse(w, "No token provided", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			JsonErrorResponse(w, "Invalid auth header format", http.StatusUnauthorized)
			return
		}
		authToken := parts[1]

		claims, err := ValidateAuthToken(authToken)
		if err == nil {
			// Auth token is valid, proceed with request
			ctx := context.WithValue(r.Context(), usernameKey, claims.Username)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// If auth token is invalid, only proceed with refresh token flow if it was an expiration error
		if !errors.Is(err, common.ErrTokenExpired) {
			JsonErrorResponse(w, "Invalid authentication token", http.StatusUnauthorized)
			return
		}

		// Try to get refresh token from cookie
		refreshCookie, err := r.Cookie("refresh_token")
		if err != nil {
			JsonErrorResponse(w, "No refresh token provided", http.StatusUnauthorized)
			return
		}

		// Validate refresh token
		refreshClaims, err := ValidateRefreshToken(refreshCookie.Value, claims.Username)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrTokenExpired):
				JsonErrorResponse(w, "Refresh token expired, please login again", http.StatusUnauthorized)
			case errors.Is(err, common.ErrInvalidSignature):
				JsonErrorResponse(w, "Invalid refresh token", http.StatusUnauthorized)
			default:
				JsonErrorResponse(w, "Invalid refresh token", http.StatusUnauthorized)
			}
			return
		}

		// Generate new auth token
		newAuthToken, err := GenerateAuthToken(refreshClaims.Username)
		if err != nil {
			JsonErrorResponse(w, "Failed to generate new auth token", http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), usernameKey, refreshClaims.Username)
		w.Header().Set("X-New-Auth-Token", newAuthToken)

		// Continue with the request using the refresh token claims
		next.ServeHTTP(w, r.WithContext(ctx))
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

func addRefreshTokentoResponse(w http.ResponseWriter, token string) {
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60,
	}
	http.SetCookie(w, cookie)
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
