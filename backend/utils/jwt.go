package utils

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var JwtKey *ecdsa.PrivateKey

const AccessTokenTTL = 15 * time.Minute
const RefreshTokenTTL = 7 * 24 * time.Hour

type Claims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func LoadECDSAPrivateKey(path string) *ecdsa.PrivateKey {
	keyData, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read ECDSA private key: %v", err)
	}
	block, _ := pem.Decode(keyData)
	if block == nil {
		log.Fatal("Failed to parse PEM block containing the key")
	}
	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("Failed to parse EC private key: %v", err)
	}
	return key
}

func SetJWTKey(key *ecdsa.PrivateKey) {
	JwtKey = key
}

func ensureKey() {
	if JwtKey == nil {
		JwtKey = LoadECDSAPrivateKey("secrets/ec256-private.pem")
	}
}

func GenerateTokens(userID int, role string) (accessToken string, refreshToken string, refreshJTI string, err error) {
	ensureKey()
	now := time.Now()

	atClaims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuid.NewString(),
		},
	}
	at := jwt.NewWithClaims(jwt.SigningMethodES256, atClaims)
	accessToken, err = at.SignedString(JwtKey)
	if err != nil {
		return
	}
	refreshJTI = uuid.NewString()
	rtClaims := &jwt.RegisteredClaims{
		Subject:   fmt.Sprint(userID),
		ExpiresAt: jwt.NewNumericDate(now.Add(RefreshTokenTTL)),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        refreshJTI,
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodES256, rtClaims)
	refreshToken, err = rt.SignedString(JwtKey)
	return
}

func ParseJWT(tokenStr string) (*Claims, error) {
	ensureKey()
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return &JwtKey.PublicKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid access token")
}

func ParseRefreshToken(tokenStr string) (*jwt.RegisteredClaims, error) {
	ensureKey()
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{},
		func(t *jwt.Token) (interface{}, error) { return &JwtKey.PublicKey, nil },
	)
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid refresh token")
}
