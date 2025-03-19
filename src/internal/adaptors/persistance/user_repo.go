package persistance

import (
    "database/sql"
    "errors"
    "task-manager/src/internal/core/user"
    "golang.org/x/crypto/bcrypt"
)

type UserPostgresRepo struct {
    db         *sql.DB
    sessionRepo *SessionPostgresRepo // Add this
}

func NewUserPostgresRepo(db *sql.DB) *UserPostgresRepo {
    repo := &UserPostgresRepo{
        db:         db,
        sessionRepo: NewSessionPostgresRepo(db), // Initialize here
    }
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
    for _, query := range queries {
        _, err := r.db.Exec(query)
        if err != nil {
            panic("Failed to initialize user DB: " + err.Error())
        }
    }
    hashedAdminPass, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
    hashedUserPass, _ := bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)
    r.db.Exec("INSERT INTO users (username, password, role) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING", "admin", string(hashedAdminPass), "admin")
    r.db.Exec("INSERT INTO users (username, password, role) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING", "user1", string(hashedUserPass), "user")
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

func (r *UserPostgresRepo) GetUserRole(userID int) (string, error) {
    var role string
    err := r.db.QueryRow("SELECT role FROM users WHERE id = $1", userID).Scan(&role)
    if err == sql.ErrNoRows {
        return "", errors.New("user not found")
    }
    return role, err
}