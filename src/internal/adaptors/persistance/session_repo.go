package persistance

import (
    "database/sql"
    "errors"
    "time"
)

type SessionPostgresRepo struct {
    db *sql.DB
}

func NewSessionPostgresRepo(db *sql.DB) *SessionPostgresRepo {
    return &SessionPostgresRepo{db: db}
}

func (r *SessionPostgresRepo) SaveSession(sessionID string, userID int, expiresAt int64) error {
    _, err := r.db.Exec("INSERT INTO sessions (session_id, user_id, expires_at) VALUES ($1, $2, $3) ON CONFLICT (session_id) DO UPDATE SET user_id = $2, expires_at = $3",
        sessionID, userID, expiresAt)
    return err
}

func (r *SessionPostgresRepo) GetSession(sessionID string) (int, error) {
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

func (r *SessionPostgresRepo) DeleteSession(sessionID string) error {
    _, err := r.db.Exec("DELETE FROM sessions WHERE session_id = $1", sessionID)
    return err
}