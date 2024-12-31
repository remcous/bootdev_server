package auth

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJwtCreate(t *testing.T) {
	knownUUIDString := "123e4567-e89b-12d3-a456-426614174000"
	secret := "secret"
	expiresIn := time.Hour
	knownUUID, err := uuid.Parse(knownUUIDString)
	if err != nil {
		t.Fatalf("Failed to parse UUID: %v", err)
	}

	token, err := MakeJWT(knownUUID, secret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create UUID: %v", err)
	}

	if token == "" || token == knownUUIDString {
		t.Fatalf("blank or unencrypted token")
	}
}

func TestJwtCreateAndRead(t *testing.T) {
	knownUUIDString := "123e4567-e89b-12d3-a456-426614174000"
	secret := "secret"
	expiresIn := time.Hour
	knownUUID, err := uuid.Parse(knownUUIDString)
	if err != nil {
		t.Fatalf("Failed to parse UUID: %v", err)
	}

	token, err := MakeJWT(knownUUID, secret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create UUID: %v", err)
	}

	userID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("Failed to read token: %v", err)
	}

	if userID != knownUUID {
		t.Fatalf("decrypted token does not match initial uuid")
	}
}

func TestJwtExpired(t *testing.T) {
	knownUUIDString := "123e4567-e89b-12d3-a456-426614174000"
	secret := "secret"
	expiresIn := time.Millisecond * 500
	knownUUID, err := uuid.Parse(knownUUIDString)
	if err != nil {
		t.Fatalf("Failed to parse UUID: %v", err)
	}

	token, err := MakeJWT(knownUUID, secret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create UUID: %v", err)
	}

	time.Sleep(time.Second)

	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Fatalf("Expected token to be expired")
	}
}

func TestJwtWrongSecret(t *testing.T) {
	knownUUIDString := "123e4567-e89b-12d3-a456-426614174000"
	secret1 := "secret"
	secret2 := "wrong"
	expiresIn := time.Hour
	knownUUID, err := uuid.Parse(knownUUIDString)
	if err != nil {
		t.Fatalf("Failed to parse UUID: %v", err)
	}

	token, err := MakeJWT(knownUUID, secret1, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create UUID: %v", err)
	}

	_, err = ValidateJWT(token, secret2)
	if err == nil {
		t.Fatalf("Should fail to parse token with wrong secret")
	}
}

func TestBearerTokenValid(t *testing.T) {
	expectedToken := "1234"
	bearer := fmt.Sprintf("Bearer %s", expectedToken)
	header := http.Header{}
	header.Add("Authorization", bearer)

	token, err := GetBearerToken(header)
	if err != nil {
		t.Fatalf("Failed to get bearer token, %v", err)
	}

	if token != expectedToken {
		t.Fatalf("expected [%s] but got [%s]", expectedToken, token)
	}
}

func TestBearerTokenUnauthorized(t *testing.T) {
	header := http.Header{}

	_, err := GetBearerToken(header)
	if err == nil {
		t.Fatalf("expected GetBearerToken to fail with no Authorization header")
	}
}
