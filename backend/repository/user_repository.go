package repository

import (
	"database/sql"
)

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	Role         string `json:"role"`
}

type UserRepository interface {
	CreateUser(username, passwordHash, role string) error
	GetUserByUsername(username string) (*User, error)
	SaveRefreshToken(tokenID string, userID int, expiresAt string) error
	DeleteRefreshToken(tokenID string) error
	GetUserIDByRefreshToken(tokenID string) (int, error)
	GetUserByID(userID int) (*User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(username, passwordHash, role string) error {
	_, err := r.db.Exec(
		"INSERT INTO users (username, password_hash, role) VALUES ($1, $2, $3)",
		username, passwordHash, role,
	)
	return err
}

func (r *userRepository) GetUserByUsername(username string) (*User, error) {
	var user User
	err := r.db.QueryRow(
		"SELECT id, username, password_hash, role FROM users WHERE username = $1",
		username,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Save a new refresh token
func (r *userRepository) SaveRefreshToken(tokenID string, userID int, expiresAt string) error {
	_, err := r.db.Exec(
		`INSERT INTO refresh_tokens (token_id, user_id, issued_at, expires_at)
		 VALUES ($1, $2, NOW(), $3)`,
		tokenID, userID, expiresAt,
	)
	return err
}

// Delete a refresh token (e.g., on refresh or logout)
func (r *userRepository) DeleteRefreshToken(tokenID string) error {
	_, err := r.db.Exec(
		`DELETE FROM refresh_tokens WHERE token_id=$1`,
		tokenID,
	)
	return err
}

// Get user ID by refresh token (and check expiry)
func (r *userRepository) GetUserIDByRefreshToken(tokenID string) (int, error) {
	var userID int
	err := r.db.QueryRow(
		`SELECT user_id FROM refresh_tokens WHERE token_id=$1 AND expires_at > NOW()`,
		tokenID,
	).Scan(&userID)
	return userID, err
}

func (r *userRepository) GetUserByID(userID int) (*User, error) {
	var user User
	err := r.db.QueryRow(
		"SELECT id, username, password_hash, role FROM users WHERE id = $1",
		userID,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
