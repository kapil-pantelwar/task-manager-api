package repository

import (
    "database/sql"
    "errors"
    "log"
    "task-manager/domain"
    "time"
	"golang.org/x/crypto/bcrypt"
)

type UserPostgresRepo struct {
    db *sql.DB
}

func NewUserPostgresRepo(db *sql.DB) *UserPostgresRepo {
    repo := &UserPostgresRepo{db: db}
    repo.initDB()
    return repo
}

func (r *UserPostgresRepo) initDB() {
    queries := []string{
        `CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            username TEXT UNIQUE NOT NULL,
            password TEXT NOT NULL,
            role TEXT NOT NULL
        )`,
        `CREATE TABLE IF NOT EXISTS sessions (
            session_id TEXT PRIMARY KEY,
            user_id INTEGER REFERENCES users(id),
            expires_at BIGINT
        )`,
    }
    for i, query := range queries {
        _, err := r.db.Exec(query)
        if err != nil {
            log.Printf("Failed query %d: %s", i, query)
            panic("Failed to initialize user DB: " + err.Error())
        }
    }
    hashedAdminPass, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
    if err != nil {
        panic("Failed to hash admin password: " + err.Error())
    }
    hashedUserPass, err := bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)
    if err != nil {
        panic("Failed to hash user password: " + err.Error())
    }
    _, err = r.db.Exec("INSERT INTO users (username, password, role) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING", "admin", string(hashedAdminPass), "admin")
    if err != nil {
        panic("Failed to insert admin user: " + err.Error())
    }
    _, err = r.db.Exec("INSERT INTO users (username, password, role) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING", "user1", string(hashedUserPass), "user")
    if err != nil {
        panic("Failed to insert user1: " + err.Error())
    }
}

func (r *UserPostgresRepo) FindUserByUsername(username string) (domain.User, error) {
    var user domain.User
    err := r.db.QueryRow("SELECT id, username, password, role FROM users WHERE username = $1", username).
        Scan(&user.ID, &user.Username, &user.Password, &user.Role)
    if err == sql.ErrNoRows {
        return domain.User{}, errors.New("user not found")
    }
    return user, err
}

func (r *UserPostgresRepo) SaveSession(sessionID string, userID int, expiresAt int64) error {
    _, err := r.db.Exec("INSERT INTO sessions (session_id, user_id, expires_at) VALUES ($1, $2, $3) ON CONFLICT (session_id) DO UPDATE SET user_id = $2, expires_at = $3",
        sessionID, userID, expiresAt)
    return err
}

func (r *UserPostgresRepo) GetSession(sessionID string) (int, error) {
    var userID int
    var expiresAt int64
    err := r.db.QueryRow("SELECT user_id, expires_at FROM sessions WHERE session_id = $1", sessionID).
        Scan(&userID, &expiresAt)
    if err == sql.ErrNoRows {
        return 0, errors.New("session not found")
    }
    if err != nil {
        return 0, err
    }
    if expiresAt < time.Now().Unix() {
        return 0, errors.New("session expired")
    }
    return userID, nil
}

func (r *UserPostgresRepo) GetUserRole(userID int) (string, error) {
    var role string
    err := r.db.QueryRow("SELECT role FROM users WHERE id = $1", userID).Scan(&role)
    if err == sql.ErrNoRows {
        return "", errors.New("user not found")
    }
    return role, err
}

func (r *UserPostgresRepo) DeleteSession(sessionID string) error {
	_,err := r.db.Exec("DELETE FROM sessions WHERE session_id = $1",sessionID)
	return err
}