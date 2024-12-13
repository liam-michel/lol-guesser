package handlejwt

import (
	"os"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret") // Set up a test secret
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))

	token, err := GenerateToken("testuser")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if token == "" {
		t.Fatalf("Expected a token, got an empty string")
	}
}

func TestParseToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret") // Set up a test secret
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))

	token, _ := GenerateToken("testuser")

	claims, err := ParseToken(token)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if claims.Username != "testuser" {
		t.Fatalf("Expected username 'testuser', got '%s'", claims.Username)
	}
}

func TestRefreshToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret") // Set up a test secret
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))

	token, _ := GenerateToken("testuser")

	newToken, err := RefreshToken(token)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if newToken == "" {
		t.Fatalf("Expected a new token, got an empty string")
	}
}
