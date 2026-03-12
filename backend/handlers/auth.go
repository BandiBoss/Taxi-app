package handlers

import (
	"Taxi-app/backend/utils"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/lib/pq"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"Taxi-app/backend/repository"
)

// RegisterRequest represents the request body for registration and login
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// @Summary Register a new user
// @Description Register a new user with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param data body RegisterRequest true "User credentials"
// @Success 201 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/register [post]
func Register(repo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.BindJSON(&req); err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err)
			return
		}

		// Hash the password
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Could not hash password", err)
			return
		}

		err = repo.CreateUser(req.Username, string(hash), "user")
		if err != nil {
			var pgErr *pq.Error
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				utils.ErrorResponse(c, http.StatusConflict, "That username is already taken. Please choose another.", err)
				return
			}
			utils.ErrorResponse(c, http.StatusInternalServerError, "Could not create user", err)
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User created"})
	}
}

// @Summary Login
// @Description Login with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param data body RegisterRequest true "User credentials"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/login [post]
func Login(repo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.BindJSON(&req); err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err)
			return
		}
		user, err := repo.GetUserByUsername(req.Username)
		if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Wrong credentials", err)
			return
		}
		accessToken, refreshToken, jti, err := utils.GenerateTokens(user.ID, user.Role)
		if err != nil {
			log.Printf("Login: token generation failed: %v", err)
			utils.ErrorResponse(c, http.StatusInternalServerError, "Token generation failed", err)
			return
		}
		expiresAt := time.Now().Add(7 * 24 * time.Hour).Format("2006-01-02 15:04:05")
		err = repo.SaveRefreshToken(jti, user.ID, expiresAt)
		if err != nil {
			log.Printf("Login: could not save refresh token: %v", err)
		}
		c.SetCookie(
			"refresh_token",
			refreshToken,
			int(utils.RefreshTokenTTL.Seconds()),
			"/",
			"",
			false,
			true,
		)
		c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
	}
}

// RefreshToken issues a new access token using a valid refresh token.
//
// @Summary      Refresh access token
// @Description  Issues a new access token if the provided refresh token (in cookie) is valid and not expired/revoked.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]string "New access token"
// @Failure      401 {object} map[string]string "Invalid or missing refresh token"
// @Failure      500 {object} map[string]string "Server error"
// @Router       /api/refresh [post]
// @Security     ApiKeyAuth
func RefreshToken(repo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		rt, err := c.Cookie("refresh_token")
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "No refresh token", err)
			return
		}

		claims, err := utils.ParseRefreshToken(rt)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid refresh token", err)
			return
		}

		userID, err := repo.GetUserIDByRefreshToken(claims.ID)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Refresh token not recognized", err)
			return
		}

		user, err := repo.GetUserByID(userID)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "User not found", err)
			return
		}

		newAccess, newRefresh, newJTI, err := utils.GenerateTokens(userID, user.Role)
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Token generation failed", err)
			return
		}

		// Remove the old refresh token and save the new one
		if err := repo.DeleteRefreshToken(claims.ID); err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete old refresh token", err)
			return
		}
		expiresAt := time.Now().Add(utils.RefreshTokenTTL).Format("2006-01-02 15:04:05")
		if err := repo.SaveRefreshToken(newJTI, userID, expiresAt); err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to save new refresh token", err)
			return
		}

		c.SetCookie("refresh_token", newRefresh, int(utils.RefreshTokenTTL.Seconds()),
			"/", "", true, true,
		)
		c.JSON(http.StatusOK, gin.H{"access_token": newAccess})
	}
}

// Logout revokes the current refresh token and clears the cookie.
//
// @Summary      Logout user
// @Description  Revokes the refresh token and clears the refresh_token cookie.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]string "Logged out"
// @Router       /api/logout [post]
// @Security     ApiKeyAuth
func Logout(repo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		rt, err := c.Cookie("refresh_token")
		if err == nil && rt != "" {
			claims, err := utils.ParseRefreshToken(rt)
			if err == nil {
				_ = repo.DeleteRefreshToken(claims.ID)
			}
		}
		// Clear the cookie
		c.SetCookie("refresh_token", "", -1, "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
	}
}
