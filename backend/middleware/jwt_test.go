package middleware

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"Taxi-app/backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Always set a test key before any test runs
	utils.SetJWTKey(generateTestKeyForInit())
}

func generateTestKeyForInit() *ecdsa.PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	return key
}

func generateValidToken(t *testing.T, userID int, role string) string {
	utils.SetJWTKey(generateTestKeyForInit())
	access, _, _, err := utils.GenerateTokens(userID, role)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	return access
}

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Missing Authorization header", func(t *testing.T) {
		r := gin.New()
		r.Use(AuthMiddleware())
		r.GET("/", func(c *gin.Context) { c.String(200, "ok") })
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, 401, w.Code)
		assert.Contains(t, w.Body.String(), "Missing or invalid auth header")
	})

	t.Run("Invalid Authorization header", func(t *testing.T) {
		r := gin.New()
		r.Use(AuthMiddleware())
		r.GET("/", func(c *gin.Context) { c.String(200, "ok") })
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "NotBearer sometoken")
		r.ServeHTTP(w, req)
		assert.Equal(t, 401, w.Code)
		assert.Contains(t, w.Body.String(), "Missing or invalid auth header")
	})

	t.Run("Invalid token", func(t *testing.T) {
		r := gin.New()
		r.Use(AuthMiddleware())
		r.GET("/", func(c *gin.Context) { c.String(200, "ok") })
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer invalid.token.here")
		r.ServeHTTP(w, req)
		assert.Equal(t, 401, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid token")
	})

	t.Run("Valid token", func(t *testing.T) {
		r := gin.New()
		r.Use(AuthMiddleware())
		called := false
		r.GET("/", func(c *gin.Context) {
			called = true
			userID, _ := c.Get("userID")
			role, _ := c.Get("role")
			c.JSON(200, gin.H{"userID": userID, "role": role})
		})
		w := httptest.NewRecorder()
		token := generateValidToken(t, 123, "admin")
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		r.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.True(t, called, "handler should be called")
		assert.Contains(t, w.Body.String(), "123")
		assert.Contains(t, w.Body.String(), "admin")
	})
}

func TestRequireRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("No role in context", func(t *testing.T) {
		r := gin.New()
		r.Use(RequireRole("admin"))
		r.GET("/", func(c *gin.Context) { c.String(200, "ok") })
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, 403, w.Code)
		assert.Contains(t, w.Body.String(), "No role in token")
	})

	t.Run("Wrong role", func(t *testing.T) {
		r := gin.New()
		r.Use(func(c *gin.Context) { c.Set("role", "user"); c.Next() })
		r.Use(RequireRole("admin"))
		r.GET("/", func(c *gin.Context) { c.String(200, "ok") })
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, 403, w.Code)
		assert.Contains(t, w.Body.String(), "Forbidden")
	})

	t.Run("Correct role", func(t *testing.T) {
		r := gin.New()
		r.Use(func(c *gin.Context) { c.Set("role", "admin"); c.Next() })
		r.Use(RequireRole("admin"))
		called := false
		r.GET("/", func(c *gin.Context) { called = true; c.String(200, "ok") })
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.True(t, called, "handler should be called")
	})
}
