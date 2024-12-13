package handlejwt

import (
	"errors"
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
	jwtSecret   = []byte(os.Getenv("JWT_SECRET"))
	jwtLifeSpan = time.Hour * 24
)

// func init() {
// 	err := godotenv.Load("../.env")
// 	if err != nil {
// 		panic(err)
// 	}
// 	//load the jwt secret from env variables
// }

func isTokenExpired(err error) bool {
	return strings.Contains(err.Error(), "expired")
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

		_, err := ParseToken(tokenString)
		if err != nil {
			if isTokenExpired(err) {
				newToken, refreshErr := RefreshToken(tokenString)
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

func GenerateToken(username string) (string, error) {
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

func RefreshToken(tokenString string) (string, error) {
	//parse the passed in token
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	//refresh the token by generating a new one
	newToken, err := GenerateToken(claims.Username)
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
