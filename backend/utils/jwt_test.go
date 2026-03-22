package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func generateTestKey(t *testing.T) *ecdsa.PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate test ECDSA key: %v", err)
	}
	return key
}

func TestJWTUtils(t *testing.T) {
	key := generateTestKey(t)
	SetJWTKey(key)

	t.Run("Generate and parse valid tokens", func(t *testing.T) {
		userID := 42
		role := "user"
		access, refresh, refreshJTI, err := GenerateTokens(userID, role)
		if err != nil {
			t.Fatalf("GenerateTokens failed: %v", err)
		}
		if access == "" || refresh == "" || refreshJTI == "" {
			t.Fatal("Tokens or JTI should not be empty")
		}

		claims, err := ParseJWT(access)
		if err != nil {
			t.Fatalf("ParseJWT failed: %v", err)
		}
		if claims.UserID != userID || claims.Role != role {
			t.Errorf("Claims mismatch: got userID=%d, role=%s", claims.UserID, claims.Role)
		}

		rclaims, err := ParseRefreshToken(refresh)
		if err != nil {
			t.Fatalf("ParseRefreshToken failed: %v", err)
		}
		if rclaims.Subject != "42" || rclaims.ID != refreshJTI {
			t.Errorf("Refresh claims mismatch: got subject=%s, id=%s", rclaims.Subject, rclaims.ID)
		}
	})

	t.Run("Invalid access token", func(t *testing.T) {
		_, err := ParseJWT("invalid.token.here")
		if err == nil {
			t.Error("Expected error for invalid token")
		}
	})

	t.Run("Invalid refresh token", func(t *testing.T) {
		_, err := ParseRefreshToken("invalid.token.here")
		if err == nil {
			t.Error("Expected error for invalid refresh token")
		}
	})

	t.Run("Tampered access token", func(t *testing.T) {
		userID := 1
		role := "admin"
		access, _, _, err := GenerateTokens(userID, role)
		if err != nil {
			t.Fatalf("GenerateTokens failed: %v", err)
		}
		
		parts := strings.Split(access, ".")
		if len(parts) != 3 {
			t.Fatal("Invalid JWT format")
		}
		parts[1] = "tampered"
		tampered := strings.Join(parts, ".")
		_, err = ParseJWT(tampered)
		if err == nil {
			t.Error("Expected error for tampered token")
		}
	})

	t.Run("Expired access token", func(t *testing.T) {
		
		userID := 2
		role := "user"
		now := time.Now()
		claims := &Claims{
			UserID: userID,
			Role:   role,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(now.Add(1 * time.Second)),
				IssuedAt:  jwt.NewNumericDate(now),
				ID:        "test-expired",
			},
		}
		tok := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
		signed, err := tok.SignedString(key)
		if err != nil {
			t.Fatalf("Failed to sign short-lived token: %v", err)
		}
		time.Sleep(2 * time.Second)
		_, err = ParseJWT(signed)
		if err == nil {
			t.Error("Expected error for expired token")
		}
	})

	t.Run("Wrong key for parsing", func(t *testing.T) {
		userID := 3
		role := "user"
		access, _, _, err := GenerateTokens(userID, role)
		if err != nil {
			t.Fatalf("GenerateTokens failed: %v", err)
		}
		
		wrongKey := generateTestKey(t)
		SetJWTKey(wrongKey)
		_, err = ParseJWT(access)
		if err == nil {
			t.Error("Expected error for token signed with different key")
		}
	})
}
